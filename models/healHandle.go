package models

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"fmt"

	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"healing2020/pkg/tools"

	"github.com/jinzhu/gorm"
)

func GetPhone(info tools.RedisUser) string {
	return info.Phone
}

type RecordResp struct {
	SongId     uint   `json:"songId"`
	SingerId   uint   `json:"singId"`
	Praise     int    `json:"praise"`
	User       string `json:"user"`
	Source     string `json:"source"`
	SongAvatar string `json:"songAvatar"`
	IsPraise   bool   `json:"isPraise"`
}
type ResultResp struct {
	VodId     uint      `json:"id"`
	Time      time.Time `json:"time"`
	Singer    string    `json:"singer"`
	More      string    `json:"more"`
	Name      string    `json:"name"`
	VodUser   string    `json:"vodUser"`
	VodAvatar string    `json:"vodAvatar"`
	VodUserId uint      `json:"vodUserId"`
	Style     string    `json:"style"`
	Language  string    `json:"language"`
	AllSongs  []RecordResp
	Err       error `json:"err"`
}

func GetRecord(id string, myID uint) ResultResp {
	intId, _ := strconv.Atoi(id)
	vodId := uint(intId)

	db := setting.MysqlConn()

	var resultResp ResultResp

	resultResp.VodId = vodId

	var vod statements.Vod
	result := db.Model(&statements.Vod{}).Where("id=?", vodId).First(&vod)
	resultResp.Err = result.Error

	if result.Error != nil {
		return resultResp
	}

	resultResp.Singer = vod.Singer
	resultResp.More = vod.More
	resultResp.Time = vod.CreatedAt
	resultResp.Name = vod.Name
	resultResp.Style = vod.Style
	resultResp.Language = vod.Language

	userId := vod.UserId
	var user statements.User
	db.Model(&statements.User{}).Select("id, sex, avatar,nick_name, setting1").Where("id =?", userId).First(&user)

	resultResp.VodUserId = user.ID
	resultResp.VodUser = user.NickName
	resultResp.VodAvatar = user.Avatar

    if user.Setting1 == 0 {
        resultResp.VodAvatar = tools.GetAvatarUrl(user.Sex)
    }

	if vod.HideName == 1 {
		resultResp.VodUser = "匿名用户"
		resultResp.VodAvatar = tools.GetAvatarUrl(user.Sex)
	}

	count := 0
	var allSong []statements.Song
	recordsToVod := db.Model(&statements.Song{}).Where("vod_id = ? and is_hide = 0", vodId).Count(&count).Find(&allSong)
	var recordResp []RecordResp = make([]RecordResp, count)

	if count == 0 {
		return resultResp
	}

	rows, _ := recordsToVod.Rows()
	defer rows.Close()

	i := 0
	for rows.Next() {
		var songRows statements.Song
		db.ScanRows(rows, &songRows)
		recordResp[i].SongId = songRows.ID
		recordResp[i].Praise = GetPraiseCount("song", songRows.ID)
		recordResp[i].Source = songRows.Source

		var userRows statements.User
		db.Model(&statements.User{}).Select("id, setting1, avatar, sex, nick_name").Where("id = ?", songRows.UserId).First(&userRows)
		recordResp[i].SingerId = userRows.ID
		recordResp[i].User = userRows.NickName
		recordResp[i].SongAvatar = userRows.Avatar

        if userRows.Setting1 == 0 {
            recordResp[i].SongAvatar = tools.GetAvatarUrl(userRows.Sex) 
        }

		recordResp[i].IsPraise, _ = HasPraise(2, myID, songRows.ID)
		i++
	}
	resultResp.AllSongs = recordResp

	return resultResp
}

func CreateRecord(id string, source string, uid uint, isHide int) (string, error) {
	intId, _ := strconv.Atoi(id)
	vodId := uint(intId)
	db := setting.MysqlConn()
	userId := uid

	var vod statements.Vod
	result1 := db.Model(&statements.Vod{}).Where("ID=?", vodId).First(&vod)
	if result1.Error != nil {
		return "", errors.New("vod_id is unvalid")
	}
	var song statements.Song
	song.VodId = vodId
	song.UserId = userId
	song.VodSend = vod.UserId
	song.Name = vod.Name
	song.Praise = 0
	song.Source = source
	song.Style = vod.Style
	song.Language = vod.Language
	song.IsHide = isHide

	err := db.Model(&statements.Song{}).Create(&song).Error

	FinishTask("3", userId)

	return song.Name, err
}

func HasPraise(types int, userid uint, id uint) (bool, uint) {
	db := setting.MysqlConn()

	var praise statements.Praise
	result := db.Model(&statements.Praise{}).Where("is_cancel = 0 and praise_id = ? and user_id = ? and type = ?", id, userid, types).First(&praise)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, 0
	}
	return true, praise.ID
}

func IsPraiseCancel(types int, userid, uint, id uint) bool {
	db := setting.MysqlConn()

	var praise statements.Praise
	result := db.Model(&statements.Praise{}).Where("is_cancel = 1 and praise_id = ? and user_id = ? and type = ?", id, userid, types).First(&praise)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func CancelPraise(userid uint, strId string, types string) error {
	intId, _ := strconv.Atoi(strId)
	id := uint(intId)
	typesInt, _ := strconv.Atoi(types)
	db := setting.MysqlConn()

	hasPraise, praiseId := HasPraise(typesInt, userid, id)
	if !hasPraise {
		return errors.New("item does not be praised")
	}

	var praise statements.Praise
	db.Model(&statements.Praise{}).Where("id = ?", praiseId).First(&praise)
	praise.IsCancel = 1
	err := db.Save(&praise).Error

	return err
}

type PraiseData struct {
	MyID     uint
	TargetID uint
	Type     string
	Msg      string
}

func AddPraise(userid uint, strId string, types string) (error, PraiseData) {
	intId, _ := strconv.Atoi(strId)
	id := uint(intId)
	typesInt, _ := strconv.Atoi(types)
	db := setting.MysqlConn()

	hasPraise, _ := HasPraise(typesInt, userid, id)
	if hasPraise {
		return errors.New("can not praise repeatedly"), PraiseData{}
	}

	var praiseData PraiseData
	praiseData.Type = types
	praiseData.MyID = userid
	praiseData.TargetID = GetTargetId(typesInt, id)
	if types == "2" {
		var song statements.Song
		db.Select("name").Where("id = ?", id).First(&song)
		praiseData.Msg = song.Name
	} else if types == "1" {
		var deliver statements.Deliver
		db.Select("text_field").Where("id = ?", id).First(&deliver)
		//split deliver to 8
		splitDeliver := strings.Split(deliver.TextField, "")
		if len(splitDeliver) > 8 {
			deliver.TextField = strings.Join(splitDeliver[:8], "") + "..."
		}
		praiseData.Msg = deliver.TextField
	}

	var praise statements.Praise
	praise.UserId = userid
	praise.Type = typesInt
	praise.PraiseId = id

	err := db.Model(&statements.Praise{}).Create(&praise).Error

	//if SyncLock(userid) {
	//	FinishTask("5", userid)
	//}
	SyncLock(userid)

	return err, praiseData
}

type Target struct {
	TargetId uint `gorm:"column:user_id"`
}

func GetTargetId(types int, id uint) uint {
	db := setting.MysqlConn()
	var target Target
	if types == 1 {
		db.Table("deliver").Select("user_id").Where("id = ?", id).Scan(&target)
		return target.TargetId
	}
	if types == 2 {
		db.Table("song").Select("user_id").Where("id = ?", id).Scan(&target)
		return target.TargetId
	}
	if types == 3 {
		db.Table("special").Select("user_id").Where("id = ?", id).Scan(&target)
		return target.TargetId
	}
	return 0
}

func SyncLock(userid uint) {
	client := setting.RedisConn()
	num := client.Incr(fmt.Sprintf("healing2020:user:%d:praised", userid)).Val()
	if num >= 3 {
		FinishTask("5", userid)
	}
}

func CreateVod(uid uint, singer string, style string, language string, name string, more string) error {
	db := setting.MysqlConn()

	var userOther statements.UserOther
	db.Select("remain_sing").Where("user_id = ?", uid).First(&userOther)
	if userOther.RemainSing <= 0 {
		return errors.New("已无点歌次数！")
	}
	var vod statements.Vod
	vod.UserId = uid
	vod.More = more
	vod.Name = name
	vod.Singer = singer
	vod.Style = style
	vod.Language = language
	err := db.Model(&statements.Vod{}).Create(&vod).Error

	FinishTask("2", uid)

	return err
}

//根据vod_id获取user_id
func SelectUserIDByVodID(vod_id uint) (uint, error) {
	db := setting.MysqlConn()

	var vod statements.Vod
	err := db.Select("user_id").Where("id=?", vod_id).First(&vod).Error
	return vod.UserId, err
}

//所有表查点赞总数的函数
func GetPraiseCount(table string, id uint) int {
	db := setting.MysqlConn()

	types := ""
	switch table {
	case "deliver":
		types = "1"
		break
	case "song":
		types = "2"
		break
	case "special":
		types = "3"
		break
	case "comment":
		types = "4"
		break
	}

	count := 0
	db.Model(&statements.Praise{}).Where("type = ? and praise_id = ? and is_cancel = 0", types, id).Count(&count)

	return count
}

func SyncPraise(id uint, table string) {
	db := setting.MysqlConn()

	if table == "" {
		return
	}
	count := GetPraiseCount(table, id)
	if table == "deliver" {
		var deliver statements.Deliver
		db.Model(&statements.Deliver{}).Where("id = ?", id).First(&deliver)
		deliver.Praise = count
		db.Save(&deliver)
	}
	if table == "song" {
		var song statements.Song
		db.Model(&statements.Song{}).Where("id = ?", id).First(&song)
		song.Praise = count
		db.Save(&song)
	}
	if table == "special" {
		var special statements.Special
		db.Model(&statements.Special{}).Where("id = ?", id).First(&special)
		special.Praise = count
		db.Save(&special)
	}
}

func AutoSyncPraise() {
	db := setting.MysqlConn()

	count := 0
	db.Model(&statements.Deliver{}).Count(&count)
	for i := 0; i < count; i++ {
		SyncPraise(uint(i+1), "deliver")
	}
	db.Model(&statements.Song{}).Count(&count)
	for i := 0; i < count; i++ {
		SyncPraise(uint(i+1), "song")
	}
	db.Model(&statements.Special{}).Count(&count)
	for i := 0; i < count; i++ {
		SyncPraise(uint(i+1), "special")
	}
}

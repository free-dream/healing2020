package models

import (
	"errors"
	"strconv"
	"time"

	//"fmt"

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
	db.Model(&statements.User{}).Select("avatar,nick_name").Where("id =?", userId).First(&user)

	resultResp.VodUser = user.NickName
	resultResp.VodAvatar = user.Avatar

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
		db.Model(&statements.User{}).Select("avatar,nick_name").Where("id = ?", songRows.UserId).First(&userRows)
		recordResp[i].User = userRows.NickName
		recordResp[i].SongAvatar = userRows.Avatar

        recordResp[i].IsPraise,_ = HasPraise(2,myID,songRows.ID)
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
	status := 0

	tx := db.Begin()
	var vod statements.Vod
	result1 := tx.Model(&statements.Vod{}).Where("ID=?", vodId).First(&vod)
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

	err := tx.Model(&statements.Song{}).Create(&song).Error
	if err != nil {
		if status < 5 {
			status++
			tx.Rollback()
		} else {
			return "", err
		}
	}
    
    FinishTask("3",userId)

	return song.Name, tx.Commit().Error
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

	var praise statements.Praise
	praise.UserId = userid
	praise.Type = typesInt
	praise.PraiseId = id

	err := db.Model(&statements.Praise{}).Create(&praise).Error

    if SyncLock(userid) {
        FinishTask("5",userid)
    }

	return err, praiseData
}

func SyncLock(userid uint) bool {
    client := setting.RedisConn()
    keyname := "healing2020:PraiseSign:"+strconv.Itoa(int(userid))
    sign := client.Get(keyname)
    if sign == nil {
        client.Set(keyname,"1",0)
        return false
    }
    if sign.String() == "1" {
        client.Set(keyname,"2",0)
        return false
    }
    if sign.String() == "2" {
        client.Set(keyname,"3",0)
        return false
    }
    if sign.String() == "3" {
        client.Del(keyname)
        return true
    }
    return false
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

    FinishTask("2",uid)

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

package models

import (
	"errors"
	"strconv"
	"time"

	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"healing2020/pkg/tools"

	"github.com/jinzhu/gorm"
)

func GetPhone(info tools.RedisUser) string {
	return info.Phone
}

type RecordResp struct {
	Praise     int    `json:"praise"`
	User       string `json:"user"`
	Source     string `json:"source"`
	SongAvatar string `json:"songAvatar"`
}
type ResultResp struct {
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

func GetRecord(id string) ResultResp {
	intId, _ := strconv.Atoi(id)
	songId := uint(intId)
	var song statements.Song

	db := setting.MysqlConn()

	var resultResp ResultResp

	result := db.Model(&statements.Song{}).Select("vod_id,source,praise,style,language,name,user_id").Where("id=?", songId).First(&song)
	resultResp.Err = result.Error
	if result.Error != nil {
		return resultResp
	}
	vodId := song.VodId

	var vod statements.Vod
	db.Model(&statements.Vod{}).Select("created_at").Where("id=?", vodId).First(&vod)

	resultResp.Time = vod.CreatedAt
	resultResp.Name = song.Name
	resultResp.Style = song.Style
	resultResp.Language = song.Language

	userId := vod.UserId
	var user statements.User
	db.Model(&statements.User{}).Select("avatar,nick_name").Where("id =?", userId).First(&user)

	resultResp.VodUser = user.NickName
	resultResp.VodAvatar = user.Avatar

	recordsToVod := db.Model(&statements.Song{}).Where("vod_id = ?", vodId).Find(&song)
	count := recordsToVod.RowsAffected
	var recordResp []RecordResp = make([]RecordResp, count)

	rows, _ := recordsToVod.Rows()

	i := 0
	for rows.Next() {
		db.ScanRows(rows, &song)
		recordResp[i].Praise = song.Praise
		recordResp[i].Source = song.Source

		db.Model(&statements.User{}).Select("avatar,nick_name").Where("id = ?", song.UserId).First(&user)
		recordResp[i].User = user.NickName
		recordResp[i].SongAvatar = user.Avatar

		i++
	}
	resultResp.AllSongs = recordResp

	return resultResp
}

func CreateRecord(id string, source string, uid uint) (string, error) {
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

	err := tx.Model(&statements.Song{}).Create(&song).Error
	if err != nil {
		if status < 5 {
			status++
			tx.Rollback()
		} else {
			return "", err
		}
	}
	return song.Name, tx.Commit().Error
}

func HasPraise(types int,userid uint,id uint) (bool,uint) {
    db := setting.MysqlConn()

    var praise statements.Praise
    result := db.Model(&statements.Praise{}).Where("is_cancel = 0 and praise_id = ? and user_id = ? and type = ?",id,userid,types).First(&praise)

    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return false,0
    }
    return true,praise.ID
}

func IsPraiseCancel(types int,userid,uint,id uint) bool {
    db := setting.MysqlConn()

    var praise statements.Praise
    result := db.Model(&statements.Praise{}).Where("is_cancel = 1 and praise_id = ? and user_id = ? and type = ?",id,userid,types).First(&praise)

    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return false
    }
    return true
}

func CancelPraise(userid uint,strId string,types string) error {
	intId, _ := strconv.Atoi(strId)
	id := uint(intId)
    typesInt, _ := strconv.Atoi(types)
	db := setting.MysqlConn()

    hasPraise,praiseId := HasPraise(typesInt,userid,id)
    if !hasPraise {
        return errors.New("item does not be praised")
    }

	tx := db.Begin()
	status := 0
	if types == "1" {
		var song statements.Song
		result := tx.Model(&statements.Song{}).Where("ID=?", id).First(&song)
		if result.Error != nil {
			return result.Error
		}
		song.Praise = song.Praise - 1
		err := tx.Save(&song).Error
		if err != nil {
			if status < 5 {
				status++
				tx.Rollback()
			} else {
				return err
			}
		}
	}
	if types == "2" {
		var deliver statements.Deliver
		result := tx.Model(&statements.Deliver{}).Where("ID=?", id).First(&deliver)
		if result.Error != nil {
			return result.Error
		}
		deliver.Praise = deliver.Praise - 1
		err := tx.Save(&deliver).Error
		if err != nil {
			if status < 5 {
				status++
				tx.Rollback()
			} else {
				return err
            }
        }
    }

    var praise statements.Praise
    tx.Model(&statements.Praise{}).Where("id = ?",praiseId).First(&praise)
    praise.IsCancel = 1
    err := tx.Save(&praise).Error
    if err != nil {
        if status < 10 {
            status++
            tx.Rollback()
        } else {
            return err
        }
    }

	return tx.Commit().Error
}

func AddPraise(userid uint,strId string, types string) error {
	intId, _ := strconv.Atoi(strId)
	id := uint(intId)
    typesInt, _ := strconv.Atoi(types)
	db := setting.MysqlConn()

    hasPraise,_ := HasPraise(typesInt,userid,id)
    if hasPraise {
        return errors.New("can not praise repeatedly")
    }

	tx := db.Begin()
	status := 0
	if types == "1" {
		var song statements.Song
		result := tx.Model(&statements.Song{}).Where("ID=?", id).First(&song)
		if result.Error != nil {
			return result.Error
		}
		song.Praise = song.Praise + 1
		err := tx.Save(&song).Error
		if err != nil {
			if status < 5 {
				status++
				tx.Rollback()
			} else {
				return err
			}
		}
	}
	if types == "2" {
		var deliver statements.Deliver
		result := tx.Model(&statements.Deliver{}).Where("ID=?", id).First(&deliver)
		if result.Error != nil {
			return result.Error
		}
		deliver.Praise = deliver.Praise + 1
		err := tx.Save(&deliver).Error
		if err != nil {
			if status < 5 {
				status++
				tx.Rollback()
			} else {
				return err
			}
		}
	}

    var praise statements.Praise
    praise.UserId = userid
    praise.Type = typesInt
    praise.PraiseId = id
    tx.Model(&statements.Praise{}).Create(&praise)

	return tx.Commit().Error
}

func CreateVod(uid uint, singer string, style string, language string, name string, more string) error {
	db := setting.MysqlConn()

	var vod statements.Vod
	vod.UserId = uid
	vod.More = more
	vod.Name = name
	vod.Singer = singer
	vod.Style = style
	vod.Language = language
	err := db.Model(&statements.Vod{}).Create(&vod).Error
	return err
}

//根据vod_id获取user_id
func SelectUserIDByVodID(vod_id uint) (uint, error) {
	db := setting.MysqlConn()

	var vod statements.Vod
	err := db.Select("user_id").Where("id=?", vod_id).First(&vod).Error
	return vod.UserId, err
}

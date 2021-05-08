package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"log"
	"strconv"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type AllSpecial struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Intro     string `json:"intro"`
	MesNumber int    `json:"mesNumber"`
	SingSong  []UserMessage
}

type UserMessage struct {
	UserID   uint   `json:"user_id"`
	Nickname string `json:"user"`
	Avatar   string `json:"avatar"`

	Id        int       `json:"song_id"`
	CreatedAt time.Time `json:"time"`
	Praise    string    `json:"praise"`
	Song      string    `json:"name"`
	Record    string    `json:"record"`
    IsPraise  bool      `json:"isPraise"`

	//无用数据
	Type   string `json:"type"`
	Photo  string `json:"photo"`
	Text   string `json:"text"`
	Source string `json:"source"`
}

func SingHome(subjectID uint) (AllSpecial, error) {
	//连接mysql
	db := setting.MysqlConn()

	//获取歌房信息
	var singSubject statements.Subject
	err := db.Table("subject").Select("id, name, intro, photo").Where("id = ? ", subjectID).First(&singSubject).Error
	var count int
	err = db.Table("comment").Where("type = 1 and song_id = ?", subjectID).Count(&count).Error
	//获取歌曲信息
	var SingHome []UserMessage
	err = db.Table("special").Select("user_id, id, created_at, praise, song, record").Where("subject_id = ? ", subjectID).Scan(&SingHome).Error
	//获取用户信息
	UserElse := make([]statements.User, len(SingHome))
	for i := 0; i < len(SingHome); i++ {
		err = db.Table("user").Select("nick_name, avatar").Where("id = ?", SingHome[i].UserID).Scan(&UserElse[i]).Error
	}

	responseSing := make([]UserMessage, len(SingHome))
	for i := 0; i < len(SingHome); i++ {
		responseSing[i] = UserMessage{
			Nickname:  UserElse[i].NickName,
			Avatar:    UserElse[i].Avatar,
			UserID:    SingHome[i].UserID,
			Id:        SingHome[i].Id,
			CreatedAt: SingHome[i].CreatedAt,
			Praise:    SingHome[i].Praise,
			Song:      SingHome[i].Song,
			Record:    SingHome[i].Record,
		}
        responseSing[i].IsPraise,_ = HasPraise(3,SingHome[i].UserID,uint(SingHome[i].Id))
	}

	allSpecial := AllSpecial{
		ID:        singSubject.ID,
		Name:      singSubject.Name,
		Intro:     singSubject.Intro,
		MesNumber: count,
		SingSong:  responseSing,
	}

	return allSpecial, err
}

//发送歌房歌曲数据
func PostSpecial(Subject_id string, Song string, User_id string, Record string) error {
	intId, _ := strconv.Atoi(Subject_id)
	subject_id := uint(intId)

	int2Id, _ := strconv.Atoi(User_id)
	user_id := uint(int2Id)

	db := setting.MysqlConn()

	status := 0
	tx := db.Begin()
	if Song != "" {
		var dev statements.Special
		var userother statements.UserOther

		dev.SubjectId = subject_id
		dev.Song = Song
		dev.UserId = user_id
		dev.Record = Record
		//发送歌曲
		err := tx.Model(&statements.Special{}).Create(&dev).Error
		if err != nil {
			if status < 5 {
				status++
				tx.Rollback()
			} else {
				return err
			}
		}

		result := tx.Model(&statements.UserOther{}).Where("user_id= ?", user_id).First(&userother)
		if result.Error != nil {
			log.Println(0)
			return result.Error
		}
		//判断完成每日任务和增加积分
		if userother.Lo4 != 1 {
			err2 := tx.Model(&statements.UserOther{}).Where("user_id = ?", user_id).Update("lo4", 1).Error
			if err2 != nil {
				log.Println(1)
				if status < 5 {
					status++
					tx.Rollback()
				} else {
					return err2
				}
			}
			var user statements.User
			result := tx.Model(&statements.User{}).Where("id= ?", user_id).First(&user)
			if result.Error != nil {
				log.Println(2)
				return result.Error
			}
			if user.Money >= 0 {
				user.Money = user.Money + 20
				err3 := tx.Save(&user).Error
				if err3 != nil {
					log.Println(3)
					if status < 5 {
						status++
						tx.Rollback()
					} else {
						return err3
					}
				}
			}
		}
	}
	return tx.Commit().Error
}

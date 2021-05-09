package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"strconv"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type AllSpecial struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Intro     string `json:"intro"`
	MesNumber int    `json:"mesNumber"`
	// HotSong   []UserMessage
	SingSong []UserMessage
}

type UserMessage struct {
	UserID   uint   `json:"user_id"`
	Nickname string `json:"user"`
	Avatar   string `json:"avatar"`

	Id        int       `json:"song_id"`
	CreatedAt time.Time `json:"time"`
	Praise    int       `json:"praise"`
	Song      string    `json:"name"`
	Record    string    `json:"source"`
	IsPraise  bool      `json:"isPraise"`

	//无用数据
	Type  string `json:"type"`
	Photo string `json:"photo"`
	Text  string `json:"text"`
	NilId string `json:"id"`
}

func SingHome(subjectID uint) (AllSpecial, error) {
	//连接mysql
	db := setting.MysqlConn()

	//获取歌房信息
	var singSubject statements.Subject
	err := db.Table("subject").Select("id, name, intro, photo").Where("id = ? ", subjectID).First(&singSubject).Error
	var count int
	err = db.Table("comment").Where("type = 1 and song_id = ?", subjectID).Count(&count).Error

	// //获取热门歌曲信息
	// var Hot []UserMessage
	// err = db.Table("special").Select("user_id, id, created_at, praise, song, record").Where("subject_id = ? ", subjectID).Order("praise DESC").Limit(3).Scan(&Hot).Error
	// //获取热门用户信息
	// HotElse := make([]statements.User, len(Hot))
	// for i := 0; i < len(Hot); i++ {
	// 	err = db.Table("user").Select("nick_name, avatar").Where("id = ?", Hot[i].UserID).Scan(&HotElse[i]).Error
	// }
	// responseHot := make([]UserMessage, len(Hot))
	// for i := 0; i < len(Hot); i++ {
	// 	if Hot[i].Praise >= 5 {
	// 		responseHot[i] = UserMessage{
	// 			Nickname:  HotElse[i].NickName,
	// 			Avatar:    HotElse[i].Avatar,
	// 			UserID:    Hot[i].UserID,
	// 			Id:        Hot[i].Id,
	// 			CreatedAt: Hot[i].CreatedAt,
	// 			Praise:    Hot[i].Praise,
	// 			Song:      Hot[i].Song,
	// 			Record:    Hot[i].Record,
	// 		}
	// 		responseHot[i].IsPraise, _ = HasPraise(3, Hot[i].UserID, uint(Hot[i].Id))
	// 	}
	// }

	//获取列表歌曲信息
	var SingHome []UserMessage
	err = db.Table("special").Select("user_id, id, created_at, praise, song, record").Where("subject_id = ? ", subjectID).Order("created_at DESC").Scan(&SingHome).Error
	//获取列表用户信息
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
		responseSing[i].IsPraise, _ = HasPraise(3, SingHome[i].UserID, uint(SingHome[i].Id))
	}

	allSpecial := AllSpecial{
		ID:        singSubject.ID,
		Name:      singSubject.Name,
		Intro:     singSubject.Intro,
		MesNumber: count,
		// HotSong:   responseHot,
		SingSong: responseSing,
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

	tx := db.Begin()
	if Song != "" {
		var dev statements.Special
		// var userother statements.UserOther

		dev.SubjectId = subject_id
		dev.Song = Song
		dev.UserId = user_id
		dev.Record = Record
		//发送歌曲
		err := tx.Model(&statements.Special{}).Create(&dev).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		result := FinishTask("4", user_id)
		if result != nil {
			tx.Rollback()
			return result
		}
	}
	return tx.Commit().Error
}

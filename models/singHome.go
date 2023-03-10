package models

import (
	"errors"
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"healing2020/pkg/tools"
	"strconv"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type AllSpecial struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Intro     string `json:"intro"`
	MesNumber int    `json:"mesNumber"`
	HotNumber int    `json:"hotNumber"`
	HotSong   []UserMessage
	SingSong  []UserMessage
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

func SingHome(belong string, pageStr string, pageStr2 string, subjectID uint, User_id string) (AllSpecial, error) {
	page, _ := strconv.Atoi(pageStr)
	page2, _ := strconv.Atoi(pageStr2)
	belongs, _ := strconv.Atoi(belong)
	int2Id, _ := strconv.Atoi(User_id)
	UserID := uint(int2Id)
	//连接mysql
	db := setting.MysqlConn()

	//获取歌房信息
	var singSubject statements.Subject
	err := db.Table("subject").Select("id, name, intro, photo").Where("id = ? ", subjectID).First(&singSubject).Error
	var count int
	err = db.Table("comment").Where("type = 1 and song_id = ?", subjectID).Count(&count).Error
	var count2 int

	var Hot []UserMessage
	var responseHot []UserMessage

	//歌房进入
	if belongs == 1 {
		//获取热门歌曲信息
		err = db.Table("special").Select("user_id, id, created_at, praise, song, record").Where("subject_id = ? ", subjectID).Order("praise DESC").Limit(3).Scan(&Hot).Error
		//获取热门用户信息
		HotElse := make([]statements.User, len(Hot))
		for i := 0; i < len(Hot); i++ {
			err = db.Table("user").Select("nick_name, avatar, setting1, sex").Where("id = ?", Hot[i].UserID).Scan(&HotElse[i]).Error
		}
		responseHot = make([]UserMessage, len(Hot))
		for i := 0; i < len(Hot); i++ {
			if Hot[i].Praise >= 5 {
				count2 = 1
				count2 = count2 + 1
				if HotElse[i].Setting1 == 0 {
					responseHot[i] = UserMessage{
						Nickname:  HotElse[i].NickName,
						Avatar:    tools.GetAvatarUrl(HotElse[i].Sex),
						UserID:    Hot[i].UserID,
						Id:        Hot[i].Id,
						CreatedAt: Hot[i].CreatedAt,
						Praise:    Hot[i].Praise,
						Song:      Hot[i].Song,
						Record:    Hot[i].Record,
					}
				} else {
					responseHot[i] = UserMessage{
						Nickname:  HotElse[i].NickName,
						Avatar:    HotElse[i].Avatar,
						UserID:    Hot[i].UserID,
						Id:        Hot[i].Id,
						CreatedAt: Hot[i].CreatedAt,
						Praise:    Hot[i].Praise,
						Song:      Hot[i].Song,
						Record:    Hot[i].Record,
					}
				}
				responseHot[i].IsPraise, _ = HasPraise(3, Hot[i].UserID, uint(Hot[i].Id))
			}
		}
	}

	//个人页进入
	if belongs == 2 {
		//获取个人歌曲信息
		err = db.Table("special").Select("user_id, id, created_at, praise, song, record").Where("subject_id = ? and user_id = ?", subjectID, UserID).Order("praise DESC").Count(&count2).Scan(&Hot).Error
		//获取个人用户信息
		HotElse := make([]statements.User, len(Hot))
		for i := 0; i < len(Hot); i++ {
			err = db.Table("user").Select("nick_name, avatar, setting1, sex").Where("id = ?", UserID).Scan(&HotElse[i]).Error
		}
		responseHot = make([]UserMessage, len(Hot))
		for i := 0; i < len(Hot); i++ {
			if HotElse[i].Setting1 == 0 {
				responseHot[i] = UserMessage{
					Nickname:  HotElse[i].NickName,
					Avatar:    tools.GetAvatarUrl(HotElse[i].Sex),
					UserID:    Hot[i].UserID,
					Id:        Hot[i].Id,
					CreatedAt: Hot[i].CreatedAt,
					Praise:    Hot[i].Praise,
					Song:      Hot[i].Song,
					Record:    Hot[i].Record,
				}
			} else {
				responseHot[i] = UserMessage{
					Nickname:  HotElse[i].NickName,
					Avatar:    HotElse[i].Avatar,
					UserID:    Hot[i].UserID,
					Id:        Hot[i].Id,
					CreatedAt: Hot[i].CreatedAt,
					Praise:    Hot[i].Praise,
					Song:      Hot[i].Song,
					Record:    Hot[i].Record,
				}
			}
			responseHot[i].IsPraise, _ = HasPraise(3, Hot[i].UserID, uint(Hot[i].Id))
		}
	}

	//获取列表歌曲信息
	var SingHome []UserMessage
	err = db.Table("special").Select("user_id, id, created_at, praise, song, record").Where("subject_id = ? ", subjectID).Order("created_at DESC").Scan(&SingHome).Error
	//获取列表用户信息
	UserElse := make([]statements.User, len(SingHome))
	for i := 0; i < len(SingHome); i++ {
		err = db.Table("user").Select("nick_name, avatar, setting1, sex").Where("id = ?", SingHome[i].UserID).Scan(&UserElse[i]).Error
	}
	responseSing := make([]UserMessage, len(SingHome))
	for i := 0; i < len(SingHome); i++ {
		if UserElse[i].Setting1 == 0 {
			responseSing[i] = UserMessage{
				Nickname:  UserElse[i].NickName,
				Avatar:    tools.GetAvatarUrl(UserElse[i].Sex),
				UserID:    SingHome[i].UserID,
				Id:        SingHome[i].Id,
				CreatedAt: SingHome[i].CreatedAt,
				Praise:    SingHome[i].Praise,
				Song:      SingHome[i].Song,
				Record:    SingHome[i].Record,
			}
		} else {
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
		}
		responseSing[i].IsPraise, _ = HasPraise(3, SingHome[i].UserID, uint(SingHome[i].Id))
	}
	var allSpecial AllSpecial
	pageResponSing, err := PageinSong(page, responseSing)
	pageResponHot, err := PageinHot(page2, responseHot)
	allSpecial = AllSpecial{
		ID:        singSubject.ID,
		Name:      singSubject.Name,
		Intro:     singSubject.Intro,
		MesNumber: count,
		HotNumber: count2,
		HotSong:   pageResponHot,
		SingSong:  pageResponSing,
	}

	return allSpecial, err
}

func PageinSong(page int, data []UserMessage) ([]UserMessage, error) {
	if (page-1)*20 > len(data) {
		return nil, errors.New("page out of range")
	}

	var result []UserMessage = make([]UserMessage, 20)
	for i := 0; i < 20; i++ {
		if (page-1)*20+i >= len(data) {
			break
		}
		result[i] = data[(page-1)*20+i]
	}
	return result, nil
}

func PageinHot(page int, data []UserMessage) ([]UserMessage, error) {
	if (page-1)*10 > len(data) {
		return nil, errors.New("page out of range")
	}

	var result []UserMessage = make([]UserMessage, 10)
	for i := 0; i < 10; i++ {
		if (page-1)*10+i >= len(data) {
			break
		}
		result[i] = data[(page-1)*10+i]
	}
	return result, nil
}

//发送歌房歌曲数据
func PostSpecial(Subject_id string, Song string, user_id uint, Record string) error {
	intId, _ := strconv.Atoi(Subject_id)
	subject_id := uint(intId)

	// int2Id, _ := strconv.Atoi(User_id)
	// user_id := uint(int2Id)

	db := setting.MysqlConn()

	tx := db.Begin()
	if Song != "" {
		var dev statements.Special
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

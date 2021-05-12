package models

import (
	"strings"
	"time"

	"healing2020/models/statements"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type RequestSongs struct { //点歌
	ID        uint      `json:"id"`
	Name      string    `json:"song"`
	CreatedAt time.Time `json:"time"`
	HideName  int       `json:"hidename"`
}
type Songs struct { //唱歌
	ID        uint      `json:"id"`
	Name      string    `json:"song"`
	CreatedAt time.Time `json:"time"`
	From      string    `json:"from"`
	Praise    int       `json:"praise"`
	IsPraise  bool      `json:"ispraise"`
	IsHide    int       `json:"ishide"`
}
type Admire struct { //点赞
	ID        uint      `json:"id"`
	Name      string    `json:"song"`
	CreatedAt time.Time `json:"time"`
	From      string    `json:"from"`
	Praise    int       `json:"number"`
}

//select并返回用户个人背景信息和剩余匿名次数
func ResponseUserOther(userID uint) (statements.UserOther, error) {
	//连接mysql
	db := setting.MysqlConn()

	//查询
	var userOther statements.UserOther
	err := db.Select("ava_background, now, remain_hide_name").Where("user_id=?", userID).First(&userOther).Error
	return userOther, err
}

//select并返回点歌信息
func ResponseVod(userID uint, my_others string) ([]RequestSongs, error) {
	//连接mysql
	db := setting.MysqlConn()
	var allVod []RequestSongs
	var err error
	//获取点歌信息
	if my_others == "my" {
		err = db.Table("vod").Select("id, name, created_at, hide_name").Where("user_id = ?", userID).Scan(&allVod).Error
	} else if my_others == "others" {
		err = db.Table("vod").Select("id, name, created_at, hide_name").Where("user_id = ? And hide_name = 0", userID).Scan(&allVod).Error
	}
	return allVod, err
}

//ResponseSongs使用
//将select到的song[]信息代入一个[]Songs结构
func songToSongs(song []statements.Song) []Songs {
	var s []Songs
	for _, value := range song {
		s = append(s, Songs{
			ID:        value.VodId,
			Name:      value.Name,
			CreatedAt: value.CreatedAt,
			From:      "治愈",
			Praise:    GetPraiseCount("song", value.ID),
			IsHide:    value.IsHide,
		})
	}
	return s
}

//ResponseSongs使用
//将select到的deliver[]信息代入一个[]Songs结构
func deliverToSongs(deliver []statements.Deliver) []Songs {
	var s []Songs
	//截取deliver的textfield
	for _, value := range deliver {
		splitDeliver := strings.Split(value.TextField, "")
		if len(splitDeliver) > 5 {
			value.TextField = strings.Join(splitDeliver[:5], "") + "..."
		}
		s = append(s, Songs{
			ID:        value.ID,
			Name:      value.TextField,
			CreatedAt: value.CreatedAt,
			From:      "投递箱",
			Praise:    GetPraiseCount("deliver", value.ID),
		})
	}
	return s
}

//ResponseSongs使用
//将select到的special[]信息代入一个[]Songs结构
func specialToSongs(special []statements.Special) []Songs {
	var s []Songs
	for _, value := range special {
		s = append(s, Songs{
			ID:        value.SubjectId,
			Name:      value.Song,
			CreatedAt: value.CreatedAt,
			From:      "歌房",
			Praise:    GetPraiseCount("special", value.ID),
		})
	}
	return s
}

//select并返回用户唱歌信息
func ResponseSongs(userID uint, myID uint, my_others string) ([]Songs, error) {
	var err error

	//连接mysql
	db := setting.MysqlConn()

	var song []statements.Song
	//获取唱歌信息
	if my_others == "my" {
		err = db.Select("id, vod_id, name, created_at, is_hide").Where("user_id=?", userID).Find(&song).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	} else if my_others == "others" {
		err = db.Select("id, vod_id, name, created_at, is_hide").Where("user_id=? AND is_hide = 0", userID).Find(&song).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}

	//获取歌房专题歌曲信息
	var special []statements.Special
	err = db.Select("id, subject_id, song, created_at").Where("user_id=?", userID).Find(&special).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	//获得投递箱信息
	var deliver []statements.Deliver
	err = db.Select("id, text_field, created_at").Where("user_id=?", userID).Find(&deliver).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	//处理不同表select下来的信息, 转换为Songs类型
	songSongs := songToSongs(song)
	specialSongs := specialToSongs(special)
	deliverSongs := deliverToSongs(deliver)
	//合并数据
	allSongs := append(append(songSongs, specialSongs...), deliverSongs...)
	for key, value := range allSongs {
		switch value.From {
		case "投递箱":
			allSongs[key].IsPraise, _ = HasPraise(1, myID, value.ID)
		case "治愈":
			allSongs[key].IsPraise, _ = HasPraise(2, myID, value.ID)
		case "歌房":
			allSongs[key].IsPraise, _ = HasPraise(3, myID, value.ID)
		}
	}
	return allSongs, nil
}

//ResponsePraise使用
//将select到的[]deliver信息代入到一个[]Admire结构
func deliverToAdmire(deliver []statements.Deliver) []Admire {
	var admire []Admire
	for _, value := range deliver {
		admire = append(admire, Admire{
			ID:        value.ID,
			Name:      value.TextField,
			CreatedAt: value.CreatedAt,
			From:      "投递箱",
			Praise:    GetPraiseCount("deliver", value.ID),
		})
	}
	return admire
}

//ResponsePraise使用
//将select到的[]special信息代入到一个[]Admire结构
func specialToAdmire(special []statements.Special) []Admire {
	var admire []Admire
	for _, value := range special {
		admire = append(admire, Admire{
			ID:        value.SubjectId,
			Name:      value.Song,
			CreatedAt: value.CreatedAt,
			From:      "歌房",
			Praise:    GetPraiseCount("special", value.ID),
		})
	}
	return admire
}

//ResponsePraise使用
//将select到的[]song信息代入到一个[]Admire结构
func songToAdmire(song []statements.Song) []Admire {
	var admire []Admire
	for _, value := range song {
		admire = append(admire, Admire{
			ID:        value.VodId,
			Name:      value.Name,
			CreatedAt: value.CreatedAt,
			From:      "治愈",
			Praise:    GetPraiseCount("song", value.ID),
		})
	}
	return admire
}

//ResponsePraise使用
//根据type提取结构体id字段
func getPraiseStructID(praise []statements.Praise) map[string][]uint {
	var deliverID []uint
	var healID []uint
	var specialID []uint
	for _, value := range praise {
		switch value.Type {
		case 1:
			deliverID = append(deliverID, value.PraiseId)
		case 2:
			healID = append(healID, value.PraiseId)
		case 3:
			specialID = append(specialID, value.PraiseId)
		}
	}
	return map[string][]uint{
		"deliver": deliverID,
		"heal":    healID,
		"special": specialID,
	}
}

//select并返回用户点赞信息
func ResponsePraise(userID uint) ([]Admire, error) {
	var err error

	//连接mysql
	db := setting.MysqlConn()

	//获取点赞对应条目
	var praise []statements.Praise
	err = db.Select("type, praise_id").Where("user_id = ? AND is_cancel = 0", userID).Find(&praise).Error
	if err != nil {
		return nil, err
	}
	allID := getPraiseStructID(praise)

	//投递箱,type=1
	var deliverInf []statements.Deliver
	err = db.Select("id, text_field, created_at").Where("id in (?)", allID["deliver"]).Find(&deliverInf).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	//治愈,type=2
	var songInf []statements.Song
	err = db.Select("id, vod_id, name, created_at").Where("id in (?)", allID["heal"]).Find(&songInf).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	//专题歌曲,type=3
	//返回id为歌房id
	var specialInf []statements.Special
	err = db.Select("id, subject_id, song, created_at").Where("id in (?)", allID["special"]).Find(&specialInf).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	allPraise := append(append(deliverToAdmire(deliverInf), songToAdmire(songInf)...), specialToAdmire(specialInf)...)
	return allPraise, err
}

//给对应点歌匿名
func HideName(vodID uint, userID uint) error {
	db := setting.MysqlConn()

	tx := db.Begin()
	err := tx.Model(&statements.Vod{}).Where("id = ?", vodID).Update(statements.Vod{HideName: 1}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(&statements.UserOther{}).Where("user_id = ?", userID).Update("remain_hide_name", gorm.Expr("remain_hide_name - ?", 1)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

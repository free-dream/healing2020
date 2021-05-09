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
	ID        uint      `json:"id" gorm:"column:vod_id"`
	Name      string    `json:"song"`
	CreatedAt time.Time `json:"time"`
	From      string    `json:"from"`
	Praise    int       `json:"praise"`
	IsPraise  bool      `json:"ispraise"`
}
type Admire struct { //点赞
	ID        uint      `json:"id"`
	Name      string    `json:"song"`
	CreatedAt time.Time `json:"time"`
	From      string    `json:"from"`
	Praise    int       `json:"number"`
}

//获取其它用户信息接口用
//select并根据id返回用户信息
func ResponseUser(userID uint) (statements.User, error) {
	//连接mysql
	db := setting.MysqlConn()

	var user statements.User
	err := db.Where("id=?", userID).First(&user).Error
	return user, err
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
func ResponseVod(userID uint) ([]RequestSongs, error) {
	//连接mysql
	db := setting.MysqlConn()

	//获取点歌信息
	var allVod []RequestSongs
	err := db.Table("vod").Select("id, name, created_at, hide_name").Where("user_id=?", userID).Scan(&allVod).Error
	return allVod, err
}

//ResponseSongs使用
//对deliver的返回进行处理，将deliver的textfield截至5个字
func handleDeliver(deliver []statements.Deliver) []statements.Deliver {
	for key := range deliver {
		splitDeliver := strings.Split(deliver[key].TextField, "")
		if len(splitDeliver) <= 5 {
			continue
		}
		deliver[key].TextField = strings.Join(splitDeliver[:5], "")
	}
	return deliver
}

//ResponseSongs使用
//将select到的deliver[]信息代入一个[]Songs结构
func deliverToSongs(deliver []statements.Deliver) []Songs {
	var s []Songs
	for _, value := range deliver {
		s = append(s, Songs{
			ID:        value.ID,
			Name:      value.TextField,
			CreatedAt: value.CreatedAt,
			From:      "投递箱",
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
		})
	}
	return s
}

//select并返回用户唱歌信息
func ResponseSongs(userID uint, myID uint) ([]Songs, error) {
	var err error

	//连接mysql
	db := setting.MysqlConn()

	//获取唱歌信息
	var singSongs []Songs
	err = db.Table("song").Select("vod_id, name, created_at").Where("user_id=?", userID).Scan(&singSongs).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	//获取歌房专题歌曲信息
	var special []statements.Special
	err = db.Select("subject_id, song, created_at").Where("user_id=?", userID).Find(&special).Error
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
	for key := range singSongs {
		singSongs[key].From = "治愈"
	}
	specialSongs := specialToSongs(special)
	deliverSongs := deliverToSongs(handleDeliver(deliver))
	//合并数据
	allSongs := append(append(singSongs, specialSongs...), deliverSongs...)
	for key, value := range allSongs {
		switch value.From {
		case "投递箱":
			allSongs[key].IsPraise, _ = HasPraise(1, myID, value.ID)
			allSongs[key].Praise = GetPraiseCount("deliver", value.ID)
		case "治愈":
			allSongs[key].IsPraise, _ = HasPraise(2, myID, value.ID)
			allSongs[key].Praise = GetPraiseCount("song", value.ID)
		case "歌房":
			allSongs[key].IsPraise, _ = HasPraise(3, myID, value.ID)
			allSongs[key].Praise = GetPraiseCount("special", value.ID)
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
			Praise:    value.Praise,
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
			Praise:    value.Praise,
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
			Praise:    value.Praise,
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
	err = db.Select("id, text_field, created_at, praise").Where("id in (?)", allID["deliver"]).Find(&deliverInf).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	//治愈,type=2
	var songInf []statements.Song
	err = db.Table("song").Select("vod_id, name, created_at, praise").Where("id in (?)", allID["heal"]).Scan(&songInf).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	//专题歌曲,type=3
	//返回id为歌房id
	var specialInf []statements.Special
	err = db.Select("subject_id, song, created_at, praise").Where("id in (?)", allID["special"]).Find(&specialInf).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	allPraise := append(append(deliverToAdmire(deliverInf), songToAdmire(songInf)...), specialToAdmire(specialInf)...)
	return allPraise, err
}

//给对应点歌匿名
func HideName(vodID uint, userID uint) error {
	db := setting.MysqlConn()

	//开始更新事务
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

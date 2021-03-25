package models

import (
	"strings"
	"time"

	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"healing2020/pkg/tools"

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
}
type Admire struct { //点赞
	ID        uint      `json:"id"`
	Name      string    `json:"song"`
	CreatedAt time.Time `json:"time"`
	From      string    `json:"from"`
	Praise    int       `json:"number"`
}

//ResponseSongs使用
//给Songs加上From，如歌房、投递、治愈
func songsFrom(someSongs []Songs, from string) []Songs {
	for i := 0; i < len(someSongs); i++ {
		someSongs[i].From = from
	}
	return someSongs
}

//ResponseSongs使用
//对deliver的返回进行处理，将deliver的textfield截至5个字
func handleDeliver(someSongs []Songs) []Songs {
	splitDeliver := make([]string, 6)
	for i := 0; i < len(someSongs); i++ {
		splitDeliver = strings.SplitN(someSongs[i].Name, "", 6)
		someSongs[i].Name = strings.Join(splitDeliver[:5], "")
	}
	return someSongs
}

//ResponseSongs使用
//将select到的deliver[]信息代入一个[]Songs结构
func deliverToSongs(deliver []statements.Deliver) []Songs {
	s := make([]Songs, len(deliver))
	for i := 0; i < len(deliver); i++ {
		s[i] = Songs{
			ID:        deliver[i].ID,
			Name:      deliver[i].TextField,
			CreatedAt: deliver[i].CreatedAt,
			From:      "投递箱",
		}
	}
	return s
}

//ResponsePraise使用
//根据praiseID来select对应表的信息，并加上From来源(歌房、治愈、投递箱）
func selectPraiseInf(db *gorm.DB, table string, from string, praiseID uint) (Admire, error) {
	var admireInf Admire
	err := db.Table(table).Select("id, name, created_at, praise").Where("id=?", praiseID).Scan(&admireInf).Error
	admireInf.From = from
	return admireInf, err
}

//ResponsePraise使用
//将select到的deliver信息代入到一个Admire结构
func deliverToAdmire(deliver statements.Deliver) Admire {
	a := Admire{
		ID:        deliver.ID,
		Name:      deliver.TextField,
		CreatedAt: deliver.CreatedAt,
		From:      "投递箱",
		Praise:    deliver.Praise,
	}
	return a
}

//获取其它用户信息接口用
//select并根据id返回用户信息
func ResponseUser(userID uint) (statements.User, error) {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	var user statements.User
	err := db.Where("id=?", userID).First(&user).Error
	return user, err
}

//select并返回用户现在使用的个人背景
func ResponseUserOther(userID uint) (string, error) {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	//查询
	var nowUserOther statements.UserOther
	err := db.Select("now").Where("user_id=?", userID).First(&nowUserOther).Error
	return tools.GetBackgroundUrl(nowUserOther.Now), err
}

//select并返回点歌信息
func ResponseVod(userID uint) ([]RequestSongs, error) {
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	//获取点歌信息
	var allVod []RequestSongs
	err := db.Table("vod").Select("id, name, created_at, hide_name").Where("user_id=?", userID).Scan(&allVod).Error
	return allVod, err
}

//select并返回用户唱歌信息
func ResponseSongs(userID uint) ([]Songs, error) {
	var err error

	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	//获取唱歌信息
	var singSongs []Songs
	err = db.Table("song").Select("id, name, created_at").Where("user_id=?", userID).Scan(&singSongs).Error
	if err != nil {
		return nil, err
	}

	//获取歌房专题歌曲信息
	var specialSongs []Songs
	err = db.Table("special").Select("id, name, created_at").Where("user_id=?", userID).Scan(&specialSongs).Error
	if err != nil {
		return nil, err
	}

	//获得投递箱信息
	var deliver []statements.Deliver
	err = db.Select("id, text_field, created_at").Where("user_id=?", userID).Find(&deliver).Error
	if err != nil {
		return nil, err
	}

	//处理不同表select下来的信息, 转换为Songs类型
	singSongs = songsFrom(singSongs, "治愈")

	specialSongs = songsFrom(specialSongs, "歌房")

	deliverSongs := deliverToSongs(deliver)
	deliverSongs = handleDeliver(deliverSongs)

	//合并数据
	allSongs := append(append(singSongs, specialSongs...), deliverSongs...)
	return allSongs, err
}

//select并返回用户点赞信息
func ResponsePraise(userID uint) ([]Admire, error) {
	var err error

	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	//获取点赞对应条目
	var praise []statements.Praise
	err = db.Select("type, praise_id").Where("user_id=?", userID).Find(&praise).Error
	if err != nil {
		return nil, err
	}
	//根据type查询不同表获得信息
	allPraise := make([]Admire, len(praise))
	for i := 0; i < len(praise); i++ {
		switch praise[i].Type {
		//投递箱
		case 1:
			var deliverInf statements.Deliver
			err = db.Select("id, text_field, created_at, praise").Where("id=?", praise[i].PraiseId).First(&deliverInf).Error
			allPraise[i] = deliverToAdmire(deliverInf)
			if err != nil {
				return nil, err
			}
		//治愈
		case 2:
			allPraise[i], err = selectPraiseInf(db, "song", "治愈", praise[i].PraiseId)
			if err != nil {
				return nil, err
			}
		//专题歌曲
		case 3:
			allPraise[i], err = selectPraiseInf(db, "special", "歌房", praise[i].PraiseId)
			if err != nil {
				return nil, err
			}
		}
	}
	return allPraise, err
}

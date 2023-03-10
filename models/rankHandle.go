package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"healing2020/pkg/tools"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"encoding/json"
	"errors"

	//"fmt"
	"strconv"
	"time"
)

type Rank struct {
	ID     uint   `json:"id"`
	User   string `json:"user"`
	Avatar string `json:"avatar"`

	Type   int    `json:"type"`
	Photo  string `json:"photo"`
	Text   string `json:"text"`
	Source string `json:"source"`

	Time     string `json:"time"`
	Praise   int    `json:"praise"`
	Name     string `json:"name"`
	VodId    uint   `json:"vodid"`
	UserId   uint   `json:"userid"`
	IsPraise bool   `json:"isPraise"`
}

type AllRank struct {
	Time string `json:"time"`
	Data []Rank `json:"data"`
}

func errRollBack(tx *gorm.DB, status *int) error {
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	if tx.Error != nil {
		if *status < 5 {
			*status++
			tx.Rollback()
		} else {
			return tx.Error
		}
	}
	return tx.Commit().Error
}

func monthTransfer(mon string) string {
	if mon == "March" {
		return "03"
	}
	if mon == "April" {
		return "04"
	}
	if mon == "May" {
		return "05"
	}
	return ""
}

func dayTransfer(day string) string {
	if len(day) == 1 {
		return "0" + day
	}
	return day
}

func SendDeliverRank() error {
	//everyday's work

	//get from mysql
	db := setting.MysqlConn()

	var deliver []statements.Deliver
	year, mon, day := time.Now().Date()
	date := strconv.Itoa(year) + "_" + monthTransfer(mon.String()) + "_" + dayTransfer(strconv.Itoa(day)) + "%"
	//fmt.Println(date)
	result := db.Model(&statements.Deliver{}).Where("created_at LIKE ?", date).Order("praise, created_at desc").Find(&deliver)
	rows := result.RowsAffected
	if rows == 0 {
		client := setting.RedisConn()

		count, _ := client.Get("healing2020:rankCount").Float64()
		keyName := "healing2020:Deliver." + strconv.FormatFloat(count/100+5.10, 'f', 2, 64)
		client.Set(keyName, "", 0)
		return errors.New("no data")
	}
	err1 := result.Error
	if err1 != nil {
		return err1
	}
	var rank []Rank = make([]Rank, 10)
	for i := 0; i < min(int(rows), 10); i++ {
		rank[i].ID = deliver[i].ID
		rank[i].Type = deliver[i].Type
		rank[i].Text = deliver[i].TextField
		rank[i].Photo = deliver[i].Photo
		rank[i].Source = deliver[i].Record
		rank[i].UserId = deliver[i].UserId

		userid := deliver[i].UserId
		var user statements.User
		result2 := db.Model(&statements.User{}).Where("id=?", userid).First(&user)
		err2 := result2.Error
		if err2 != nil {
			return err2
		}
		rank[i].User = user.NickName
		rank[i].Avatar = user.Avatar

		if user.Setting1 == 0 {
			rank[i].Avatar = tools.GetAvatarUrl(user.Sex)
		}
	}
	jsonRank, _ := json.Marshal(rank)

	//set in redis
	client := setting.RedisConn()
	count, _ := client.Get("healing2020:rankCount").Float64()
	keyName := "healing2020:Deliver." + strconv.FormatFloat(count/100+5.10, 'f', 2, 64)
	client.Set(keyName, jsonRank, 0)

	return nil
}

func GetDeliverRank(userid uint) ([]AllRank, string) {
	result := make([]AllRank, 10)
	client := setting.RedisConn()
	count, _ := client.Get("healing2020:rankCount").Float64()
	var i float64 = 0
	for j := 0; ; j++ {
		rank := make([]Rank, 10)
		//fmt.Println(count)
		//fmt.Println(i)
		if i*100 > count {
			break
		}
		var date float64 = 5.10 + i
		dateStr := strconv.FormatFloat(date, 'f', 2, 64)
		keyname := "healing2020:Deliver." + dateStr
		data, err := client.Get(keyname).Bytes()
		if err != nil {
			result[j].Data = rank
			result[j].Time = dateStr
			i = i + 0.01
			continue
			// return nil, err.Error()
		}
		json.Unmarshal(data, &rank)
		// ???????????????????????????
		for k := 0; k < len(rank); k++ {
			rank[k].IsPraise, _ = HasPraise(1, userid, rank[k].ID)
			rank[k].Praise = GetPraiseCount("deliver", rank[k].ID)
		}
		i = i + 0.01

		result[j].Data = rank
		result[j].Time = dateStr
	}

	//fmt.Println(result)
	return result, ""
}

func SendSongRank() error {
	//everyday's work

	//get from mysql
	db := setting.MysqlConn()

	var song []statements.Song
	year, mon, day := time.Now().Date()
	date := strconv.Itoa(year) + "_" + monthTransfer(mon.String()) + "_" + dayTransfer(strconv.Itoa(day))
	result := db.Model(&statements.Song{}).Where("is_hide = 0 and created_at LIKE ?", date+"%").Order("praise, created_at desc").Find(&song)
	rows := result.RowsAffected
	if rows == 0 {
		client := setting.RedisConn()

		count, _ := client.Get("healing2020:rankCount").Float64()
		keyName := "healing2020:Song." + strconv.FormatFloat(count/100+5.10, 'f', 2, 64)
		client.Set(keyName, "", 0)
		return errors.New("no data")
	}
	err1 := result.Error
	if err1 != nil {
		return err1
	}
	var rank []Rank = make([]Rank, 10)
	for i := 0; i < min(int(rows), 10); i++ {
		rank[i].ID = song[i].ID
		rank[i].Name = song[i].Name
		rank[i].Praise = GetPraiseCount("song", song[i].ID)
		rank[i].Time = date
		rank[i].VodId = song[i].VodId
		rank[i].Source = song[i].Source

		userid := song[i].UserId
		var user statements.User
		result2 := db.Model(&statements.User{}).Where("id=?", userid).First(&user)
		err2 := result2.Error
		if err2 != nil {
			return err2
		}
		rank[i].User = user.NickName
		rank[i].Avatar = user.Avatar
		rank[i].UserId = userid

		if user.Setting1 == 0 {
			rank[i].Avatar = tools.GetAvatarUrl(user.Sex)
		}
	}
	jsonRank, _ := json.Marshal(rank)

	//set in redis
	client := setting.RedisConn()
	count, _ := client.Get("healing2020:rankCount").Float64()
	keyName := "healing2020:Song." + strconv.FormatFloat(count/100+5.10, 'f', 2, 64)
	client.Set(keyName, jsonRank, 0)

	return nil
}

func GetSongRank(userid uint) ([]AllRank, string) {
	result := make([]AllRank, 10)
	client := setting.RedisConn()
	count, _ := client.Get("healing2020:rankCount").Float64()
	var i float64 = 0
	for j := 0; ; j++ {
		var rank []Rank
		if i*100 > count {
			break
		}
		var date float64 = 5.10 + i
		dateStr := strconv.FormatFloat(date, 'f', 2, 64)
		keyname := "healing2020:Song." + dateStr
		data, err := client.Get(keyname).Bytes()
		if err != nil {
			return nil, "Unexpected data"
		}
		json.Unmarshal(data, &rank)

		// ???????????????????????????
		for k := 0; k < len(rank); k++ {
			rank[k].IsPraise, _ = HasPraise(2, userid, rank[k].ID)
			rank[k].Praise = GetPraiseCount("song", rank[k].ID)
		}
		i = i + 0.01

		result[j].Data = rank
		result[j].Time = dateStr
	}

	//fmt.Println(result)
	return result, ""
}

func SendUserRank() error {
	//everyday's work

	//get from mysql
	db := setting.MysqlConn()

	var user []statements.User
	var allRank [][]Rank = make([][]Rank, 5)
	for i := 0; i < 5; i++ {
		var rank []Rank = make([]Rank, 10)
		pattern := []string{"", "??????", "??????", "????????????", "????????????"}
		var result *gorm.DB
		if i == 0 {
			result = db.Model(&statements.User{}).Order("money desc,created_at").Find(&user)
		} else {
			result = db.Model(&statements.User{}).Where("campus=?", pattern[i]).Order("money desc,created_at").Find(&user)
		}
		rows := result.RowsAffected
		if rows == 0 {
			allRank[i] = rank
			continue
		}
		if result.Error != nil {
			return result.Error
		}
		for i := 0; i < min(int(rows), 10); i++ {
			rank[i].ID = user[i].ID
			rank[i].UserId = user[i].ID
			rank[i].User = user[i].NickName
			rank[i].Avatar = user[i].Avatar

			if user[i].Setting1 == 0 {
				rank[i].Avatar = tools.GetAvatarUrl(user[i].Sex)
			}
		}

		allRank[i] = rank
	}
	jsonRank, _ := json.Marshal(allRank)

	//set in redis
	client := setting.RedisConn()
	keyName := "healing2020:User"
	client.Set(keyName, jsonRank, 0)

	return nil
}

type AllUserRank struct {
	Data [][]Rank `json:"data"`
}

func GetAllUserRank() (AllUserRank, string) {
	var result AllUserRank
	var rank [][]Rank
	client := setting.RedisConn()

	keyname := "healing2020:User"
	data, err := client.Get(keyname).Bytes()
	if err != nil {
		return AllUserRank{}, "Unexpected data"
	}
	json.Unmarshal(data, &rank)

	result.Data = rank
	return result, ""
}

type UserRank struct {
	Rank int `json:"rank"`
}

func GetUserRank(id string) (UserRank, error) {
	intId, _ := strconv.Atoi(id)
	userId := uint(intId)
	db := setting.MysqlConn()

	rows, err := db.Model(&statements.User{}).Order("money desc,created_at").Rows()
	rank := 0
	if err != nil {
		return UserRank{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var user statements.User
		db.ScanRows(rows, &user)
		if user.ID == userId {
			break
		}
		rank++
	}
	var userRank UserRank
	userRank.Rank = rank
	return userRank, err
}

func UpdateRankCount() {
	client := setting.RedisConn()

	count, _ := client.Get("healing2020:rankCount").Float64()
	count++
	client.Set("healing2020:rankCount", count, 0)
}

func min(a int, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

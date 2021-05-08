package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"encoding/json"
	"errors"
	"fmt"
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

	Time   string `json:"time"`
	Praise int    `json:"praise"`
	Name   string `json:"name"`
    UserId uint   `json:"userid"`
    IsPraise bool `json:"isPraise"`
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

	status := 0
	tx := db.Begin()

	var deliver []statements.Deliver
	year, mon, day := time.Now().Date()
	date := strconv.Itoa(year) + "_" + monthTransfer(mon.String()) + "_" + dayTransfer(strconv.Itoa(day)) + "%"
	//fmt.Println(date)
	result := tx.Model(&statements.Deliver{}).Where("created_at LIKE ?", date).Order("praise, created_at desc").Find(&deliver)
	err1 := errRollBack(tx, &status)
	if err1 != nil {
		return err1
	}
	rows := result.RowsAffected
	if rows == 0 {
		client := setting.RedisConn()

		count, _ := client.Get("healing2020:rankCount").Float64()
		keyName := "healing2020:Deliver." + strconv.FormatFloat(count/100+4.20, 'f', 2, 64)
		client.Set(keyName, "", 0)
		return errors.New("no data")
	}
	var rank []Rank = make([]Rank, 1)
	for i := 0; i < min(int(rows), 10); i++ {
		rank[i].ID = deliver[i].ID
		rank[i].Type = deliver[i].Type
		rank[i].Text = deliver[i].TextField
		rank[i].Photo = deliver[i].Photo
		rank[i].Source = deliver[i].Record
        rank[i].UserId = deliver[i].UserId

		userid := deliver[i].UserId
		var user statements.User
		status = 0
		tx2 := db.Begin()
		result2 := tx2.Model(&statements.User{}).Where("id=?", userid).First(&user)
		err2 := errRollBack(result2, &status)
		if err2 != nil {
			return err2
		}
		rank[i].User = user.NickName
		rank[i].Avatar = user.Avatar
	}
	jsonRank, _ := json.Marshal(rank)

	//set in redis
	client := setting.RedisConn()
	count, _ := client.Get("healing2020:rankCount").Float64()
	keyName := "healing2020:Deliver." + strconv.FormatFloat(count/100+4.20, 'f', 2, 64)
	client.Set(keyName, jsonRank, 0)

	return nil
}

func GetDeliverRank(userid uint) ([]AllRank, string) {
	result := make([]AllRank, 10)
	client := setting.RedisConn()
	count, _ := client.Get("healing2020:rankCount").Float64()
	var i float64 = 0
	for j := 0; ; j++ {
		rank := make([]Rank, 1)
		//fmt.Println(count)
		//fmt.Println(i)
		if i*100 > count {
			break
		}
		var date float64 = 4.20 + i
		dateStr := strconv.FormatFloat(date, 'f', 2, 64)
		keyname := "healing2020:Deliver." + dateStr
		//fmt.Println(dateStr)
		data, err := client.Get(keyname).Bytes()
		if err != nil {
			fmt.Println(err)
			return nil, "Unexpected data"
		}
		json.Unmarshal(data, &rank)
        // 把是否点赞的项拼上
        for k:=0;k<len(rank);k++ {
            rank[k].IsPraise,_ = HasPraise(1,userid,rank[k].ID)
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

	status := 0
	tx := db.Begin()

	var song []statements.Song
	year, mon, day := time.Now().Date()
	date := strconv.Itoa(year) + "_" + monthTransfer(mon.String()) + "_" + dayTransfer(strconv.Itoa(day))
	result := tx.Model(&statements.Song{}).Where("created_at LIKE ?", date+"%").Order("praise, created_at desc").Find(&song)
	err1 := errRollBack(tx, &status)
	if err1 != nil {
		return err1
	}
	rows := result.RowsAffected
	if rows == 0 {
		client := setting.RedisConn()

		count, _ := client.Get("healing2020:rankCount").Float64()
		keyName := "healing2020:Song." + strconv.FormatFloat(count/100+4.20, 'f', 2, 64)
		client.Set(keyName, "", 0)
		return errors.New("no data")
	}
	var rank []Rank = make([]Rank, 10)
	for i := 0; i < min(int(rows), 10); i++ {
		rank[i].ID = song[i].ID
		rank[i].Name = song[i].Name
		rank[i].Praise = song[i].Praise
		rank[i].Time = date
		rank[i].Source = song[i].Source

		userid := song[i].UserId
		var user statements.User
		status = 0
		tx2 := db.Begin()
		result2 := tx2.Model(&statements.User{}).Where("id=?", userid).First(&user)
		err2 := errRollBack(result2, &status)
		if err2 != nil {
			return err2
		}
		rank[i].User = user.NickName
		rank[i].Avatar = user.Avatar
	}
	jsonRank, _ := json.Marshal(rank)

	//set in redis
	client := setting.RedisConn()
	count, _ := client.Get("healing2020:rankCount").Float64()
	keyName := "healing2020:Song." + strconv.FormatFloat(count/100+4.20, 'f', 2, 64)
	client.Set(keyName, jsonRank, 0)

	return nil
}

func GetSongRank(userid uint) ([]AllRank, string) {
	result := make([]AllRank, 1)
	client := setting.RedisConn()
	count, _ := client.Get("healing2020:rankCount").Float64()
	var i float64 = 0
	for j := 0; ; j++ {
		var rank []Rank
		if i*100 > count {
			break
		}
		var date float64 = 4.20 + i
		dateStr := strconv.FormatFloat(date, 'f', 2, 64)
		keyname := "healing2020:Song." + dateStr
		data, err := client.Get(keyname).Bytes()
		if err != nil {
			return nil, "Unexpected data"
		}
		json.Unmarshal(data, &rank)

        // 把是否点赞的项拼上
        for k:=0;k<len(rank);k++ {
            rank[k].IsPraise,_ = HasPraise(1,userid,rank[k].ID)
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
	var rank []Rank = make([]Rank, 1)
	var allRank [][]Rank = make([][]Rank, 3)
	for i := 0; i < 3; i++ {
		pattern := []string{"", "中大", "华工"}
		var result *gorm.DB
		if i == 0 {
			result = db.Model(&statements.User{}).Order("money, created_at desc").Find(&user)
		} else {
			result = db.Model(&statements.User{}).Where("campus=?", pattern[i]).Order("money, created_at desc").Find(&user)
		}
		rows := result.RowsAffected
		if rows == 0 {
			return nil
		}
		if result.Error != nil {
			return result.Error
		}
		for i := 0; i < min(int(rows), 10); i++ {
			rank[i].ID = user[i].ID
			rank[i].User = user[i].NickName
			rank[i].Avatar = user[i].Avatar
		}

		allRank[i] = rank
	}
	jsonRank, _ := json.Marshal(allRank)

	//set in redis
	client := setting.RedisConn()
	count, _ := client.Get("healing2020:rankCount").Float64()
	keyName := "healing2020:User." + strconv.FormatFloat(count/100+4.20, 'f', 2, 64)
	client.Set(keyName, jsonRank, 0)

	return nil
}

type AllUserRank struct {
	Time string   `json:"time"`
	Data [][]Rank `json:"data"`
}

func GetAllUserRank() ([]AllUserRank, string) {
	result := make([]AllUserRank, 10)
	client := setting.RedisConn()
	count, _ := client.Get("healing2020:rankCount").Float64()
	var i float64 = 0
	for j := 0; ; j++ {
		var rank [][]Rank
		if i*100 > count {
			fmt.Println(j)
			break
		}
		var date float64 = 4.20 + i
		dateStr := strconv.FormatFloat(date, 'f', 2, 64)
		keyname := "healing2020:User." + dateStr
		data, err := client.Get(keyname).Bytes()
		if err != nil {
			return nil, "Unexpected data"
		}
		json.Unmarshal(data, &rank)
		i = i + 0.01

		result[j].Data = rank
		result[j].Time = dateStr
	}

	//fmt.Println(result)
	return result, ""
}

type UserRank struct {
	Rank int `json:"rank"`
}

func GetUserRank(id string) (UserRank, error) {
	intId, _ := strconv.Atoi(id)
	userId := uint(intId)
	db := setting.MysqlConn()

	rows, err := db.Model(&statements.User{}).Order("Money desc").Rows()
	rank := 0
	if err != nil {
		return UserRank{}, err
	}
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

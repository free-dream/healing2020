package models

import (
	"database/sql"
	"encoding/json"
	"strconv"

	//"fmt"
	"errors"
	"strings"
	"time"

	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"healing2020/pkg/tools"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type MainMsg struct {
	Sing   []SongMsg `json:"sing"`
	Listen []SongMsg `json:"listen"`
}

type SongMsg struct {
	User   string    `json:"user"`
	Name   string    `json:"name"`
	Avatar string    `json:"avatar"`
	Time   time.Time `json:"time"`
	Sex    int       `json:"sex"`
	More   string    `json:"more"`

	Id       uint   `json:"id"`
	SongId   uint   `json:"songId"`
	Like     int    `json:"like"`
	Style    string `json:"style"`
    Language string `json:"language"`
	Source   string `json:"source"`
	Singer   string `json:"singer"`
	UserId   uint   `json:"userid"`
    Phone    string `json:"phone"`
	IsPraise bool   `json:"isPraise"`
}

func SendMainMsg() {
	client := setting.RedisConn()

	var sortArr = []string{"0", "1"}
	var keyArr = []string{"", "rap", "韩语","ACG", "流行", "古风", "民谣", "摇滚", "抖音热歌", "其他", "国语", "英语", "日语", "粤语"}
	//var sortArr = []string{"0"}
	//var keyArr = []string{""}
	for _, sort := range sortArr {
		for _, key := range keyArr {
			listenRaw := LoadSongMsg(sort, key, "")
			singRaw := LoadVodMsg(sort, key, "")
			listen, _ := json.Marshal(listenRaw)
			sing, _ := json.Marshal(singRaw)

			keyname := "healing2020:Main:" + key + "ListenMsg" + sort
			client.Set(keyname, listen, 0)
			keyname = "healing2020:Main:" + key + "SingMsg" + sort
			client.Set(keyname, sing, 0)
			//fmt.Println(result.Err())
		}
	}
}

func max(a int, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func LoadSongMsg(sort string, key string, userTags string) []SongMsg {
	db := setting.MysqlConn()
	count := 0
	db.Raw("select count(*) from song").Row().Scan(&count)
	var songList []SongMsg = make([]SongMsg, max(8, count))
	i := 0

	var rows *sql.Rows
	var result *gorm.DB
	if key == "" || key == "推荐" {
		if sort == "0" {
			result = db.Raw("select id,user_id,vod_id,name,praise,source,style,language,created_at from song where is_hide = 0 order by rand()")
			rows, _ = result.Rows()
		} else {
			result = db.Raw("select id,user_id,vod_id,name,praise,source,style,language,created_at from song where is_hide = 0 order by created_at desc")
			rows, _ = result.Rows()
		}
	} else {
		if sort == "0" {
			result = db.Raw("select id,user_id,vod_id,name,praise,source,style,language,created_at from song where is_hide = 0 and style=? or is_hide = 0 and language=? order by rand()", key, key)
			rows, _ = result.Rows()
		} else {
			result = db.Raw("select id,user_id,vod_id,name,praise,source,style,language,created_at from song where is_hide = 0 and style=? or is_hide = 0 and language=? order by created_at desc", key, key)
			rows, _ = result.Rows()
		}
	}

	if rows == nil {
		return songList
	}

	defer rows.Close()

	for rows.Next() {
		var song statements.Song
		db.ScanRows(rows, &song)
		if key == "推荐" && !recommendFilter(song.Style, song.Language, userTags) {
			continue
		}
		songList[i].SongId = song.ID
		songList[i].Source = song.Source
		songList[i].Name = song.Name
		songList[i].Style = song.Style
        songList[i].Language = song.Language
		songList[i].Time = song.CreatedAt
		userId := song.UserId
		vodId := song.VodId
		songList[i].Id = vodId

		var user statements.User
		db.Model(&statements.User{}).Select("setting1,nick_name,sex,avatar").Where("id=?", userId).Find(&user)
		songList[i].User = user.NickName
		songList[i].Sex = user.Sex
		songList[i].Avatar = user.Avatar
		songList[i].UserId = userId

        if user.Setting1 == 0 {
            songList[i].Avatar = tools.GetAvatarUrl(user.Sex)
        }

		var vod statements.Vod
		db.Model(&statements.Vod{}).Select("more").Where("id=?", vodId).Find(&vod)
		songList[i].More = vod.More

		i++
	}

	return songList
}

func LoadVodMsg(sort string, key string, userTags string) []SongMsg {
	db := setting.MysqlConn()
	count := 0
	db.Raw("select count(*) from vod").Row().Scan(&count)
	var vodList []SongMsg = make([]SongMsg, max(8, count))

	var rows *sql.Rows
	var result *gorm.DB
	if key == "" || key == "推荐" {
		if sort == "0" {
			result = db.Raw("select id,user_id,name,singer,more,style,language,created_at,hide_name from vod order by rand()")
			rows, _ = result.Rows()
		} else {
			result = db.Raw("select id,user_id,name,singer,more,style,language,created_at,hide_name from vod order by created_at desc")
			rows, _ = result.Rows()
		}
	} else {
		if sort == "0" {
			result = db.Raw("select id,user_id,name,singer,more,created_at,style,language,hide_name from vod where style=? or language=? order by rand()", key, key)
			rows, _ = result.Rows()
		} else {
			result = db.Raw("select id,user_id,name,singer,more,created_at,style,language,hide_name from vod where style=? or language=? order by created_at desc", key, key)
			rows, _ = result.Rows()
		}
	}

	if rows == nil {
		return vodList
	}

	defer rows.Close()

	i := 0
	for rows.Next() {
		var vod statements.Vod
		db.ScanRows(rows, &vod)
		if key == "推荐" && !recommendFilter(vod.Style, vod.Language, userTags) {
			continue
		}
		vodList[i].Id = vod.ID
		vodList[i].Name = vod.Name
		vodList[i].More = vod.More
		vodList[i].Time = vod.CreatedAt
		vodList[i].Singer = vod.Singer
		vodList[i].UserId = vod.UserId
        vodList[i].Style = vod.Style
        vodList[i].Language = vod.Language

		userid := vod.UserId

		var user statements.User
		db.Model(&statements.User{}).Select("phone, nick_name,sex,avatar,setting1").Where("id=?", userid).Find(&user)
		vodList[i].User = user.NickName
		vodList[i].Sex = user.Sex
		vodList[i].Avatar = user.Avatar
        vodList[i].Phone = user.Phone

        if user.Setting1 == 0 {
            vodList[i].Avatar = tools.GetAvatarUrl(user.Sex)
        }

		if vod.HideName == 1 {
			vodList[i].UserId = 0
			vodList[i].User = "匿名用户"
			vodList[i].Avatar = tools.GetAvatarUrl(user.Sex)
		}

		i++
	}

	return vodList
}

func GetMainMsg(pageStr string, sort string, key string, tags string, userid uint) (MainMsg, error) {
	page, _ := strconv.Atoi(pageStr)
	var result MainMsg
	client := setting.RedisConn()

    // 推荐部分先处理
	if tags != "" {
		// listen := LoadSongMsg(sort, "推荐", tags)
		// sing := LoadVodMsg(sort, "推荐", tags)
		// resultSing, resultListen, err := Paging(page, sing, listen)

		// if err != nil {
		// 	return result, errors.New("page out of range")
		// }

        // 从redis拿'全部'的数据
        data1, err1 := client.Get("healing2020:Main:" + "" + "SingMsg" + sort).Bytes()
        if data1 == nil {
            return MainMsg{}, nil
        }
        if err1 != nil {
            return MainMsg{}, err1
        }
        data2, err2 := client.Get("healing2020:Main:" + "" + "ListenMsg" + sort).Bytes()
        if data2 == nil {
            return MainMsg{}, nil
        }
        if err2 != nil {
            return MainMsg{}, err2
        }

        var sing []SongMsg
        var listen []SongMsg
        json.Unmarshal(data1, &sing)
        json.Unmarshal(data2, &listen)
        resultSing, resultListen, err := Paging(page, sing, listen)

        if err != nil {
            return result, errors.New("page out of range")
        }

		//塞进是否点赞,排除用户不喜欢的项
		for i := 0; i < len(resultListen); i++ {
            if !recommendFilter(resultListen[i].Style,resultListen[i].Language,tags) {
                resultListen[i].Id = 0
                continue
            }
			resultListen[i].IsPraise, _ = HasPraise(2, userid, resultListen[i].SongId)
			resultListen[i].Like = GetPraiseCount("song", resultListen[i].SongId)
		}

		for i := 0; i < len(resultSing); i++ {
            if !recommendFilter(resultSing[i].Style,resultSing[i].Language,tags) {
                resultSing[i].Id = 0
                continue
            }
		}

		result.Sing = resultSing
		result.Listen = resultListen

		return result, nil
	}

	data1, err1 := client.Get("healing2020:Main:" + key + "SingMsg" + sort).Bytes()
	if data1 == nil {
		return MainMsg{}, nil
	}
	if err1 != nil {
		return MainMsg{}, err1
	}
	data2, err2 := client.Get("healing2020:Main:" + key + "ListenMsg" + sort).Bytes()
	if data2 == nil {
		return MainMsg{}, nil
	}
	if err2 != nil {
		return MainMsg{}, err2
	}

	var sing []SongMsg
	var listen []SongMsg
	json.Unmarshal(data1, &sing)
	json.Unmarshal(data2, &listen)
	resultSing, resultListen, err := Paging(page, sing, listen)

	if err != nil {
		return result, errors.New("page out of range")
	}

	//塞进是否点赞
	for i := 0; i < len(resultListen); i++ {
		resultListen[i].IsPraise, _ = HasPraise(2, userid, resultListen[i].SongId)
		resultListen[i].Like = GetPraiseCount("song", resultListen[i].SongId)
	}

	result.Sing = resultSing
	result.Listen = resultListen

	return result, nil
}

func Paging(page int, data1 []SongMsg, data2 []SongMsg) ([]SongMsg, []SongMsg, error) {
	var result1 []SongMsg
	var result2 []SongMsg
	if (page-1)*20 > len(data1) {
		result1 = make([]SongMsg, 1)
	} else {
		result1 = make([]SongMsg, 20)
	}
	if (page-1)*20 > len(data2) {
		result2 = make([]SongMsg, 1)
	} else {
		result2 = make([]SongMsg, 20)
	}
	if len(result1) == 1 && len(result2) == 1 {
		return result1, result2, errors.New("page out of page")
	}
	for i := 0; i < 20; i++ {
		if (page-1)*20+i >= len(data1) {
			break
		}
		result1[i] = data1[(page-1)*20+i]
	}
	for i := 0; i < 20; i++ {
		if (page-1)*20+i >= len(data2) {
			break
		}
		result2[i] = data2[(page-1)*20+i]
	}

	return result1, result2, nil
}

func isListNil(result *gorm.DB) bool {
	count := result.RowsAffected
	if count == 0 {
		return true
	}
	return false
}

func tagsSplit(tags string) []string {
	// tags样式 "流行，国语，古风，..."
	return strings.Split(tags, ",")
}

func recommendFilter(style string, language string, userTags string) bool {
	// 把不是用户爱好的过滤
	tags := tagsSplit(userTags)
	for _, tag := range tags {
		if tag == style || tag == language {
			return true
		}
	}

	return false
}

type SearchResp struct {
	User []UserResp
	Song []SongResp
	Vod  []VodResp
	Err  string
}

type UserResp struct {
	UserId   uint   `json:"userid"`
	UserName string `json:"userName"`
	Avatar   string `json:"avatar"`
	More     string `json:"more"`
	Bg       int    `json:"background"`
}

type SongResp struct {
	SongId   uint      `json:"songid"`
	VodId    uint      `json:"vodid"`
	SongName string    `json:"name"`
	Avatar   string    `json:"avatar"`
	Praise   int       `json:"like"`
	Source   string    `json:"source"`
	Singer   string    `json:"user"`
	Time     time.Time `json:"time"`
}

type VodResp struct {
	VodId   uint      `json:"vodid"`
	VodName string    `json:"name"`
	VodUser string    `json:"user"`
    VodUserId uint    `json:"userid"`
    Phone   string    `json:"phone"`
	Avatar  string    `json:"avatar"`
	More    string    `json:"more"`
	Sex     int       `json:"sex"`
	Time    time.Time `json:"time"`
}

func GetSearchResult(search string) SearchResp {
	db := setting.MysqlConn()

	var searchResp SearchResp
	var songResp []SongResp
	var userResp []UserResp
	var vodResp []VodResp

	var songCount int = 0
	result := db.Model(&statements.Song{}).Where("name = ? and is_hide = 0", search).Select("id,vod_id,source,created_at,user_id").Count(&songCount)
	if songCount != 0 && result.Error == nil {
		rows, _ := result.Rows()
		defer rows.Close()

		songResp = make([]SongResp, songCount)

		i := 0
		for rows.Next() {
			var song statements.Song
			db.ScanRows(rows, &song)
			songResp[i].SongId = song.ID
			songResp[i].VodId = song.VodId
			songResp[i].Praise = GetPraiseCount("song", song.ID)
			songResp[i].Source = song.Source
			songResp[i].SongName = search
			songResp[i].Time = song.CreatedAt

			var user statements.User
			db.Model(&statements.User{}).Select("nick_name,avatar").Where("id = ?", song.UserId).First(&user)
			songResp[i].Singer = user.NickName
			songResp[i].Avatar = user.Avatar

			i++
		}
	} else {
		if result.Error != nil {
			searchResp.Err = result.Error.Error()
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			searchResp.Err = ""
		}
	}

	var vodCount int = 0
	result = db.Model(&statements.Vod{}).Where("name = ?", search).Select("id, more, created_at, user_id, hide_name").Count(&vodCount)

	if vodCount != 0 && result.Error == nil {
		rows, _ := result.Rows()
		defer rows.Close()

		vodResp = make([]VodResp, vodCount)

		i := 0
		for rows.Next() {
			var vod statements.Vod
			db.ScanRows(rows, &vod)
			vodResp[i].VodId = vod.ID
			vodResp[i].VodName = search
			vodResp[i].Time = vod.CreatedAt
			vodResp[i].More = vod.More
            vodResp[i].VodUserId = vod.UserId

			var user statements.User
			db.Model(&statements.User{}).Select("phone, setting1, sex, more, avatar, nick_name").Where("id = ?", vod.UserId).First(&user)
			vodResp[i].VodUser = user.NickName
			vodResp[i].Avatar = user.Avatar
			vodResp[i].Sex = user.Sex
            vodResp[i].Phone = user.Phone

            if user.Setting1 == 0 {
                vodResp[i].Avatar = tools.GetAvatarUrl(user.Sex)
            }

			if vod.HideName == 1 {
				vodResp[i].VodUser = "匿名用户"
				vodResp[i].Avatar = tools.GetAvatarUrl(user.Sex)
                vodResp[i].VodUserId = 0
			}

			i++
		}
	} else {
		if result.Error != nil {
			searchResp.Err = result.Error.Error()
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			searchResp.Err = ""
		}
	}

	userCount := 0
	result = db.Model(&statements.User{}).Select("id,nick_name,more,avatar").Where("nick_name=? or true_name = ? or phone =?", search, search, search).Count(&userCount)
	if userCount != 0 && result.Error == nil {
		rows, _ := result.Rows()
		defer rows.Close()

		userResp = make([]UserResp, userCount)

		i := 0
		for rows.Next() {
			var user statements.User
			db.ScanRows(rows, &user)
			userResp[i].UserId = user.ID
			userResp[i].More = user.More
			userResp[i].Avatar = user.Avatar
			userResp[i].UserName = user.NickName

            if user.Setting1 == 0 {
                userResp[i].Avatar = tools.GetAvatarUrl(user.Sex)
            }

			var userOther statements.UserOther
			db.Model(&statements.UserOther{}).Select("now").Where("user_id = ?", user.ID).First(&userOther)
			userResp[i].Bg = userOther.Now
			i++
		}
	} else {
		if result.Error != nil {
			searchResp.Err = result.Error.Error()
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			searchResp.Err = ""
		}
	}

	searchResp.User = userResp
	searchResp.Song = songResp
	searchResp.Vod = vodResp

	return searchResp
}

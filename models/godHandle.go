package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"

	//"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"strconv"
)

func CreateDelivers(id string, types int, textfield string, photo string, record string, praise int) error {
	intId, _ := strconv.Atoi(id)
	userid := uint(intId)
	var deliver statements.Deliver
	deliver.UserId = userid
	deliver.Type = types
	deliver.TextField = textfield
	deliver.Photo = photo
	deliver.Record = record
	deliver.Praise = praise
	db := setting.MysqlConn()
	defer db.Close()
	result := db.Model(&statements.Deliver{}).Create(&deliver)
	return result.Error
}

//func CreateBgs(id string, now int, b1 int, b2 int, b3 int, b4 int, b5 int) error {
//	intId, _ := strconv.Atoi(id)
//	userid := uint(intId)
//	var bg statements.Background
//	bg.UserId = userid
//	bg.Now = now
//	bg.B1 = b1
//	bg.B2 = b2
//	bg.B3 = b3
//	bg.B4 = b4
//	bg.B5 = b5
//	db := setting.MysqlConn()
//	defer db.Close()
//	result := db.Model(&statements.Background{}).Create(&bg)
//	return result.Error
//}

func CreateComments(id string, types int, id1 string, id2 string, content string) error {
	intId1, _ := strconv.Atoi(id1)
	intId2, _ := strconv.Atoi(id2)
	intId, _ := strconv.Atoi(id)
	userid := uint(intId)
	songid := uint(intId1)
	deliverid := uint(intId2)
	var comment statements.Comment
	comment.SongId = songid
	comment.UserId = userid
	comment.DeliverId = deliverid
	comment.Type = types
	comment.Content = content
	db := setting.MysqlConn()
	defer db.Close()
	result := db.Model(&statements.Comment{}).Create(&comment)
	return result.Error
}

func CreateLotterys(pid int, uid int, weight int) error {
	prizeId := uint(pid)
	userId := uint(uid)
	var lottery statements.Lottery
	lottery.PrizeId = prizeId
	lottery.UserId = userId
	lottery.Weight = weight
	db := setting.MysqlConn()
	defer db.Close()
	result := db.Model(&statements.Lottery{}).Create(&lottery)
	return result.Error
}

func CreateMailboxs(message string) error {
	var mailbox statements.Mailbox
	mailbox.Message = message
	db := setting.MysqlConn()
	defer db.Close()
	result := db.Model(&statements.Mailbox{}).Create(&mailbox)
	return result.Error
}

func CreateMessages(id1 string, id2 string, types int, content string, url string) error {
	var message statements.Message
	intId1, _ := strconv.Atoi(id1)
	intId2, _ := strconv.Atoi(id2)
	send := uint(intId1)
	receive := uint(intId2)
	message.Send = send
	message.Receive = receive
	message.Type = types
	message.Content = content
	message.Url = url
	db := setting.MysqlConn()
	defer db.Close()
	result := db.Model(&statements.Message{}).Create(&message)
	return result.Error
}

func CreatePraises(id1 string, types int, id2 string) error {
	var praise statements.Praise
	intId1, _ := strconv.Atoi(id1)
	intId2, _ := strconv.Atoi(id2)
	userid := uint(intId1)
	praiseid := uint(intId2)
	praise.UserId = userid
	praise.Type = types
	praise.PraiseId = praiseid
	db := setting.MysqlConn()
	defer db.Close()
	result := db.Model(&statements.Praise{}).Create(&praise)
	return result.Error
}

func CreatePrizes(name string, intro string, photo string, weight int) error {
	var prize statements.Prize
	prize.Name = name
	prize.Intro = intro
	prize.Photo = photo
	prize.Weight = weight
	db := setting.MysqlConn()
	defer db.Close()
	result := db.Model(&statements.Prize{}).Create(&prize)
	return result.Error
}

func CreateRanks(campus string, allrank string, partrank string) error {
	var rank statements.Rank
	rank.Campus = campus
	rank.AllRank = allrank
	rank.PartRank = partrank
	db := setting.MysqlConn()
	defer db.Close()
	result := db.Model(&statements.Rank{}).Create(&rank)
	return result.Error
}

func CreateSongs(id1 string, id2 string, id3 string, name string, praise int, source string, style string, language string) error {
	intId1, _ := strconv.Atoi(id1)
	intId2, _ := strconv.Atoi(id2)
	intId3, _ := strconv.Atoi(id3)
	userid := uint(intId1)
	vodid := uint(intId2)
	vodsend := uint(intId3)
	var song statements.Song
	song.UserId = userid
	song.VodId = vodid
	song.VodSend = vodsend
	song.Name = name
	song.Praise = praise
	song.Source = source
	song.Style = style
	song.Language = language
	db := setting.MysqlConn()
	defer db.Close()
	result := db.Model(&statements.Song{}).Create(&song)
	return result.Error
}

func CreateSpecials(id1 string, id2 string, name string, praise int, song string) error {
	intId1, _ := strconv.Atoi(id1)
	intId2, _ := strconv.Atoi(id2)
	subjectid := uint(intId1)
	userid := uint(intId2)
	var special statements.Special
	special.SubjectId = subjectid
	special.UserId = userid
	special.Name = name
	special.Praise = praise
	special.Song = song
	db := setting.MysqlConn()
	defer db.Close()
	result := db.Model(&statements.Special{}).Create(&special)
	return result.Error
}

func CreateSubjects(name string, intro string) error {
	var subject statements.Subject
	subject.Name = name
	subject.Intro = intro
	db := setting.MysqlConn()
	defer db.Close()
	result := db.Model(&statements.Subject{}).Create(&subject)
	return result.Error
}

func CreateUsers(openid string, nick string, name string, more string, avatar string, campus string, phone string, sex int, hobby string, money int, setting1 int, setting2 int, setting3 int) error {
	var user statements.User
	user.OpenId = openid
	user.NickName = nick
	user.TrueName = name
	user.More = more
	user.Avatar = avatar
	user.Campus = campus
	user.Phone = phone
	user.Sex = sex
	user.Hobby = hobby
	user.Money = money
	user.Setting1 = setting1
	user.Setting2 = setting2
	user.Setting3 = setting3
	db := setting.MysqlConn()
	defer db.Close()
	result := db.Model(&statements.User{}).Create(&user)
	return result.Error
}

func CreateVods(id string, more string, name string, singer string, style string, language string) error {
	intId, _ := strconv.Atoi(id)
	userid := uint(intId)
	var vod statements.Vod
	vod.UserId = userid
	vod.Name = name
	vod.More = more
	vod.Singer = singer
	vod.Style = style
	vod.Language = language
	db := setting.MysqlConn()
	defer db.Close()
	result := db.Model(&statements.Vod{}).Create(&vod)
	return result.Error
}

func CreateFakeUserOther(id string, user_id uint, remainHideName int, remainVod int) error {
	intId, _ := strconv.Atoi(id)
	ID := uint(intId)
	var user_other statements.UserOther
	user_other.ID = ID
	user_other.UserId = user_id
	user_other.RemainHideName = remainHideName
	user_other.RemainSing = remainVod
	db := setting.MysqlConn()
	defer db.Close()
	result := db.Model(&statements.Vod{}).Create(&user_other)
	return result.Error
}

func TableCleanUp() {
	db := setting.MysqlConn()
	defer db.Close()

	//db.Exec("TRUNCATE TABLE background")
	db.Exec("TRUNCATE TABLE comment")
	db.Exec("TRUNCATE TABLE deliver")
	db.Exec("TRUNCATE TABLE lottery")
	db.Exec("TRUNCATE TABLE mailbox")
	//db.Exec("TRUNCATE TABLE message")
	db.Exec("TRUNCATE TABLE praise")
	db.Exec("TRUNCATE TABLE prize")
	//db.Exec("TRUNCATE TABLE rank")
	db.Exec("TRUNCATE TABLE song")
	db.Exec("TRUNCATE TABLE special")
	db.Exec("TRUNCATE TABLE subject")
	db.Exec("TRUNCATE TABLE user")
	db.Exec("TRUNCATE TABLE vod")
}

package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type AllComment struct {
	CommentResponse Comment
	NickName        string
}

type Comment struct {
	UserID    uint   `json:"UserID"`
	Type      int    `json:"Type"`
	SongID    uint   `json:"song_id"`
	DeliverID uint   `json:"deliver_id"`
	Content   string `json:"content"`
}

func GetComment(strID string, Type string) ([]AllComment, error) {
	var err error
	intID, _ := strconv.Atoi(strID)
	id := uint(intID)
	//连接mysql
	db := setting.MysqlConn()
	defer db.Close()

	//获取评论其他信息
	var commentElse []Comment
	if Type == "2" { //投递评论
		err = db.Table("comment").Select("user_id, type, deliver_id, content").Where("type = 2 AND deliver_id = ?", id).Scan(&commentElse).Error
		if err != nil {
			return nil, err
		}
	}
	if Type == "1" { //歌房评论
		err = db.Table("comment").Select("user_id, type, song_id, content").Where("type = 1 AND song_id = ?", id).Scan(&commentElse).Error
		if err != nil {
			return nil, err
		}
	}
	
	//获取评论人昵称信息
	commentName := make([]statements.User, len(commentElse))
	for i := 0; i < len(commentElse); i++ {
		err = db.Table("user").Select("nick_name").Where("id = ?", commentElse[i].UserID).First(&commentName[i]).Error
		if err != nil {
			return nil, err
		}
	}

	responseComment := make([]AllComment, len(commentElse))
	for i := 0; i < len(commentElse); i++ {
		responseComment[i] = AllComment{
			CommentResponse: commentElse[i],
			NickName:        commentName[i].NickName,
		}
	}
	return responseComment, err
}

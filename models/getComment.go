package models

import (
	"healing2020/models/statements"
	"healing2020/pkg/setting"
	"strconv"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type AllComment struct {
	CommentResponse Comment
	NickName        string `json:"nickname"`
	Avatar          string `json:"avatar"`
	IsPraise  bool      `json:"isPraise"`
}

type Comment struct {
	Id 		  int 		`json:"comment_id"`
	UserID    uint      `json:"UserID"`
	Type      int       `json:"Type"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
	Praise    int    	`json:"praise"`
}

func GetComment(strID string, Type string) ([]AllComment, error) {
	var err error
	intID, _ := strconv.Atoi(strID)
	id := uint(intID)
	//连接mysql
	db := setting.MysqlConn()

	//获取评论其他信息
	var commentElse []Comment
	if Type == "2" { //投递评论
		err = db.Table("comment").Select("id, user_id, type, created_at, content, praise").Where("type = 2 AND deliver_id = ?", id).Scan(&commentElse).Error
		if err != nil {
			return nil, err
		}
	}
	if Type == "1" { //歌房评论
		err = db.Table("comment").Select("id, user_id, type, created_at, content, praise").Where("type = 1 AND song_id = ?", id).Scan(&commentElse).Error
		if err != nil {
			return nil, err
		}
	}

	//获取评论人昵称信息
	commentName := make([]statements.User, len(commentElse))
	for i := 0; i < len(commentElse); i++ {
		err = db.Table("user").Select("nick_name, avatar").Where("id = ?", commentElse[i].UserID).First(&commentName[i]).Error
		if err != nil {
			return nil, err
		}
	}

	responseComment := make([]AllComment, len(commentElse))
	for i := 0; i < len(commentElse); i++ {
		responseComment[i] = AllComment{
			CommentResponse: commentElse[i],
			NickName:        commentName[i].NickName,
			Avatar:          commentName[i].Avatar,
		}
		responseComment[i].IsPraise, _ = HasPraise(4, commentElse[i].UserID, uint(commentElse[i].Id))
	}
	return responseComment, err
}

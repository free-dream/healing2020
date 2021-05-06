package controller

import (
	"healing2020/models"
	"healing2020/pkg/e"

	"github.com/gin-gonic/gin"
)

func autoCreate() {
	models.CreateUsers("", "robot", "truename", "more", "avatar", "华工", "12312312311", 1, "摇滚", 100, 0, 0, 0)
	models.CreateUsers("", "robot2", "truename2", "more", "avatar2", "华工", "12312312312", 1, "古风", 30, 1, 0, 0)
	models.CreateUsers("", "robot3", "truename3", "more", "avatar3", "中大", "12182312312", 0, "古风", 130, 1, 1, 0)
	models.CreateUsers("", "robot4", "truename4", "more", "avatar4", "中大", "18182311310", 1, "日语", 120, 0, 1, 0)
	models.CreateUsers("", "robot5", "truename5", "more", "avatar5", "华工", "12312123412", 0, "爵士", 300, 1, 1, 1)

	//    models.CreateBgs("1",1,1,1,1,1,1)
	//    models.CreateBgs("2",2,1,1,1,1,1)
	//    models.CreateBgs("3",3,1,1,1,1,1)
	//    models.CreateBgs("4",4,1,1,1,1,1)
	//    models.CreateBgs("5",5,1,1,1,1,1)

	models.CreateVods("1", "no more", "abc", "tobor", "", "")
	models.CreateVods("2", "no more", "abcd", "2tobor", "", "")
	models.CreateVods("3", "more", "abde", "3tobor", "", "")
	models.CreateVods("4", "more", "bde", "4tobor", "", "")
	models.CreateVods("5", "more", "bdecc", "3tobor", "", "")

	models.CreateSongs("1", "2", "2", "abcd", 0, "source1", "", "")
	models.CreateSongs("3", "2", "2", "abcd", 0, "source2", "", "")
	models.CreateSongs("3", "1", "1", "abc", 0, "source3", "", "")
	models.CreateSongs("4", "1", "1", "abc", 0, "source4", "", "")
	models.CreateSongs("5", "3", "3", "abde", 4, "source5", "", "")
	models.CreateSongs("2", "5", "5", "bdecc", 3, "source6", "", "")

	models.CreatePraises("1", 1, "5")
	models.CreatePraises("2", 1, "5")
	models.CreatePraises("4", 1, "5")
	models.CreatePraises("3", 1, "5")
	models.CreatePraises("3", 1, "6")
	models.CreatePraises("1", 1, "6")
	models.CreatePraises("2", 1, "6")
	models.CreatePraises("4", 2, "1")

	models.CreateDelivers("2", 1, "i am robot", "", "", 1)
	models.CreateDelivers("5", 1, "i am root", "", "", 1)
	models.CreateDelivers("3", 1, "i am robt", "", "", 1)
	models.CreateDelivers("3", 2, "i am robot2", "photo1", "", 0)
	models.CreateDelivers("1", 2, "am robo2", "photo2", "", 0)
	models.CreateDelivers("1", 2, "i am robo2", "photo3", "", 0)
	models.CreateDelivers("4", 3, "i m robot5", "", "source7", 0)
	models.CreateDelivers("2", 3, "i m rot2", "", "source8", 0)
	models.CreateDelivers("3", 3, "i m rbot2", "", "source9", 0)
	models.CreateDelivers("1", 3, "i m obot2", "", "source0", 0)

	models.CreateComments("5", 1, "3", "", "wonderful")
	models.CreateComments("5", 2, "", "1", "me too")

	models.CreateFakeUserOther("1", 1, 0, 3)
	models.CreateFakeUserOther("2", 2, 1, 2)
	models.CreateFakeUserOther("3", 3, 2, 3)
	models.CreateFakeUserOther("4", 4, 2, 3)
}

func LoadTestData() {
	models.TableCleanUp()
	autoCreate()
}

func PostSubject(c *gin.Context) {
	ID := c.Query("subject_id")
	Name := c.Query("name")
	Photo := c.Query("photo")
	Intro := c.Query("intro")
	err := models.PostSubject(ID, Name, Photo, Intro)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to add subject"})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: "发送歌房成功！"})
}

func PostSpecial(c *gin.Context) {
	Subject_id := c.Query("subject_id")
	Song := c.Query("song")
	User_id := c.Query("user_id")
	err := models.PostSpecial(Subject_id, Song, User_id)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "Fail to add special"})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: "发送歌房成功！"})
}

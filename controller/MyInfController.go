package controller

import (
	"time"
)

type RequestSongs struct {	//点歌
	Song string `json:"song"`
	Time time.Time `json:"time"`
}
type SingSongs struct {		//唱歌
	Song string `json:"song"`
	Time time.Time `json:"time"`
	From string `json:"from"`
}
type Admire struct {		//点赞
	Song string `json:"song"`
	Time time.Time `json:"time"`
	From string `json:"from"`
	Number int `json:"number"`
}

type PersonalPage struct {
	NickName string `json:"name"`
	Campus string `json:"school"`
	More string  `json:"more"`
	Setting1 int `json:"setting1"`
	Setting2 int `json:"setting2"`
	Setting3 int `json:"setting3"`
	Avatar string `json:"avatar"`
	Vod RequestSongs `json:"requestSongs"`
	Songs SingSongs `json:"Songs"`
	Praise Admire `json:"admire"` 
}
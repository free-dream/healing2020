package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"healing2020/models"
	"healing2020/models/statements"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"
)

type BroadcastContent struct {
	Text string `json:"message"`
}

type Message struct {
	Type       int    `josn:"type"`
	FromUserID uint   `json:"fromUserID"`
	ToUserID   uint   `json:"toUserID" validate:"required"`
	Content    string `json:"content" validate:"required"`
}

type WsConnection struct {
	ws        *websocket.Conn
	userID    uint
	closeChan chan byte
	mutex     sync.Mutex
	isClosed  bool
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var MessageQueue = make(map[int](chan *Message))
var broadcastChan = make(chan string)

//@Title Broadcast
//@Description 广播
//@Tags message
//@Produce json
//@Param json body BroadcastContent true "广播信息"
//@Router /broadcast [post]
//@Success 200 {object} e.ErrMsgResponse
//@Failure 403 {object} e.ErrMsgResponse
func Broadcast(c *gin.Context) {
	json := BroadcastContent{}
	c.BindJSON(&json)

	broadcastChan <- json.Text

	err := models.CreateMailBox(json.Text)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "广播失败!"})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: e.GetMsg(e.SUCCESS)})
}

//@Title SendMessage
//@Description 发送消息并保存于数据库
//@Tags message
//@Produce json
//@Router /message [post]
//@Success 200 {object} e.ErrMsgResponse
//@Failure 403 {object} e.ErrMsgResponse
func SendMessage(c *gin.Context) {
	var msg Message
	c.BindJSON(&msg)
	msg.Type = 2 //此接口只用于发送文本信息
	msgDB := statements.Message{
		Send:    msg.FromUserID,
		Receive: msg.FromUserID,
		Content: msg.Content,
		Type:    msg.Type,
	}

	MessageQueue[int(msg.ToUserID)] <- &msg

	err := models.CreateMessage(msgDB)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "保存信息失败！"})
		return
	}

	c.JSON(200, e.ErrMsgResponse{Message: e.GetMsg(e.SUCCESS)})
}

//处理ws连接
func WsHandle(c *gin.Context) {
	user := tools.GetUser()

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("ws init failed")
		return
	}

	wsConn := &WsConnection{
		ws:        ws,
		userID:    user.ID,
		closeChan: make(chan byte),
		isClosed:  false,
	}

	//	go wsConn.heartbeat()
	go wsConn.readWs()
	go wsConn.writeWs()
	go wsConn.writeBroadCast()
}

func (wsConn *WsConnection) close() {
	wsConn.ws.Close()
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}

func (wsConn *WsConnection) heartbeat() {
	for {
		time.Sleep(2 * time.Second)
		if err := wsConn.ws.WriteMessage(websocket.TextMessage, []byte("hb")); err != nil {
			wsConn.close()
			return
		}
	}
}

func (wsConn *WsConnection) writeWs() {
	if _, ok := MessageQueue[int(wsConn.userID)]; !ok {
		MessageQueue[int(wsConn.userID)] = make(chan *Message, 1000)
	}
	for {
		select {
		case msg := <-MessageQueue[int(wsConn.userID)]:
			user := tools.GetUser()
			if wsConn.userID != user.ID {
				continue
			}
			responseMsg, _ := json.Marshal(msg)
			if err := wsConn.ws.WriteMessage(websocket.TextMessage, []byte(responseMsg)); err != nil {
				fmt.Println("write websocket fail")
				wsConn.close()
				return

			}
		case <-wsConn.closeChan:
			return
		}
	}
}

func (wsConn *WsConnection) writeBroadCast() {
	for {
		select {
		case msg := <-broadcastChan:
			if err := wsConn.ws.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				fmt.Println("write websocket fail")
				wsConn.close()
				return

			}
			if models.CreateMailBox(msg) != nil {

			}
		case <-wsConn.closeChan:
			return
		}
	}
}

func (wsConn *WsConnection) readWs() {
	for {
		_, rawData, err := wsConn.ws.ReadMessage()
		if err != nil {
			wsConn.close()
			return
		}
		data := Message{
			Type:       1,
			FromUserID: wsConn.userID,
		}
		err = json.Unmarshal(rawData, &data)
		if err != nil {
			fmt.Println("json.unmarshal failed")
			fmt.Println("rawData: " + string(rawData))
			continue
		}
		if _, ok := MessageQueue[int(data.ToUserID)]; !ok {
			MessageQueue[int(data.ToUserID)] = make(chan *Message, 1000)
		}
		select {
		case MessageQueue[int(data.ToUserID)] <- &data:
		case <-wsConn.closeChan:
			return
		}
		//TODO: 确认接收
	}
}

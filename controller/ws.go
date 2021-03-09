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
	"healing2020/pkg/tools"
)

type BroadCastContent struct {
	Text string `json:"message"`
}

type Message struct {
	Type       int
	FromuserID int
	TouserID   int    `validate:"required"`
	Content    string `validate:"required"`
}

type WsConnection struct {
	ws        *websocket.Conn
	userID    int
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

func Broadcast(c *gin.Context) {
	json := BroadCastContent{}
	c.BindJSON(&json)
	broadcastChan <- json.Text
}

func WsHandle(c *gin.Context) {
	user := tools.GetUser()

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("ws init failed")
		return
	}

	wsConn := &WsConnection{
		ws:        ws,
		userID:    int(user.ID),
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
	if _, ok := MessageQueue[wsConn.userID]; !ok {
		MessageQueue[wsConn.userID] = make(chan *Message, 1000)
	}
	for {
		select {
		case msg := <-MessageQueue[wsConn.userID]:
			if err := wsConn.ws.WriteMessage(websocket.TextMessage, []byte(msg.Content)); err != nil {
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
			FromuserID: wsConn.userID,
		}
		err = json.Unmarshal(rawData, &data)
		if err != nil {
			fmt.Println("json.unmarshal failed")
			fmt.Println("rawData: " + string(rawData))
			continue
		}
		if _, ok := MessageQueue[data.TouserID]; !ok {
			MessageQueue[data.TouserID] = make(chan *Message, 1000)
		}
		select {
		case MessageQueue[data.TouserID] <- &data:
		case <-wsConn.closeChan:
			return
		}
		//TODO: 确认接收
	}
}

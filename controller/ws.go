package controller

import (
	"crypto/sha1"
	"encoding/json"
	"log"
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

type ACK struct {
	ACKID string `json:"ack"`
}

type Message struct {
	ID         string `json:"id"`
	Type       int    `json:"type"`
	Time       string `json:"time"`
	FromUserID uint   `json:"fromUserID"`
	ToUserID   uint   `json:"toUserID" validate:"required"`
	Content    string `json:"content" validate:"required"`
	URL        string `json:"url"`
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
var ACKchan = make(map[string](chan *ACK))
var MysqlCreate = make(chan *Message, 1000)
var MysqlDelete = make(chan *Message, 1000)

//turn Message to statements.Message
func msgToSMessage(msg *Message) statements.Message {
	message := statements.Message{
		MsgID:   msg.ID,
		Send:    msg.FromUserID,
		Receive: msg.ToUserID,
		Type:    msg.Type,
		Content: msg.Content,
		Url:     msg.URL,
		Time:    msg.Time,
	}
	return message
}

//create msg chan for user
func createUserMsgChan(userID uint) {
	if _, ok := MessageQueue[int(userID)]; !ok {
		MessageQueue[int(userID)] = make(chan *Message, 1000)
	}
}

//run in main.go
//turn message in mysql to chan
func MysqltoChan() {
	allMessage, err := models.SelectAllMessage()
	if err != nil {
		log.Println("get message from mysql fail!")
		return
	}
	for _, value := range allMessage {
		msg := Message{
			ID:         value.MsgID,
			FromUserID: value.Send,
			ToUserID:   value.Receive,
			Type:       value.Type,
			Content:    value.Content,
			URL:        value.Url,
			Time:       value.Time,
		}
		createUserMsgChan(msg.ToUserID)
		MessageQueue[int(msg.ToUserID)] <- &msg
		if msg.Type == 1 {
			createUserMsgChan(msg.FromUserID)
			MessageQueue[int(msg.FromUserID)] <- &msg
		}
	}
}

//@Title Broadcast
//@Description 广播
//@Tags message
//@Produce json
//@Param json body Message true "广播信息(只要content)"
//@Router /broadcast [post]
//@Success 200 {object} e.ErrMsgResponse
//@Failure 403 {object} e.ErrMsgResponse
func Broadcast(c *gin.Context) {
	msg := Message{
		URL:        "",
		Time:       time.Now().Format("2006-01-02 15:04:05"),
		Type:       0,
		FromUserID: 0,
	}
	c.BindJSON(&msg)
	hash, ok := c.Get("hash")
	if !ok || hash != getHash() {
		c.JSON(403, e.ErrMsgResponse{Message: "rejected"})
		return
	}
	msg.ID = tools.Md5String(msg.Time)
	userCount, err := models.GetUserNum()
	for i := 1; i <= userCount; i++ {
		createUserMsgChan(uint(i))
		msg.ToUserID = uint(i)
		MessageQueue[i] <- &msg
	}
	err = models.CreateMailBox(msg.Content)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "存储广播信息失败"})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: e.GetMsg(e.SUCCESS)})
}

//处理ws连接
func WsHandle(c *gin.Context) {
	user := tools.GetUser(c)

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("ws init failed")
		return
	}

	wsConn := &WsConnection{
		ws:        ws,
		userID:    user.ID,
		closeChan: make(chan byte),
		isClosed:  false,
	}

	go wsConn.heartbeat()
	go wsConn.readWs(c)
	go wsConn.writeWs(c)
	go wsConn.MsgMysql()
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

//心跳
func (wsConn *WsConnection) heartbeat() {
	for {
		time.Sleep(2 * time.Second)
		if err := wsConn.ws.WriteMessage(websocket.TextMessage, []byte("hb")); err != nil {
			wsConn.close()
			return
		}
	}
}

//向客户端发送用户消息
func (wsConn *WsConnection) writeWs(c *gin.Context) {
	userID := tools.GetUser(c).ID
	createUserMsgChan(userID)
	timeoutNum := 0
	for {
		select {
		case msg := <-MessageQueue[int(wsConn.userID)]:
			if wsConn.userID != userID {
				continue
			}
			//wait 0.5s to response ack to front-end
			time.Sleep(time.Duration(500) * time.Millisecond)
			//send message
			responseMsg, _ := json.Marshal(msg)
			if err := wsConn.ws.WriteMessage(websocket.TextMessage, []byte(responseMsg)); err != nil {
				log.Println("write websocket fail")
				wsConn.close()
				return
			}
			//create ack chan for every message
			if _, ok := ACKchan[msg.ID]; !ok {
				ACKchan[msg.ID] = make(chan *ACK)
			}
			select {
			//if timeout 2s, drop msg back to the message chan
			case <-time.After(time.Second * 2):
				log.Println("timeout, msg is not be received")
				timeoutNum += 1
				MessageQueue[int(wsConn.userID)] <- msg
				//if no response from front-end for long time, close ws
				if timeoutNum > 10 {
					log.Println("no response from front-end for long time, ws close")
					timeoutNum = 0
					wsConn.close()
					return
				}
			case ack := <-ACKchan[msg.ID]:
				//judge ack
				if ack.ACKID == msg.ID {
					//log.Println("he get it")
					MysqlDelete <- msg
					//if get ack, resetting timeoutNum
					timeoutNum = 0
					continue
				} else {
					log.Println("塞进该消息ack通道的ackID与消息不符合！")
				}
			}

		case <-wsConn.closeChan:
			return
		}
	}
}

//接收客户端发送的消息和ack
func (wsConn *WsConnection) readWs(c *gin.Context) {
	userID := tools.GetUser(c).ID
	for {
		_, rawData, err := wsConn.ws.ReadMessage()
		if err != nil {
			log.Println("read ws fail, ws close")
			wsConn.close()
			return
		}

		receiveACK := ACK{}
		data := Message{}

		json.Unmarshal(rawData, &receiveACK)
		json.Unmarshal(rawData, &data)
		//log.Println(data)
		//log.Println(receiveACK)

		if receiveACK != (ACK{}) {
			if _, ok := ACKchan[receiveACK.ACKID]; !ok {
				log.Println("未知的ack报文")
				wsConn.ws.WriteMessage(websocket.TextMessage, []byte("未知的ack报文"))
				continue
			}
			ACKchan[receiveACK.ACKID] <- &receiveACK
		} else if data != (Message{}) && data.Type == 2 {
			//judge data FromUserID
			if userID != data.FromUserID {
				wsConn.ws.WriteMessage(websocket.TextMessage, []byte("FromUserID和用户id不同"))
				log.Println("FromUserID is not same as userID")
				data.FromUserID = userID
			}
			MysqlCreate <- &data
			//if get message, response ack
			responseACK, _ := json.Marshal(ACK{ACKID: data.ID})
			wsConn.ws.WriteMessage(websocket.TextMessage, []byte(responseACK))
			//ready to response msg
			createUserMsgChan(data.ToUserID)
			select {
			case MessageQueue[int(data.ToUserID)] <- &data:
			case <-wsConn.closeChan:
				return
			}
		} else {
			log.Println("json.unmarshal failed")
			wsConn.ws.WriteMessage(websocket.TextMessage, []byte("json.unmarshal failed"))
			log.Println("rawData: " + string(rawData))
			continue
		}
	}
}

//持久化存储
func (wsConn *WsConnection) MsgMysql() {
	mysqlError := 0
	for {
		select {
		case msg := <-MysqlCreate:
			message := msgToSMessage(msg)
			if models.SaveMessage(message) != nil {
				MysqlCreate <- msg
				mysqlError += 1
			} else {
				mysqlError = 0
			}
		case msg := <-MysqlDelete:
			message := msgToSMessage(msg)
			if models.DeleteMessage(message) != nil {
				MysqlDelete <- msg
				mysqlError += 1
			} else {
				mysqlError = 0
			}
		}
		//if try to save or create mant times, close ws
		if mysqlError > 10 {
			log.Println("can not connect to the mysql, ws close")
			wsConn.close()
		}
	}
}

func getHash() string {
	t := time.Now()
	t_zero := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Format("2006-01-02 15:04")
	sh := sha1.New()
	sh.Write([]byte(t_zero))
	sum := sh.Sum([]byte("healing2020"))
	return string(sum)
}

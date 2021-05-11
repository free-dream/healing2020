package controller

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

var mutex sync.Mutex
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
	mutex.Lock()
	if _, ok := MessageQueue[int(userID)]; !ok {
		MessageQueue[int(userID)] = make(chan *Message, 1000)
	}
	mutex.Unlock()
}

//create ack chan for message
func createACKchan(msgID string) {
	mutex.Lock()
	if _, ok := ACKchan[msgID]; !ok {
		ACKchan[msgID] = make(chan *ACK)
	}
	mutex.Unlock()
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

//持久化存储
func (wsConn *WsConnection) MsgMysql() {
	mysqlError := 0
	for {
		select {
		case msg := <-MysqlCreate:
			message := msgToSMessage(msg)
			err := models.SaveMessage(message)
			if err != nil {
				MysqlCreate <- msg
				mysqlError += 1
			} else {
				mysqlError = 0
			}
		case msg := <-MysqlDelete:
			message := msgToSMessage(msg)
			err := models.DeleteMessage(message)
			if err != nil {
				MysqlDelete <- msg
				mysqlError += 1
			} else {
				mysqlError = 0
			}
		case <-wsConn.closeChan:
			return
		}
		//if try to save or create mant times, close ws
		if mysqlError > 10 {
			log.Println("ws save fail many times, can not connect to the mysql, ws close")
			wsConn.close()
			return
		}

	}
}

type BroadcastReq struct {
	Content string `json:"content" binding:"required"`
	Hash    string `json:"hash" binding:"required"`
	Type    int    `json:"type"`
	Url     string `json:"url"`
}

//@Title Broadcast
//@Description 广播
//@Tags message
//@Produce json
//@Router /broadcast [post]
//@Success 200 {object} e.ErrMsgResponse
//@Failure 403 {object} e.ErrMsgResponse
func Broadcast(c *gin.Context) {
	msg := Message{
		Time:       time.Now().Format("2006-01-02 15:04:05"),
		FromUserID: 0,
	}

	// 鉴权
	var form BroadcastReq
	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: e.GetMsg(e.INVALID_PARAMS)})
		return
	}
	if form.Hash != getHash() {
		c.JSON(403, e.ErrMsgResponse{Message: "rejected"})
		return
	}
	msg.Content = form.Content
	msg.Type = form.Type
	msg.URL = form.Url
	msg.ID = tools.Md5String(msg.Time)

	// 开始广播
	userCount, err := models.GetUserNum()
	for i := 1; i <= userCount; i++ {
		msg.ToUserID = uint(i)
		MysqlCreate <- &msg
		createUserMsgChan(uint(i))
		MessageQueue[i] <- &msg
	}
	err = models.CreateMailBox(msg.Content)
	if err != nil {
		c.JSON(403, e.ErrMsgResponse{Message: "存储广播信息失败"})
		return
	}
	c.JSON(200, e.ErrMsgResponse{Message: e.GetMsg(e.SUCCESS)})
}

//处理ws连接，http协议升级
func WsHandle(c *gin.Context) {
	user := tools.GetUser(c)

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(strconv.Itoa(int(user.ID)) + "ws init failed")
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
				log.Println("write websocket fail " + strconv.Itoa(int(userID)))
				wsConn.close()
				return
			}
			createACKchan(msg.ID)
			select {
			//if timeout 2s, drop msg back to the message chan
			case <-time.After(time.Second * 2):
				timeoutNum += 1
				createUserMsgChan(wsConn.userID)
				MessageQueue[int(wsConn.userID)] <- msg
				//if no response from front-end for long time, close ws
				if timeoutNum > 10 {
					log.Println(strconv.Itoa(int(wsConn.userID)) + " no response from front-end for long time, ws close")
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
					log.Println("bad ackID:" + ack.ACKID + " 塞进该消息ack通道的ackID与消息不符合！")
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

		select {
		case <-wsConn.closeChan:
			return
		default:
		}
		_, rawData, err := wsConn.ws.ReadMessage()
		if err != nil {
			wsConn.close()
			return
		}

		receiveACK := ACK{}
		data := Message{}

		json.Unmarshal(rawData, &receiveACK)
		json.Unmarshal(rawData, &data)
		//log.Println(data)
		//log.Println(receiveACK)

		if receiveACK != (ACK{}) { //receive ack
			if _, ok := ACKchan[receiveACK.ACKID]; !ok {
				log.Println("from:" + strconv.Itoa(int(userID)) + "未知的ack报文")
				wsConn.ws.WriteMessage(websocket.TextMessage, []byte("未知的ack报文"))
				continue
			}
			ACKchan[receiveACK.ACKID] <- &receiveACK
		} else if data != (Message{}) && data.Type == 2 { //receive msg data
			//judge data FromUserID
			if userID != data.FromUserID {
				wsConn.ws.WriteMessage(websocket.TextMessage, []byte("FromUserID和用户id不同"))
				log.Println("FromUserID " + strconv.Itoa(int(data.FromUserID)) + "is not same as userID " + strconv.Itoa(int(userID)))
				data.FromUserID = userID
			}
			MysqlCreate <- &data
			//if get message, response ack
			responseACK, _ := json.Marshal(ACK{ACKID: data.ID})
			wsConn.ws.WriteMessage(websocket.TextMessage, []byte(responseACK))
			//ready to response msg
			createUserMsgChan(data.ToUserID)
			MessageQueue[int(data.ToUserID)] <- &data
		} else { //not ack or msg
			log.Println("ws json.unmarshal failed")
			wsConn.ws.WriteMessage(websocket.TextMessage, []byte("json.unmarshal failed"))
			log.Println("rawData: " + string(rawData))
			continue
		}

	}
}

func getHash() string {
	t := time.Now()
	t_zero := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Format("2006-01-02 15:04")
	sh := sha1.New()
	sh.Write([]byte(t_zero))
	sum := sh.Sum([]byte("healing2020"))
	return fmt.Sprintf("%x", sum)
}

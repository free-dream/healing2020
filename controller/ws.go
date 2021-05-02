package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"healing2020/models"
	"healing2020/pkg/e"
	"healing2020/pkg/tools"
)

type ACK struct {
	ACKID string `json:"ack"`
}

type ServerMsg struct {
	Text string `json:"message"`
	Time string `json:"time"`
	Type int    `json:"type"`
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
var ServerMsgChan = make(chan *ServerMsg)
var ACKchan = make(map[string](chan *ACK))

//@Title Broadcast
//@Description 广播
//@Tags message
//@Produce json
//@Param json body ServerMsg true "广播信息"
//@Router /broadcast [post]
//@Success 200 {object} e.ErrMsgResponse
//@Failure 403 {object} e.ErrMsgResponse
func Broadcast(c *gin.Context) {
	json := ServerMsg{
		Type: 0,
	}
	c.BindJSON(&json)

	userCount, err := models.GetUserNum()
	for i := 0; i < userCount; i++ {
		ServerMsgChan <- &json
	}

	err = models.CreateMailBox(json.Text)
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
	go wsConn.readWs()
	go wsConn.writeWs(c)
	go wsConn.writeServerMsg(c)
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
	if _, ok := MessageQueue[int(wsConn.userID)]; !ok {
		MessageQueue[int(wsConn.userID)] = make(chan *Message, 1000)
	}
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
func (wsConn *WsConnection) readWs() {
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
				continue
			}
			ACKchan[receiveACK.ACKID] <- &receiveACK
		} else if data != (Message{}) {
			//if get message, response ack
			responseACK, _ := json.Marshal(ACK{ACKID: data.ID})
			wsConn.ws.WriteMessage(websocket.TextMessage, []byte(responseACK))
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

//广播
func (wsConn *WsConnection) writeServerMsg(c *gin.Context) {
	for {
		select {
		case msg := <-ServerMsgChan:
			//send broadcast msg
			responseMsg, _ := json.Marshal(msg)
			if err := wsConn.ws.WriteMessage(websocket.TextMessage, []byte(responseMsg)); err != nil {
				log.Println("write websocket fail")
				wsConn.close()
				return
			}
			//save the broadcast inf
			if models.CreateMailBox(msg.Text) != nil {
				log.Println("create mailbox in mysql fail")
			}
		case <-wsConn.closeChan:
			return
		}
	}
}

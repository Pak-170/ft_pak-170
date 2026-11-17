package server

import (
	"context"
	"go_pakito/dbp"
	"go_pakito/models"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub		*Hub
	ws		*websocket.Conn
	info	*models.User
	online	bool

	sendMsg			chan *models.ChatView
	sendChatInfo	chan *models.Chat
	cancel			context.CancelFunc
}

func newClient(h *Hub, ws *websocket.Conn, info *models.User) *Client {
	return &Client{
		hub:			h,
		ws:				ws,
		info:			info,
		online:			false,

		sendMsg:		make(chan *models.ChatView),
		sendChatInfo:	make(chan *models.Chat),
	}
}


// ----------------------------------------------------------------------------


func isValidMsg(msg *models.ChatView, err *error) (int, string){
	if err != nil {
		if websocket.IsCloseError(*err, websocket.CloseGoingAway) {
			return websocket.CloseNormalClosure, "Bye!"
		} else {
			return websocket.CloseInvalidFramePayloadData, (*err).Error()
		}
	}

	if validate.Struct(msg) != nil {
		return websocket.ClosePolicyViolation, "Invalid JSON"
	}
	return 0, ""
}

func (c *Client) readFromWS(ctx context.Context) {

	defer c.ws.Close()
	defer c.cancel()
	defer func() {c.hub.leave <- c}()

	for {
		var msg models.ChatView
		err := c.ws.ReadJSON(&msg)
		errCode, errMsg := isValidMsg(&msg, &err)
		if errCode != 0 {
			log.Printf("Client %v (%v) disconnected with Code %v (%v)\n",
				c.info.Username, c.info.ID, errCode, errMsg)
			c.ws.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(errCode, string(errMsg)),
				time.Now().Add(1 * time.Second))
			break
		}
		select {
		case <-ctx.Done():
			log.Println("CTX DONE EN READPUMP")
			return
		case c.hub.broadcastChan <- &msg:
			err = dbp.InsertMsg(ctx, c.hub.DB, &msg)
			if err != nil {
				log.Printf("ERROR INSERTING INTO DB: %v\n", err)
				return
			}
		}
	}
}

func (c *Client) writeToWS(ctx context.Context) {

	defer c.ws.Close()
	defer c.cancel()

	for {
		select {
		case <-ctx.Done():
			log.Println("CTX DONE EN WRITEPUMP")
			return
		case msg, ok := <-c.sendMsg:
			if !ok {
				log.Println("User connection closed")
				return
			}
			err := c.ws.WriteJSON(msg)
			if err != nil {
				log.Printf("WS error msg: %v\n", err)
				return
			}
		case chatInfo, ok := <-c.sendChatInfo:
			if !ok {
				log.Println("User connection closed")
				return
			}
			err := c.ws.WriteJSON(chatInfo)
			if err != nil {
				log.Printf("WS error chat: %v\n", err)
				return
			}
		}
	}
}

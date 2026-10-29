package server

import (
	"context"
	"go_pakito/dbp"
	"go_pakito/models"
	"log"

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

func (c *Client) readFromWS(ctx context.Context) {

	defer c.ws.Close()
	defer c.cancel()

	var msg models.ChatView
	for {
		err := c.ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("ERROR READING JSON FROM WS: %v\n", err)
			break
		}
		c.hub.broadcastChan <- &msg
		dbp.InsertMsg(ctx, c.hub.DB, &msg)
	}
}

func (c *Client) writeToWS(ctx context.Context) {

//	defer c.ws.Close()
	defer func() {c.hub.leave <- c}()

	for {
		select {
		case <-ctx.Done():
			log.Println("CTX DONE EN WRITEPUMP")
			return
		case msg, ok := <-c.sendMsg: // hemos recibimos un msg en el canal send
			log.Printf("(%v) 1 ->> OK: %v\n", c.info.ID, ok)
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
			log.Printf("(%v) 2 ->> OK: %v\n", c.info.ID, ok)
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

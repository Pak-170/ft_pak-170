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
	notification	chan *any
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

	var msg models.ChatView
	for {
		err := c.ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("ERROR READING JSON FROM WS: %v\n", err)
			break
		}
		log.Printf("(%v)->>>> msg received: %v.\n", c.info.ID, msg)
		c.hub.broadcastChan <- &msg
		dbp.InsertMsg(ctx, c.hub.DB, &msg)
	}
}

func (c *Client) writeToWS(ctx context.Context) {

	defer c.ws.Close()
	defer func() {c.hub.leave <- c}()

	for { //? Falta el ticker, comprobar si es necesario.
		select {
		case <-ctx.Done():
			log.Println("CTX DONE EN WRITEPUMP")
			return
		case msg, ok := <-c.sendMsg: // hemos recibimos un msg en el canal send
			if !ok {
				log.Println("User closed connection")
				return
			}
			err := c.ws.WriteJSON(msg)
			if err != nil {
				log.Printf("WS error msg: %v\n", err)
				return
			}
		case chatInfo := <-c.sendChatInfo:
			log.Printf("WS recibido para un nuevo chat en el cliente %v: %v\n", c.info.ID, chatInfo)
			err := c.ws.WriteJSON(chatInfo)
			if err != nil {
				log.Printf("WS error chat: %v\n", err)
				return
			}
		}
	}
}

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

	send chan *models.ChatView
}

func newClient(h *Hub, ws *websocket.Conn, info *models.User) Client {
	return Client{
		hub:	h,
		ws:		ws,
		info:	info,
		send:	make(chan *models.ChatView),
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
			//break
		}
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
		case msg, ok := <-c.send: // hemos recibimos un msg en el canal send
			if !ok {
				log.Println("User closed connection")
				return
			}
			err := c.ws.WriteJSON(msg)
			if err != nil {
				log.Printf("WS error: %v\n", err)
				return
			}
		}
	}
}

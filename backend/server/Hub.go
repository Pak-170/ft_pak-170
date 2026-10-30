package server

import (
	"context"
	"go_pakito/dbp"
	"go_pakito/models"
	"log"

	"gorm.io/gorm"
)

type Hub struct { // Es el único que puede gestionar y trabajar realmente con los elementos compartidos entre clientes
	DB				*gorm.DB
	chats			map[int][]*Client

	join			chan *Client
	leave			chan *Client
	broadcastChan	chan *models.ChatView
}

func newHub(db *gorm.DB) Hub {
	return Hub{
		DB:				db,
		chats:			make(map[int][]*Client),
		join:			make(chan *Client),
		leave:			make(chan *Client),
		broadcastChan:	make(chan *models.ChatView),
	}
}

// ----------------------------------------------------------------------------

func (h *Hub) removeUser(ctx context.Context, user *Client) {

	chatIds, err := dbp.SelectChatIds(ctx, h.DB, int(user.info.ID))

	if err != nil {
		log.Println(err)
		return
	}

	for _, chatID := range chatIds {
		for i, userInChat := range h.chats[chatID] {
			if userInChat == user {
				h.chats[chatID] = append(h.chats[chatID][:i], h.chats[chatID][(i+1):]...)
				break
			}
		}
	}
}

func (h *Hub) addUser(ctx context.Context, user *Client) {

	chatIds, err := dbp.SelectChatIds(ctx, h.DB, int(user.info.ID))

	if err != nil {
		log.Println(err)
		return
	}

	for _, chatID := range chatIds {
		h.chats[chatID] = append(h.chats[chatID], user)
	}
}

func (h *Hub) broadcast(ctx context.Context, msg *models.ChatView) {

	for _, user := range h.chats[int(msg.ChatID)] {
		if user.info.Username == msg.Username {
			continue
		}
		select {
		case user.send <- msg: // escribimos en el canal del usuario para los msgs, no mandamos el mensaje por ws.
		default:
			close(user.send)
			h.removeUser(ctx, user)
		}
	}
}

func initChats(hub *Hub) error {

	var chats map[int][]*Client = make(map[int][]*Client, 0)

	chatIds, err := dbp.SelectChatIds(context.Background(), hub.DB, dbp.ALL) //!

	if err != nil {
		return err
	}
	for _, chatID := range chatIds {
		chats[chatID] = make([]*Client, 0)
	}
	hub.chats = chats
	return err
}

// ----------------------------------------------------------------------------

// *{"ChatID":1,"Username":"yuri","Msg":"sadasdasd","SendTime":"2026-10-19T08:47:38.363Z"}
func (h *Hub) run(ctx context.Context) {

	if err := initChats(h); err != nil {
		log.Printf("ERROR INITIALIZING CHATS: %v\n", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("conexion closed")
			return
		case user := <-h.join:
			h.addUser(ctx, user)
		case user := <-h.leave:
			h.removeUser(ctx, user)
		case msg := <-h.broadcastChan:
			log.Println("MESNAJE RECIBIDO PARA MANDAR POR BROADCAAAAAST")
			h.broadcast(ctx, msg)
		}
	}
}

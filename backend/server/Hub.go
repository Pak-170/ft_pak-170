package server

import (
	"context"
	"go_pakito/dbp"
	"go_pakito/models"
	"log"
	"gorm.io/gorm"
)


type CreateChatData struct {
	InfoType	string	`json:"infoType"`
	Name		string	`json:"name"`
	ChatType	string	`json:"type"`
	UserIDs		[]uint	`json:"userIDs"`
}

type Hub struct { // Es el único que puede gestionar y trabajar realmente con los elementos compartidos entre clientes
	DB				*gorm.DB
	chats			map[uint]*models.Chat
	clients			map[uint]*Client

	join			chan *Client
	leave			chan *Client
	newChat			chan *models.Chat
	broadcastChan	chan *models.ChatView
}

func newHub(db *gorm.DB) Hub {
	return Hub{
		DB:				db,
		chats:			make(map[uint]*models.Chat),
		clients:		make(map[uint]*Client),
		join:			make(chan *Client),
		leave:			make(chan *Client),
		newChat:		make(chan *models.Chat),
		broadcastChan:	make(chan *models.ChatView),
	}
}

// ----------------------------------------------------------------------------


func (h *Hub) broadcast(msg *models.ChatView) {

	newMsgNotification := newMsgNotification{msg}

	for _, user := range h.chats[msg.ChatID].Participants {
		client := h.clients[user.ID]
		if user.Username == msg.Username {
			continue
		}
		newMsgNotification.notify(client)
	}
}


func (hub *Hub)initChats(ctx context.Context) error {

	var chats map[uint]*models.Chat = make(map[uint]*models.Chat, 0)

	chatsDB, err := dbp.SelectAny[models.Chat](ctx, hub.DB)

	if err != nil {
		return err
	}
	for _, chatDB := range chatsDB {
		participants, err := dbp.SelectUsersOfChat(ctx, hub.DB, chatDB.ID)
		if err != nil {
			return err
		}
		chats[chatDB.ID] = &chatDB
		chats[chatDB.ID].Participants = participants
	}
	hub.chats = chats
	return err
}

func (hub *Hub)initClients(ctx context.Context) error {
	var clients map[uint]*Client = make(map[uint]*Client, 0)

	clientsDB, err := dbp.SelectAny[models.User](ctx, hub.DB)

	if err != nil {
		return err
	}
	for _, clientDB := range clientsDB {
		clients[clientDB.ID] = newClient(hub, nil, &clientDB)
	}
	hub.clients = clients
	return err
}


func	(h *Hub) insertUsersInChat(ctx context.Context, chatID uint, users []uint) error {

	var err error

	for _, user := range users {
		err = dbp.InsertUserInChat(ctx, h.DB, chatID, user)
		if err != nil {
			return err
		}
	}

	participants, err := dbp.SelectUsersOfChat(ctx, h.DB, chatID)
	if err != nil {
		return err
	}
	h.chats[chatID].Participants = participants

	return err
}

// ----------------------------------------------------------------------------

func (h *Hub) run(ctx context.Context) {

	if err := h.initChats(ctx); err != nil {
		log.Printf("ERROR INITIALIZING CHATS: %v\n", err)
		return
	}
	if err := h.initClients(ctx); err != nil {
		log.Printf("ERROR INITIALIZING CLIENTS: %v\n", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("conexion closed")
			return
		case user := <-h.join:
			log.Printf("hub join channel: %v\n", user)
			h.clients[user.info.ID].online = true
		case user := <-h.leave:
			log.Printf("hub leave channel: %v\n", user)
			h.clients[user.info.ID].online = false
		case msg := <-h.broadcastChan:
			log.Printf("hub broadcast channel: %v\n", msg)
			h.broadcast(msg)
		}
	}
}

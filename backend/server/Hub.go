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
	chats			map[uint][]*Client
	clients			map[uint]*Client

	join			chan *Client
	leave			chan *Client
	newChat			chan *models.Chat
	broadcastChan	chan *models.ChatView
}

func newHub(db *gorm.DB) Hub {
	return Hub{
		DB:				db,
		chats:			make(map[uint][]*Client),
		clients:		make(map[uint]*Client),
		join:			make(chan *Client),
		leave:			make(chan *Client),
		newChat:		make(chan *models.Chat),
		broadcastChan:	make(chan *models.ChatView),
	}
}

// ----------------------------------------------------------------------------

func (h *Hub) removeClient(ctx context.Context, client *Client) {

	chatIds, err := dbp.SelectChatIds(ctx, h.DB, int(client.info.ID))

	if err != nil {
		log.Println(err)
		return
	}

	for _, chatID := range chatIds {
		for i, clientInChat := range h.chats[chatID] {
			if clientInChat == client {
				h.chats[chatID] = append(h.chats[chatID][:i], h.chats[chatID][(i+1):]...)
				break
			}
		}
	}
}

func (h *Hub) addClient(ctx context.Context, client *Client) {

	h.clients[client.info.ID] = client
	chatIds, err := dbp.SelectChatIds(ctx, h.DB, int(client.info.ID))

	if err != nil {
		log.Println(err)
		return
	}

	for _, chatID := range chatIds {
		h.chats[chatID] = append(h.chats[chatID], client)
	}
}

func (h *Hub) broadcast(ctx context.Context, msg *models.ChatView) {

	for _, client := range h.chats[msg.ChatID] {
		if client.info.Username == msg.Username {
			continue
		}
		select {
		case client.sendMsg <- msg: // escribimos en el canal del usuario para los msgs, no mandamos el mensaje por ws.
		default:
			close(client.sendMsg)
			close(client.sendChatInfo)
			h.removeClient(ctx, client)
		}
	}
}


func (hub *Hub)initChats() error {

	var chats map[uint][]*Client = make(map[uint][]*Client, 0)

	chatIds, err := dbp.SelectChatIds(context.Background(), hub.DB, dbp.ALL)

	if err != nil {
		return err
	}
	for _, chatID := range chatIds {
		chats[chatID] = make([]*Client, 0)
	}
	hub.chats = chats
	return err
}

func (hub *Hub)initClients() error {
	var clients map[uint]*Client = make(map[uint]*Client, 0)

	clientsDB, err := dbp.SelectAny[models.User](context.Background(), hub.DB)

	if err != nil {
		return err
	}
	for _, clientDB := range clientsDB {
		clients[clientDB.ID] = newClient(hub, nil, &clientDB)
	}
	hub.clients = clients
	return err
}

// ----------------------------------------------------------------------------

func (h *Hub) run(ctx context.Context) {

	if err := h.initChats(); err != nil {
		log.Printf("ERROR INITIALIZING CHATS: %v\n", err)
		return
	}
	if err := h.initClients(); err != nil {
		log.Printf("ERROR INITIALIZING CLIENTS: %v\n", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("conexion closed")
			return
		case user := <-h.join:
			h.addClient(ctx, user)
		case user := <-h.leave:
			h.removeClient(ctx, user)
		case msg := <-h.broadcastChan:
			h.broadcast(ctx, msg)
		}
	}
}

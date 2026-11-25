package server

import (
	"context"
	"go_pakito/Notification"
	"go_pakito/Session"
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
	//clients			map[uint]*Client
	sessions		map[string]session.Session

	join			chan *session.Session
	leave			chan *session.Session
	newChat			chan *models.Chat
	broadcastChan	chan *models.ChatView
}

func newHub(db *gorm.DB) Hub {
	return Hub{
		DB:				db,
		chats:			make(map[uint]*models.Chat),
		//clients:		make(map[uint]*Client),
		sessions:		make(map[string]session.Session),

		join:			make(chan *session.Session),
		leave:			make(chan *session.Session),
		newChat:		make(chan *models.Chat),
		broadcastChan:	make(chan *models.ChatView),
	}
}

// ----------------------------------------------------------------------------

// -- Session interface implementation:

func (h *Hub) Broadcast(ctx context.Context, msg *models.ChatView) error{
	select {
	case <-ctx.Done():
		log.Println("CTX DONE EN READPUMP - BROADCAST")
		return ctx.Err()
	case h.broadcastChan <- msg:
		return nil
	}
}

func (h *Hub) Join(s *session.Session) {
	h.join <- s // qUIZÁS esto en un select
}

func (h *Hub) Leave(s *session.Session) {
	h.leave <- s
}

// ----------------------------------------------------------------------------

func (h *Hub) broadcast(ctx context.Context, msg *models.ChatView) {

	newMsg := Notification.NewMsg{Msg: msg}

	sessions, err := dbp.SelectSessionsFromChatID(ctx, h.DB, msg.ChatID)
	if err != nil {
		log.Println("DB error while trying to broadcast a msg...")
	}

	for _, sessionModel := range sessions {

		sessionSession, exists := h.sessions[sessionModel.ID] // Qué clase de nombres son estos Yolanda por favor
		if exists == false || sessionSession.GetHashedToken() == msg.SenderHashToken {
			continue
		}
		newMsg.Notify(&sessionSession)
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
	// if err := h.initClients(ctx); err != nil {
	// 	log.Printf("ERROR INITIALIZING CLIENTS: %v\n", err)
	// 	return
	// }

	for {
		select {
		case <-ctx.Done():
			log.Println("conection closed --------------- ")
			return
		case s := <-h.join:
			h.sessions[s.GetHashedToken()] = *s
		case s := <-h.leave:
			delete(h.sessions, s.GetHashedToken())
		case msg := <-h.broadcastChan:
			h.broadcast(ctx, msg)
		}
	}
}

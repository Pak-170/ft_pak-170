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
	sessions		map[string]session.Session

	join			chan *session.Session
	leave			chan *session.Session
	newChat			chan *models.Chat
	broadcastChan	chan *models.ChatView
}

func newHub(db *gorm.DB) Hub {
	return Hub{
		DB:				db,
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

func (h *Hub) Join(ctx context.Context, s *session.Session) error{
	select {
	case <-ctx.Done():
		return ctx.Err()
	case h.join <- s:
		return nil
	}
}

func (h *Hub) Leave(ctx context.Context, s *session.Session) error{
	select {
	case <-ctx.Done():
		return ctx.Err()
	case h.leave <- s:
		return nil
	}
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


func	(h *Hub) insertUsersInChat(ctx context.Context, chatID uint, users []uint) error {

	var err error

	for _, user := range users {
		err = dbp.InsertUserInChat(ctx, h.DB, chatID, user)
		if err != nil {
			return err
		}
	}
	return err
}

// ----------------------------------------------------------------------------

func (h *Hub) run(ctx context.Context) {

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

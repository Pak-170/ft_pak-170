package session

import (
	"context"
	"github.com/gorilla/websocket"
	"go_pakito/models"
	"go_pakito/dbp"
	"gorm.io/gorm"
	"go_pakito/validation"
	"log"
	"time"
)

type hub interface {
	Broadcast(ctx context.Context, msg *models.ChatView) error
	Join(ctx context.Context, s *Session) error
	Leave(ctx context.Context, s *Session) error
}

// mapa token-session <-> Session.
type Session struct {
	token	string
	ws		*websocket.Conn
	Info	models.User // Esto quizás que sea privado luego

	SendMsg			chan *models.ChatView
	SendChatInfo	chan *models.Chat
	cancel			context.CancelFunc
}

func (s *Session)GetHashedToken() string {
	return s.token
}


func NewSession(Token string, WS *websocket.Conn, Info *models.User, Cancel *context.CancelFunc) Session {
	return Session{
		token: Token,
		ws: WS,
		Info: *Info,

		SendMsg: make(chan *models.ChatView),
		SendChatInfo: make(chan *models.Chat),
		cancel: *Cancel,
	}
}

func (s *Session)CloseSession() {
	s.cancel()
	s.ws.Close()
}


func (s *Session)ReadFromWS(ctx context.Context, h hub, db *gorm.DB) {

	defer s.CloseSession()
	defer h.Leave(ctx, s)

	for {
		var msg models.ChatView
		err := s.ws.ReadJSON(&msg)
		errCode, errMsg := validation.IsValidMsg(&msg, err)
		if errCode != 0 {
			log.Printf("Client %v (%v) disconnected with Code %v (%v)\n",
				s.Info.Username, s.Info.ID, errCode, errMsg)
			s.ws.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(errCode, string(errMsg)),
				time.Now().Add(1 * time.Second))
			break
		}
		err = dbp.InsertMsg(ctx, db, &msg)
		if err != nil {
			log.Printf("ERROR INSERTING INTO DB: %v\n", err)
			break
		}
		msg.SenderHashToken = s.token
		err = h.Broadcast(ctx, &msg)
		if err != nil {
			log.Printf("ERROR BROADCASTING: %v\n", err)
			break
		}
	}
}

func (s *Session)WriteToWS(ctx context.Context) {

	defer s.CloseSession()

	for {
		select {
		case <-ctx.Done():
			log.Println("CTX DONE EN WRITEPUMP")
			return
		case msg, ok := <-s.SendMsg:
			if !ok {
				log.Println("User connection closed")
				return
			}
			err := s.ws.WriteJSON(msg)
			if err != nil {
				log.Printf("WS error msg: %v\n", err)
				return
			}
		case chatInfo, ok := <-s.SendChatInfo:
			if !ok {
				log.Println("User connection closed")
				return
			}
			err := s.ws.WriteJSON(chatInfo)
			if err != nil {
				log.Printf("WS error chat: %v\n", err)
				return
			}
		}
	}
}
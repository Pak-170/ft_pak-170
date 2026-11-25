package Notification

import (
	"go_pakito/Session"
	"go_pakito/models"
)

/*
	- Mandar actualización por ws [x]
	- Mandar Notificación por ws
*/

type notification interface{
	Notify(s *session.Session)
}


type NewChat struct{
	Chat	*models.Chat
}

func (n *NewChat)Notify(s *session.Session) {
	select {
	case s.SendChatInfo <- n.Chat:
	default:
		s.CloseSession()
	}
}


type NewMsg struct{
	Msg	*models.ChatView
}
func (n *NewMsg)Notify(s *session.Session) {

	select {
	case s.SendMsg <- n.Msg:
	default:
		s.CloseSession()
	}
}



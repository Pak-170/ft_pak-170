package Notification

import (
	"go_pakito/Session"
	"go_pakito/models"
	"log"
)

/*
	- Mandar actualización por ws [x]
	- Mandar Notificación por ws
*/

type notification interface{
	Notify([]uint)
}


type NewChat struct{
	Chat	*models.Chat
}

func (n *NewChat)Notify(s *session.Session) {
	select {
	case s.SendChatInfo <- n.Chat:
		log.Println("NOTIFY ENVIAMOS A UN CANAAAAAAAAAAAAAAAAAAAAAAAAAAAL")
	default:
		log.Println("NOTIFY default ;( why")
		s.CloseSession()
	}
}


type NewMsg struct{
	Msg	*models.ChatView
}
func (n *NewMsg)Notify(s *session.Session) {

	select {
	case s.SendMsg <- n.Msg:
		log.Println("NOTIFY ENVIAMOS A UN CANAAAAAAAAAAAAAAAAAAAAAAAAAAAL")
	default:
		log.Println("NOTIFY default ;( why")
		s.CloseSession()
	}
}



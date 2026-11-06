package server

import (
	// "context"
	// "go_pakito/dbp"
	"go_pakito/models"
	"log"
)

const (
	MSG = 1
	NEWCHAT = 2
)

type Notification interface{
	notify(*Client)
}

type newChatNotification struct{
	chat	*models.Chat
}
func (n *newChatNotification)notify(c *Client) {

	if (c.online == false) {
		return
	}
	log.Printf("***** newChatNotification for client %v: %v\n", c.info.ID, n)
	c.sendChatInfo <- n.chat
}



// type newChat struct {
// 	n	Notification[T]
// }


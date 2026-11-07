package server

import (
	"go_pakito/models"
	"log"
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


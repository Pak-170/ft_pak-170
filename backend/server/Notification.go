package server

import (
	"go_pakito/models"
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

	select {
	case c.sendChatInfo <- n.chat:
	default:
		close(c.sendMsg)
		close(c.sendChatInfo)
		c.online = false
	}
}


type newMsgNotification struct{
	msg	*models.ChatView
}
func (n *newMsgNotification)notify(c *Client) {

	if (c.online == false) {
		return
	}

	select {
	case c.sendMsg <- n.msg:
	default:
		close(c.sendMsg)
		close(c.sendChatInfo)
		c.online = false
	}
}



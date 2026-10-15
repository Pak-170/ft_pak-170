package web

import (
	"log"
	"go_pakito/models"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {} //!Gestionar el tema de conexión por proxy

const (
	join = 0
	leave = 1
	chatMsg = 2
)

type user struct {
	ws		*websocket.Conn
	db		*models.User
}

type chat struct {
	activeUsers []*user
}


func (c *chat)addUser(user *user) {
	c.activeUsers = append(c.activeUsers, user)
	err := user.ws.WriteMessage(1, []byte("Welcome " + user.db.Username + "!"))
	if err != nil {
		log.Println(err)
		return
	}
	c.broadcast("", user, join)
}


func (c *chat)reader(user *user) {

	ws := user.ws

	c.addUser(user)
	for {
		_, msg, err := ws.ReadMessage()

		if err != nil {
			log.Println(err);
			break
		}
		c.broadcast(string(msg), user, chatMsg)
	}
}





func (c *chat)broadcast(msg string, sender *user, msgType int) {

	// if msg == "quit" {
	// 	chat.removeUser(sender)
	// 	sender.ws.Close()
	// }

	var finalMsg []byte

	if msgType == join {
		finalMsg = []byte("user " + sender.db.Username + " joined the chat")
	} else if msgType == leave {
		finalMsg = []byte("user " + sender.db.Username + " left the chat")
	} else {
		finalMsg = []byte(sender.db.Username + ": " + msg)
	}

	for _,user := range c.activeUsers {
		if user == sender {
			continue
		}
		err := user.ws.WriteMessage(1, finalMsg)
		if err != nil {
			log.Println(err)
		}
	}
}
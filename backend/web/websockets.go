package web

import (
	"log"
	"gorm.io/gorm"
	"go_pakito/models"
	"github.com/gorilla/websocket"
	"net/http"
	"go_pakito/dbp"
)

var upgrader = websocket.Upgrader {CheckOrigin: func(r *http.Request) bool { return true }} //!Gestionar el tema de conexión por proxy

const (
	join = 0
	leave = 1
	chatMsg = 2
)

type user struct {
	ws		*websocket.Conn
	info	*models.User
	dbconn	*gorm.DB //! crear un canal para gestionar concurrencia
}

// type chat struct {
// 	chatInfo	*models.Chat
// 	activeUsers []*user // ? maybe change to a map to check if an user is already connected to the chat??
// 	msgs		[]*models.Message
// }


// func (c *chat)addUser(user *user) {
// 	c.activeUsers = append(c.activeUsers, user)
// 	err := user.ws.WriteMessage(1, []byte("Welcome " + user.info.Username + "!"))
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	c.broadcast("", user, join)
// }


func (u *user)reader() {

	ws := u.ws

	allChats, err := dbp.SelectAllChatsOfUser(u.dbconn, u.info.ID)
	if err != nil {
		log.Printf("DB ERROR: %v\n", err)
		return
	}
	log.Print(allChats)
	u.ws.WriteJSON(allChats)
	//c.addUser(u)
	for {
		_, _, err := ws.ReadMessage() //! HABRÁ QUE PARSEAR DESPUÉS PARA RECIBIR LOS MENSAJES COMO JSON
		if err != nil {
			log.Println(err);
			break
		}
		//! Falta meter los mensajes en la BD
		//c.broadcast(string(msg), u, chatMsg)
	}
}


// func (c *chat)broadcast(msg string, sender *user, msgType int) {

// 	// if msg == "quit" {
// 	// 	chat.removeUser(sender)
// 	// 	sender.ws.Close()
// 	// }

// 	var finalMsg []byte

// 	if msgType == join {
// 		finalMsg = []byte("user " + sender.db.Username + " joined the chat")
// 	} else if msgType == leave {
// 		finalMsg = []byte("user " + sender.db.Username + " left the chat")
// 	} else {
// 		finalMsg = []byte(sender.db.Username + ": " + msg)
// 	}

// 	for _,user := range c.activeUsers {
// 		if user == sender {
// 			continue
// 		}
// 		err := user.ws.WriteMessage(1, finalMsg)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}
// }
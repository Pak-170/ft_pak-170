package web

import (
	//"encoding/json"
	"go_pakito/dbp"
	"go_pakito/models"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {CheckOrigin: func(r *http.Request) bool { return true }} //!Gestionar el tema de conexión por proxy

const (
	join = 0
	leave = 1
	chatMsg = 2
)

type User struct {
	ws		*websocket.Conn
	info	*models.User
}



func (h *Hub)addUser(user *User) {
	chatIds, err := dbp.SelectChatIds(h.DB, int(user.info.ID))

	for _, chatID := range chatIds {
		h.chats[chatID] = append(h.chats[chatID], user)
	}
	if err != nil {
		log.Println(err)
		return
	}
}


//*{"ChatID":1,"Username":"yuri","Msg":"sadasdasd","SendTime":"2026-10-19T08:47:38.363Z"}
func	(h *Hub)reader(u *User) {

	ws := u.ws

	allChats, err := dbp.SelectAllChatsOfUser(h.DB, u.info.ID)
	if err != nil {
		log.Printf("DB ERROR: %v\n", err)
		return
	}
	log.Print(allChats)
	u.ws.WriteJSON(allChats)
	h.addUser(u)
	var	msg models.ChatView
	for {
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("ERROR FORMATING WS MSG TO JSON: %v\n", err);
			continue
		}
		if dbp.InsertMsg(h.DB, &msg) != nil {
			u.ws.WriteMessage(1, []byte("ERROR SAVING MESSAGE INSIDE DB"))
			continue
		}
		h.broadcast(&msg, u)
	}
}


func (h *Hub)broadcast(msg *models.ChatView, sender *User) {
	for _,user := range h.chats[int(msg.ChatID)] {
		if user == sender {
			continue
		}
		err := user.ws.WriteJSON(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
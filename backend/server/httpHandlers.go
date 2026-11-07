package server

import (
	"context"
	"go_pakito/dbp"
	"go_pakito/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// - Chat Handlers ------------------------------------------------------------

func (h *Hub) sendUsers(c *gin.Context) {
	users, err := dbp.SelectAny[models.User](c.Request.Context(), h.DB)
	if err != nil {
		log.Printf("ERROR GETTING USERS FROM DB: %v\n", err)
		return
	}
	log.Println(users)
	c.JSON(http.StatusOK, users)
}

func (h *Hub) sendChats(c *gin.Context) {
	userID, exits := c.Get("userID")
	if exits != true {
		return
	}

	chats, _, err := dbp.SelectAllChatsOfUser(c.Request.Context(), h.DB, userID.(uint))
	if err != nil {
		log.Printf("ERROR GETTING CHATS FROM DB: %v\n", err)
		return
	}
	// if chats == nil {
	// 	c.JSON(http.StatusOK, make([]models.Chat, 0))
	// } else {
		c.JSON(http.StatusOK, chats)
	//}
	log.Printf("---> %v <---\n", chats)
}


func	(h *Hub) createChat(c *gin.Context) {
	var chat CreateChatData
	var err error = c.BindJSON(&chat)
	if err != nil {
		log.Println("ERROR: Coun't Bind chat creation JSON")
		return
	}

	ownerID, exists := c.Get("userID")
	if exists == false {
		log.Println("ERROR: User doesn't exist")
		return
	}

	chatID, err := dbp.InsertChat(c.Request.Context(), h.DB, chat.Name, ownerID.(uint))
	if err != nil {
		log.Printf("ERROR: Couldn't insert into DB -> %v\n", err)
		return
	}
	h.chats[chatID] = make([]*Client, 0)

	err = dbp.InsertUserInChat(c.Request.Context(), h.DB, chatID, ownerID.(uint))
	if err != nil {
		log.Printf("ERROR: Couldn't insert into DB -> %v\n", err)
		return
	}
	for _, user := range chat.UserIDs {
		err = dbp.InsertUserInChat(c.Request.Context(), h.DB, chatID, uint(user))
		if err != nil {
			log.Printf("ERROR: Couldn't insert into DB -> %v\n", err)
			return
		}
	}

	participants, err := dbp.SelectUsersOfChat(c.Request.Context(), h.DB, chatID)
	if err != nil {
		log.Printf("ERROR: DB -> %v\n", err)
		return
	}

	var newChatNotification newChatNotification = newChatNotification{chat: &models.Chat{
		ID: chatID, 
		Name: chat.Name, 
		Participants: participants,
	}}

	for _, user := range chat.UserIDs {
		log.Printf("**** USER: %v\n", user)
		newChatNotification.notify(h.clients[user])
	}
	c.JSON(http.StatusOK, chat)
}

// - WS Handler----------------------------------------------------------------

func (h *Hub) sendPreviousMsgs(chatsInfo *[]models.Chat, priorMsgs *[]models.ChatView, client *Client) {

	log.Println("------------------------------------------------------")
	for _, chatInfo := range *chatsInfo {
		client.sendChatInfo <- &chatInfo
		log.Println(chatInfo)
	}
	log.Println("------------------------------------------------------")
	
	for _, msg := range *priorMsgs {
		client.sendMsg <- &msg
		log.Println(msg)
	}
	log.Println("------------------------------------------------------")

}

func (h *Hub) handleWS(c *gin.Context) {
	ctx := context.Background()

	// Upgrading connection.
	wsconn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Printf("ERROR UPGRADING HTTP CONECTION: %v\n", err)
		return
	}

	// First, we authenticate the user by their session cookie.
	userID, exists := c.Get("userID")
	if exists == false {
		log.Printf("ERROR: Couldn't get userID\n")
		return
	}
	client := h.clients[userID.(uint)]
	client.ws = wsconn

	chatInfo, priorMsgs, err := dbp.SelectAllChatsOfUser(ctx, h.DB, client.info.ID)

	if err != nil {
		log.Printf("ERROR RETRIEVING OLS MSGS FOR USER %s\n", userID)
		return
	}
	h.join <- client
	go client.readFromWS(ctx)
	go client.writeToWS(ctx)
	h.sendPreviousMsgs(&chatInfo, &priorMsgs, client)
}


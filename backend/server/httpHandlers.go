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
	chats, err := dbp.SelectAny[models.Chat](c.Request.Context(), h.DB)
	if err != nil {
		log.Printf("ERROR GETTING CHATS FROM DB: %v\n", err)
		return
	}
	log.Println(chats)
	c.JSON(http.StatusOK, chats)
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
	//h.chats[chatID] = append(h.chats[chatID], h.clients[ownerID.(uint)])

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

	//? Is this needed now that we use cookies?
	// Then, we check if the user exists
	var queryUser models.User
	err = h.DB.First(&queryUser, userID).Error
	if err != nil {
		log.Printf("ERROR. USER %s DOESN'T EXISTS\n", userID)
		return
	}

	var client Client = newClient(h, wsconn, &queryUser)

	chatInfo, priorMsgs, err := dbp.SelectAllChatsOfUser(ctx, h.DB, client.info.ID)

	if err != nil {
		log.Printf("ERROR RETRIEVING OLS MSGS FOR USER %s\n", userID)
		return
	}
	h.join <- &client
	go client.readFromWS(ctx)
	go client.writeToWS(ctx)
	h.sendPreviousMsgs(&chatInfo, &priorMsgs, &client)
}


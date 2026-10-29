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

func (h *Hub) getUsers(c *gin.Context) {
	users, err := dbp.SelectAny[models.User](c.Request.Context(), h.DB)
	if err != nil {
		log.Printf("ERROR GETTING USERS FROM DB: %v\n", err)
		return
	}
	log.Println(users)
	c.JSON(http.StatusOK, users)
}

func (h *Hub) getChats(c *gin.Context) {
	userID, exits := c.Get("userID")
	if exits != true {
		return
	}

	chats, _, err := dbp.SelectAllChatsOfUser(c.Request.Context(), h.DB, userID.(uint))
	if err != nil {
		log.Printf("ERROR GETTING CHATS FROM DB: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, chats)
}

func	(h *Hub)postChats(c *gin.Context) {

	var chat CreateChatData

	// We obtain info from the http petition
	var err error = c.BindJSON(&chat)
	if err != nil {
		log.Println("ERROR: Coun't Bind chat creation JSON")
		return
	}

	// We obtain the ID from the petition sender, the chat owner
	ownerID, exists := c.Get("userID")
	if exists == false {
		log.Println("ERROR: User doesn't exist")
		return
	}

	// Chat creation and insertion into DB and HUB
	chatID, err := dbp.InsertChat(c.Request.Context(), h.DB, chat.Name, ownerID.(uint))
	if err != nil {
		log.Printf("ERROR: Couldn't insert into DB -> %v\n", err)
		return
	}
	h.chats[chatID] = models.NewChat(chatID, chat.Name, ownerID.(uint), nil)

	// Insertion of the users
	err = h.insertUsersInChat(c.Request.Context(), chatID, append(chat.UserIDs, ownerID.(uint)))
	if err != nil {
		log.Printf("ERROR: couldn't insert user into DB -> %v\n", err)
	}

	var newChatNotification newChatNotification = newChatNotification{h.chats[chatID]}

	for _, user := range chat.UserIDs {
		newChatNotification.notify(h.clients[user])
	}
	c.JSON(http.StatusOK, chat)
}

// - WS Handler----------------------------------------------------------------

func (h *Hub) sendPreviousMsgs(chatsInfo *[]models.Chat, priorMsgs *[]models.ChatView, client *Client) {

	for _, chatInfo := range *chatsInfo {
		client.sendChatInfo <- &chatInfo
	}

	for _, msg := range *priorMsgs {
		client.sendMsg <- &msg
	}
}

func cancelPriorContext(client *Client) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	if client.ws != nil {
		log.Printf("SE HA RECARGADO LA PÁGINAAAAAAAAAAAA\n\n");
		client.ws.Close()
		client.cancel()
	}
	client.cancel = cancel
	return ctx
}

func (h *Hub) handleWS(c *gin.Context) {

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
		wsconn.Close()
		return
	}

	log.Printf("USUARIO IDENTIFICADOOOOOO: %v\n", userID)
	client := h.clients[userID.(uint)]

	ctx := cancelPriorContext(client)
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


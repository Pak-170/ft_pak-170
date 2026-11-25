package server

import (
	"context"
	"go_pakito/Notification"
	"go_pakito/Session"
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Hub) getChats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if exists != true {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}

	chats, _, err := dbp.SelectAllChatsOfUser(c.Request.Context(), h.DB, userID.(uint))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
		return
	}
	c.JSON(http.StatusOK, chats)
}

func (h *Hub) postChats(c *gin.Context) {

	var chat models.CreateChatData

	// We obtain info from the http petition
	var err error = c.BindJSON(&chat)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "JSON no válido"})
		log.Println("ERROR: Coun't Bind chat creation JSON")
		return
	}

	// We obtain the ID from the petition sender, the chat owner
	ownerID, exists := c.Get("userID")
	if exists == false {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error interno de servidor"})
		log.Println("ERROR: User doesn't exist")
		return
	}

	// Chat creation and insertion into DB and HUB
	chatID, err := dbp.InsertChat(c.Request.Context(), h.DB, chat.Name, ownerID.(uint))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error interno de servidor"})
		log.Printf("ERROR: Couldn't insert into DB -> %v\n", err)
		return
	}

	// Insertion of the users
	err = h.insertUsersInChat(c.Request.Context(), chatID, append(chat.UserIDs, ownerID.(uint)))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error interno de servidor"})
		log.Printf("ERROR: couldn't insert user into DB -> %v\n", err)
		return
	}

	chatDB, err := dbp.SelectAnyWhere[models.Chat](c.Request.Context(), h.DB, "id = ?", chatID)
	if err != nil || len(chatDB) != 1 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error interno de servidor"})
		log.Printf("error finding chat in DB -> %v\n", err)
		return
	}

	var newChat Notification.NewChat = Notification.NewChat{Chat: &chatDB[0]}

	sessions, err := dbp.SelectSessionsFromChatID(c.Request.Context(), h.DB, chatID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error interno de servidor"})
		log.Printf("ERROR: Database Petition Failed -> %v\n", err)
		return
	}
	sessionToken, err := getHashToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error interno de servidor"})
		log.Printf("ERROR: There was some problem retrieving needed information from the petition...\n")
		return
	}
	for _, sessionModel := range sessions {
		sessionSession, exists := h.sessions[sessionModel.ID] // Qué clase de nombres son estos Yolanda por favor
		if exists == false || sessionSession.GetHashedToken() == sessionToken {
			continue
		}
		newChat.Notify(&sessionSession)
	}

	c.JSON(http.StatusCreated, chatDB[0])
}

// - WS Handler----------------------------------------------------------------

func (h *Hub) sendPreviousMsgs(chatsInfo *[]models.Chat, priorMsgs *[]models.ChatView, s *session.Session) {

	for _, chatInfo := range *chatsInfo {
		s.SendChatInfo <- &chatInfo
	}

	for _, msg := range *priorMsgs {
		s.SendMsg <- &msg
	}
}

func (h *Hub) handleWS(c *gin.Context) {

	// Upgrading connection.
	wsconn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ERROR UPGRADING HTTP CONECTION: %v\n", err)
		return
	}
	// First, we authenticate the user by their session cookie.
	sessionToken, err := getHashToken(c)
	if err != nil {
		log.Printf("ERROR: There was some problem retrieving needed information from the petition...\n")
		wsconn.Close()
		return
	}

	user, err := dbp.SelectUserFromToken(c.Request.Context(), h.DB, sessionToken)
	if err != nil {
		log.Printf("ERROR: There was some problem retrieving needed information from the petition...\n")
		wsconn.Close()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	s := session.NewSession(sessionToken, wsconn, &user, &cancel)

	chatInfo, priorMsgs, err := dbp.SelectAllChatsOfUser(ctx, h.DB, s.Info.ID)
	if err != nil {
		log.Printf("ERROR RETRIEVING OLS MSGS FOR USER %v\n", s.Info.ID)
		return
	}
	if err := h.Join(ctx, &s); err != nil {
		log.Println(err)
		return
	}
	go s.ReadFromWS(ctx, h, h.DB)
	go s.WriteToWS(ctx)

	h.sendPreviousMsgs(&chatInfo, &priorMsgs, &s)
}

package web

import (
	"go_pakito/models"
//	"go_pakito/web"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /user/:id
func (h *httpHandler) getUserById(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := h.gormdb.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "User not found"})
	} else {
		c.JSON(http.StatusOK, user)
	}
}

// GET /chat/:id/users
func (h *httpHandler) getUsersInChat(c *gin.Context) {
	chatId := c.Param("id")
	var usersInChat []models.UsersInChatView

	if err := h.gormdb.DB.Find(&usersInChat, chatId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Chat not found"})
	} else {
		c.JSON(http.StatusOK, usersInChat)
	}
}

func (h *httpHandler) home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"title": "PRUEBA OAUTH"})
}


// - WS Handler----------------------------------------------------------------

func (h *httpHandler) handleWS(c *gin.Context) {

	wsconn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Printf("ERROR UPGRADING HTTP CONECTION: %v\n", err)
		return
	}
	defer wsconn.Close()

	var chatID = c.Param("chatID")
	var modelChat []models.ChatView

	err = h.gormdb.DB.Where("chat_id = ?", chatID).Find(&modelChat).Error

	if err != nil {
		log.Printf("ERROR. NO CHAT WAS FOUND WITH CHAT ID %s\n", chatID)
		return
	}
	wsconn.WriteJSON(modelChat)

	var username string = c.Query("username")
	if username == "" {
		log.Printf("ERROR. USERNAME NEEDED AS QUERY PARAMETER TO CONNECT TO CHAT\n")
		return
	}

	var queryUser models.User
	err = h.gormdb.DB.Where("username = ?", username).First(&queryUser).Error
	if err != nil {
		log.Printf("ERROR. NO USER WAS FOUND WITH USER ID %s\n", username)
		return
	}

	var userws user = user{ws: wsconn, db: &queryUser}

	var newchat chat

	newchat.reader(&userws)
}

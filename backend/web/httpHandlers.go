package web

import (
	"go_pakito/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /user/:id
func (db *GormDB) getUserById(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "User not found"})
	} else {
		c.JSON(http.StatusOK, user)
	}
}

// GET /chat/:id/users
func (db *GormDB) getUsersInChat(c *gin.Context) {
	chatId := c.Param("id")
	var usersInChat []models.UsersInChatView

	if err := db.DB.Find(&usersInChat, chatId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Chat not found"})
	} else {
		c.JSON(http.StatusOK, usersInChat)
	}
}

func (db *GormDB) home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"title": "PRUEBA OAUTH"})
}


// - WS Handler----------------------------------------------------------------

func (db *GormDB) handleWS(c *gin.Context) {

	// Upgrading connection.
	wsconn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Printf("ERROR UPGRADING HTTP CONECTION: %v\n", err)
		return
	}
	defer wsconn.Close()

	// First, we check if we received an userID.
	var userID string = c.Query("user-id")
	if userID == "" {
		log.Printf("ERROR. USERID NEEDED AS QUERY PARAMETER TO CONNECT TO CHAT\n")
		return
	}

	// Then, we check if the user exists
	var queryUser models.User
	err = db.DB.First(&queryUser, userID).Error
	if err != nil {
		log.Printf("ERROR. USER %s DOESN'T EXISTS\n", userID)
		return
	}

	var userws user = user{ws: wsconn, info: &queryUser, dbconn: db.DB}
	userws.reader()
}

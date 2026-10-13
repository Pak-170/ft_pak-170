package web

import (
	"go_pakito/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /user/:id
func (h *http_handler) getUserById(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := h.gormdb.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "User not found"})
	} else {
		c.JSON(http.StatusOK, user)
	}
}

// GET /chat/:id/users
func (h *http_handler) getUsersInChat(c *gin.Context) {
	chatId := c.Param("id")
	var usersInChat []models.UsersInChatView

	if err := h.gormdb.DB.Find(&usersInChat, chatId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Chat not found"})
	} else {
		c.JSON(http.StatusOK, usersInChat)
	}
}

func (h *http_handler) home(c *gin.Context) {
	c.JSON(http.StatusOK,gin.H{"title": "PRUEBA OAUTH"})
}

package web

import (
	"go_pakito/dbp"
	"go_pakito/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//GET /user/id/:id
func (h *http_handler)getUserById(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := h.gormdb.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error":"User not found"})
	} else {
		c.JSON(http.StatusOK, user)
	}
}

//GET /user/username/:username INCOMPLETA
func (h *http_handler)getUserByUsername(c *gin.Context) {
	c.Params.ByName("username")

}

//GET /chat/:id/users
func (h *http_handler)getUsersInChat(c *gin.Context) {
	chatId := c.Param("id")
	var usersInChat []models.UsersInChatView

	if err := h.gormdb.DB.Find(&usersInChat, chatId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error":"Chat not found"})
	} else {
		for _, user := range usersInChat {
			c.JSON(http.StatusOK, user) //! Falta formatear bien el JSON
		}
	}
}


func (h *http_handler) chatHandler(w http.ResponseWriter, r *http.Request) {
	var chatName string = r.URL.Path[len("/chat/"):]
	chat, err := dbp.SelectChatByChatname(h.gormdb.DB, chatName)

	if err != nil {
		fmt.Fprintf(w, "Error: chat with chatname %s was not found", chatName)
	} else {
		fmt.Fprintf(w, "%s\n\n", chatName)
		for _,msg := range chat {
			printChat(&w, &msg)
		}
	}
}
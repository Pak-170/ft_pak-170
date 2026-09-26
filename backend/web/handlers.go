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

//GET /user/username/:username
func (h *http_handler)getUserByUsername(c *gin.Context) {
	c.Params.ByName("username")

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
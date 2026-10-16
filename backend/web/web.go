package web

import (
	"fmt"
	"gorm.io/gorm"
	"go_pakito/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GormDB struct {
	DB *gorm.DB
}
// ------------------------ printers

func printUser(w *http.ResponseWriter, user *models.User) {

	fmt.Fprintf(*w, "ID: %d\nUSERNAME: %s\nALIAS:%s\n",
		user.ID, user.Username, *(user.Alias))

}

func printChat(w *http.ResponseWriter, chat *models.ChatView) {
	fmt.Fprintf(*w, "%s: %s (%v)\n", chat.Username, chat.Msg, chat.SendTime)
}

// ------------------------

func Web(db *GormDB) {

	r := gin.Default()
	r.StaticFile("/favicon.ico", "../frontend/public/favicon.png") //!ruta relativa a la de ejecución ahora que no se usa docker aún, beware en el futuroooo

	r.GET("/", db.home)
	r.GET("/user/:id", db.getUserById)
	r.GET("/chat/:id/users", db.getUsersInChat)
	r.GET("/ws/api/chat", db.handleWS)

	r.Run(":8080")
}

package web

import (
	"fmt"
	"go_pakito/dbp"
	"go_pakito/models"
//	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type http_handler struct {
	gormdb *dbp.GormDB
}

// ------------------------ printers

func printUser(w *http.ResponseWriter, user *models.User) {

	alias := ""
	if user.Alias != nil {
		alias = *user.Alias
	}
	fmt.Fprintf(*w, "ID: %d\nUSERNAME: %s\nALIAS:%s\n",
		user.ID, user.Username, alias)
		
}

func printChat(w *http.ResponseWriter, chat *models.ChatView) {
	fmt.Fprintf(*w, "%s: %s (%v)\n", chat.Username, chat.Msg, chat.SendTime)
}


// ------------------------

func Web(db *dbp.GormDB) {
	var h http_handler = http_handler{db}

	r := gin.Default()
	r.StaticFile("/favicon.ico", "../frontend/public/favicon.png") //!ruta relativa a la de ejecución ahora que no se usa docker aún, beware en el futuroooo

	r.GET("/user/id/:id", h.getUserById)
	r.GET("/chat/:id/users", h.getUsersInChat)
	r.GET("/user/username/:username", h.getUserByUsername)
	//http.HandleFunc("/user/", h.getUser)
	//http.HandleFunc("/chat/", h.chatHandler)

	//log.Fatal(http.ListenAndServe(":8080", nil))
	r.Run(":8080")
}

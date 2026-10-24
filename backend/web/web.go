package web

import (
	"go_pakito/dbp"
	"log"
	"gorm.io/gorm"
	"go_pakito/models"
//	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)


type	Hub struct {
	DB		*gorm.DB //! crear un canal para gestionar concurrencia ???
	chats	map[int][]*User
}


// ------------------------

func	initChats(hub *Hub) error {
	var chats	map[int][]*User = make(map[int][]*User, 0)
	chatIds, err := dbp.SelectChatIds(hub.DB, dbp.ALL)

	if err != nil {
		return err
	}
	for _, chatID := range chatIds {
		chats[chatID] = make([]*User, 0)
	}
	hub.chats = chats
	return err
}

func Web(db *gorm.DB) {
	var hub	Hub = Hub{DB: db}

	if err := initChats(&hub); err != nil {
		log.Printf("ERROR INITIALIZING CHATS IDS: %v\n", err)
		return
	}

	r := gin.Default()
	r.StaticFile("/favicon.ico", "../frontend/public/favicon.png")

	// r.GET("/", db.home)
	// r.GET("/user/:id", db.getUserById)
	// r.GET("/chat/:id/users", db.getUsersInChat)
	r.GET("/api/chat", hub.handleWS)

	// --- autentication (publics) ---
	r.POST("/api/register", Register(db))
	r.POST("/api/login", Login(db))
	r.POST("/api/logout", Logout(db))

	protected := r.Group("/api", Auth(db))
	{
		protected.GET("/me", func(c *gin.Context) {
			var u models.User
			if err := db.First(&u, UserID(c)).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
				return
			}
			c.JSON(http.StatusOK, u)
		})
		// Chat endpoints go in here
	}

	//log.Fatal(http.ListenAndServe(":8080", nil))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // para ejecutar fuera de Docker
	}
	r.Run(":" + port)
}

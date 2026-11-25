package server

import (
	"go_pakito/models"
	"net/http"
	"os"

	"gorm.io/gorm"

	"context"

	"github.com/gin-gonic/gin"
)




func NewRouter(ctx context.Context, db *gorm.DB) *gin.Engine{
	var hub Hub = newHub(db)

	go hub.run(ctx)

	r := gin.Default()
	r.StaticFile("/favicon.ico", "../frontend/public/favicon.png")


	// --- autentication (publics) ---
	r.POST("/api/register", hub.Register(db))
	r.POST("/api/login", Login(db))
	r.POST("/api/logout", Logout(db))

	protected := r.Group("/api", Auth(db))
	{
		protected.GET("/me", func(c *gin.Context) {
			var u models.User
			userID := UserID(c)
			if err := db.First(&u, userID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
				return
			}
			// if hub.clients[userID] != nil {
			// 	hub.clients[userID].online = true
			// }
			c.JSON(http.StatusOK, u)
		})
		protected.POST("/chats", hub.postChats)
		protected.GET("/chats", hub.getChats)
		protected.GET("/users", hub.getUsers)
		protected.GET("/chat", hub.handleWS)
	}
	return r
}

func InitServer(ctx context.Context, db *gorm.DB) {
	r := NewRouter(ctx, db)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // para ejecutar fuera de Docker
	}
	r.Run(":" + port)
}

package server

import (
	"gorm.io/gorm"
	"go_pakito/models"
//	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"context"
)

func InitServer(ctx context.Context, db *gorm.DB) {

	var hub Hub = newHub(db)

	go hub.run(ctx)

	r := gin.Default()
	r.StaticFile("/favicon.ico", "../frontend/public/favicon.png")

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

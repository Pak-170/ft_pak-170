package web

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go_pakito/models"
)

const (
	cookieName = "session"
	sessionTTL = 24 * time.Hour
)

type credentials struct {
	Username string  `json:"username" form:"username"`
	Password string  `json:"password" form:"password"`
	Alias    *string `json:"alias" form:"alias"`
}

func hashToken(t string) string {
	sum := sha256.Sum256([]byte(t))
	return hex.EncodeToString(sum[:])
}

func newSession(db *gorm.DB, userID uint) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	s := models.Session{
		ID:        hashToken(token),
		UserID:    userID,
		ExpiresAt: time.Now().Add(sessionTTL),
	}
	return token, db.Create(&s).Error
}

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cred credentials
		if err := c.ShouldBind(&cred); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			return
		}
		if l := len(cred.Username); l < 3 || l > 20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "el usuario debe tener entre 3 y 20 caracteres"})
			return
		}
		if l := len(cred.Password); l < 8 || l > 72 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "la contraseña debe tener entre 8 y 72 bytes"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(cred.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
			return
		}

		user := models.User{Username: cred.Username, Alias: cred.Alias, PasswordHash: string(hash)}
		if err := db.Create(&user).Error; err != nil {
			var me *mysql.MySQLError
			if errors.As(err, &me) && me.Number == 1062 {
				c.JSON(http.StatusConflict, gin.H{"error": "ese usuario ya existe"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
			return
		}

		if err := startSession(c, db, user.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "usuario creado, pero falló el inicio de sesión"})
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cred credentials
		if err := c.ShouldBind(&cred); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
			return
		}

		var user models.User
		err := db.Where("username = ?", cred.Username).First(&user).Error
		if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(cred.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario o contraseña incorrectos"})
			return
		}

		if err := startSession(c, db, user.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func Logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Cookie(cookieName); err == nil {
			db.Delete(&models.Session{}, "id = ?", hashToken(cookie))
		}
		http.SetCookie(c.Writer, &http.Cookie{Name: cookieName, Value: "", Path: "/", MaxAge: -1, HttpOnly: true})
		c.Status(http.StatusNoContent)
	}
}

func startSession(c *gin.Context, db *gorm.DB, userID uint) error {
	token, err := newSession(db, userID)
	if err != nil {
		return err
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:	cookieName,
		Value:	token,
		Path:	"/",
		MaxAge:	int(sessionTTL.Seconds()),
		HttpOnly:	true,
		SameSite:	http.SameSiteLaxMode,
	})
	return nil
}

// Middleware: hooks with group.Use(Auth(db)).
func Auth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(cookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no has iniciado sesión"})
			return
		}

		var s models.Session
		err = db.Where("id = ? AND expires_at > ?", hashToken(cookie), time.Now()).First(&s).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "sesión inválida o caducada"})
			return
		}

		c.Set("userID", s.UserID)
		c.Next()
	}
}

// Inside a protected handler: uid := UserID(c)
func UserID(c *gin.Context) uint {
	return c.GetUint("userID")
}
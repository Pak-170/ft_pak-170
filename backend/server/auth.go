package server

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

// newToken returns 64 random bits from the cryptographic generator of the system.
// Never use math/rand for this: is predictable.
func newToken() (uint64, error) {
	var b [8]byte
	for {
		if _, err := rand.Read(b[:]); err != nil {
			return 0, err
		}
		if t := binary.BigEndian.Uint64(b[:]); t != 0 { // 0 se reserva como "sin token"
			return t, nil
		}
	}
}

// parseToken conveerts the value of the cookie into uint64.
// If it isn't valid, it doesn't make a request to the DB.
func parseToken(s string) (uint64, bool) {
	t, err := strconv.ParseUint(s, 10, 64)
	return t, err == nil && t != 0
}

func newSession(db *gorm.DB, userID uint) (uint64, error) {
	token, err := newToken()
	if err != nil {
		return 0, err
	}
	s := models.Session{
		ID:        token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(sessionTTL),
	}
	return token, db.Create(&s).Error
}

func (h *Hub) Register(db *gorm.DB) gin.HandlerFunc {
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
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				c.JSON(http.StatusConflict, gin.H{"error": "ese usuario ya existe"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno"})
			return
		}
		h.clients[user.ID] = newClient(h, nil, &user)
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
			if id, ok := parseToken(cookie); ok {
				db.Delete(&models.Session{}, "id = ?", id)
			}
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
		Name:     cookieName,
		Value:    strconv.FormatUint(token, 10),
		Path:     "/",
		MaxAge:   int(sessionTTL.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
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
		id, ok := parseToken(cookie)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "sesión inválida o caducada"})
			return
		}

		var s models.Session
		if err := db.Where("id = ? AND expires_at > ?", id, time.Now()).First(&s).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "sesión inválida o caducada"})
			return
		}

		c.Set("userID", s.UserID)
		c.Set("sessionID", s.ID) // uint64
		c.Next()
	}
}

// Inside a protected handler: uid := UserID(c)
func UserID(c *gin.Context) uint {
	return c.GetUint("userID")
}

// Inside a protected handler: sid := SessionID(c)
func SessionID(c *gin.Context) uint64 {
	return c.GetUint64("sessionID")
}
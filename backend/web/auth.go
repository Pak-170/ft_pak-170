package web

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go_pakito/models"
)

const (
	cookieName = "session"
	sessionTTL = 24 * time.Hour
)

type ctxKey string

const UserIDKey ctxKey = "userID"

type credentials struct {
	Username	string `json:"username"`
	Password	string `json:"password"`
	Alias		*string `json:"alias"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string){
	writeJSON(w, status, map[string]string{"error":msg})
}

func hasToken(t string) string {
	sum := sha256.Sum256([]byte(t))
	return hex.EncodeToString(sum[:])
}

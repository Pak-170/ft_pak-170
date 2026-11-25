package test

import (
	"encoding/json"

	"go_pakito/models"
	"go_pakito/server"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"

	"log"
)

var cookie *http.Cookie

// func TestGetChats(t *testing.T) {
// 	_, r, w := setupTest(t)

// 	req, _ := http.NewRequest("GET", "/api/chats", nil)
// 	req.Header.Set("Content-Type", "application/json")
// 	r.ServeHTTP(w, req)
// 	assert.Equal(t, http.StatusOK, w.Code)
// }

const NEWUSERNAME = "userTest"

func TestNewUser(t *testing.T) {

	_, r, w := setupTest(t)

	jsonUser := `{"username":"userTest", "password":"12345678", "alias":null}`

	t.Run("register", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/register", strings.NewReader(jsonUser))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

	})

	t.Run("login", func(t *testing.T) {
		w = httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/login", strings.NewReader(jsonUser))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		cookie = w.Result().Cookies()[0]
	})

	t.Run("get-chats", func(t *testing.T) {
		w = httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/chats", nil)
		req.AddCookie(cookie)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("get-users", func(t *testing.T) {
		w = httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/users", nil)
		req.AddCookie(cookie)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("post-new-chat", func(t *testing.T) {
		w = httptest.NewRecorder()
		chat := server.CreateChatData {
			InfoType: "chat",
			Name: "test-chat",
			ChatType: "group",
			UserIDs: []uint{1}, //pepe
		}
		jsonchat, _ := json.Marshal(chat)
		req, _ := http.NewRequest("POST", "/api/chats", strings.NewReader(string(jsonchat)))
		req.AddCookie(cookie)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
	server := httptest.NewServer(r)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/chat"

	var newchatid int
	t.Run("ws-new-chat", func(t *testing.T) {
		conn, resp, err := websocket.DefaultDialer.Dial(wsURL,
			http.Header{"Cookie": []string{cookie.String()}})
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		if err != nil {
			log.Println("error")
			return
		}
		defer conn.Close()
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, readmsg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
		} else {
			newchatid = int(readmsg[2])
		}
	})
	t.Run("ws-invalid-chat", func(t *testing.T) {
		conn, resp, err := websocket.DefaultDialer.Dial(wsURL,
			http.Header{"Cookie": []string{cookie.String()}})
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		if err != nil {
			log.Println("error")
			return
		}
		defer conn.Close()
		msg := models.ChatView{
			ChatID:		2,
			Username:	NEWUSERNAME,
			Msg:		"hola",
			SendTime:	time.Now(),
		}
		// Ahora puedes escribir y leer mensajes como si fuera un navegador
		jsonmsg, _ := json.Marshal(msg)
		err = conn.WriteMessage(websocket.TextMessage, jsonmsg)
		if err == nil {
			t.Log("ESTE TEST DEBE FALLAR POR NO ESTAR USER IN CHAT")
		}
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		_, readmsg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
		} else {
			log.Println(string(readmsg))
		}
	})

	t.Run("ws-ok-chat", func(t *testing.T) {
		conn, resp, err := websocket.DefaultDialer.Dial(wsURL,
			http.Header{"Cookie": []string{cookie.String()}})
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		defer conn.Close()
		if err != nil {
			log.Println("error")
		}
		msg := models.ChatView{
			ChatID:		uint(newchatid),
			Username:	NEWUSERNAME,
			Msg:		"hola",
			SendTime:	time.Now(),
		}

		jsonmsg, _ := json.Marshal(msg)
		err = conn.WriteMessage(websocket.TextMessage, jsonmsg)
		if err == nil {
			t.Log("ESTE TEST DEBERÍA FUNCIONAR")
		}
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		_, readmsg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
		} else {
			log.Println(string(readmsg))
		}
	})
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

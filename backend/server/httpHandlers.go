package server

import (
	"context"
	"go_pakito/dbp"
	"go_pakito/models"
	"log"

	"github.com/gin-gonic/gin"
)

// - WS Handler----------------------------------------------------------------

func (h *Hub) sendPreviousMsgs(priorMsgs *[]models.ChatView, client *Client) { //? Necesitamos el contexto aquí??
	for _, msg := range *priorMsgs {
		client.send <- &msg
		log.Println(msg)
	}
}


func (h *Hub) handleWS(c *gin.Context) {
	ctx := context.Background()
	log.Printf("WEBSOCKET CONNECTED!!!!\n")
	// Upgrading connection.
	wsconn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Printf("ERROR UPGRADING HTTP CONECTION: %v\n", err)
		return
	}

	// First, we check if we received an userID.
	var userID string = c.Query("user-id")
	if userID == "" {
		log.Printf("ERROR. USERID NEEDED AS QUERY PARAMETER TO CONNECT TO CHAT\n")
		return
	}

	// Then, we check if the user exists
	var queryUser models.User
	err = h.DB.First(&queryUser, userID).Error
	if err != nil {
		log.Printf("ERROR. USER %s DOESN'T EXISTS\n", userID)
		return
	}

	var client Client = newClient(h, wsconn, &queryUser)

	priorMsgs, err := dbp.SelectAllChatsOfUser(ctx, h.DB, client.info.ID)
	if err != nil {
		log.Printf("ERROR RETRIEVING OLS MSGS FOR USER %s\n", userID)
		return
	}
	h.join <- &client
	go client.readFromWS(ctx)
	go client.writeToWS(ctx)
	h.sendPreviousMsgs(&priorMsgs, &client)
}

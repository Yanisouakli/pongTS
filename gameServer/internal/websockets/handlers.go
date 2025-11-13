package websocket

import (
	"log"
	"net/http"
	"pongServer/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebsocketHandler(hub *Hub, c *gin.Context, gm *handlers.GameManager) {
	id := c.Query("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("failed upgrading connection", err)
	}

	client := &Client{
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: id,
		Gm:     gm,
	}

	hub.Clients[id] = client
	hub.Register <- client

	go client.Read()
	go client.Write()

}

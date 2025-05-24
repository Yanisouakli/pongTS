package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserID string // UserID can be a string or any other type that suits your application
	GameID string // GameID can be a string or any other type that suits your application
}

func (c *Client) Read() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	var isUrlCorrect bool

	for {
		if isUrlCorrect {
			//game data handling here
			//#TODO later:
		}

	}
}

func (c *Client) Write() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("error writing message to client:", err)
		}

	}

}

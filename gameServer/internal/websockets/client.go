package websocket

import (
	"encoding/json"
	"log"
	"pongServer/internal/models"
	"pongServer/internal/services"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub                    *Hub
	Conn                   *websocket.Conn
	Send                   chan []byte
	UserID                 string
	userConnectedn, GameID string
}

func (c *Client) Read() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	var isAuthenticated bool

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("error reading event:", err)
			break
		}

		var eventBody map[string]json.RawMessage
		bodyError := json.Unmarshal(msg, &eventBody)
		if bodyError != nil {
			log.Printf("error Unmarshaling eventBody %v", bodyError)
			continue
		}

		var msgType string
		if err := json.Unmarshal(eventBody["type"], &msgType); err != nil {
			log.Printf("error Unmarshaling msgType %v", err)
		}

		if !isAuthenticated {
			//here i will be checking for the game existence
			if msgType == "init" {
				var initEvent models.WsEvent[models.InitEvent]
				if err := json.Unmarshal(msg, &initEvent); err != nil {
					log.Printf("invalid init json %v", err)
					continue
				}

				errorHandshake := models.WsEvent[models.ErrorEvent]{
					Type: "error",
					Params: models.ErrorEvent{
						Error: "error handshake",
					},
				}

				jsonErrorHandshake, _ := json.Marshal(errorHandshake)
				if err := services.CheckConnectedUser(initEvent.Params.GameID, initEvent.Params.PlayerInit.PlayerID, initEvent.Params.PlayerInit.XPos, initEvent.Params.PlayerInit.YPos); err != nil {
					log.Printf("error while joining the game %v", err)
					c.Hub.Clients[initEvent.Params.PlayerInit.PlayerID].Send <- jsonErrorHandshake
				}
				succesHandshake := models.WsEvent[models.SuccesInitEvent]{
					Type: "init",
					Params: models.SuccesInitEvent{
						Message: "succes handshake",
					},
				}

				jsonSuccesHandshake, handshakeErr := json.Marshal(succesHandshake)

				if handshakeErr != nil {
					log.Printf("error marshaling succes event %v", handshakeErr)
				}
				c.Hub.Clients[initEvent.Params.PlayerInit.PlayerID].Send <- jsonSuccesHandshake

				isAuthenticated = true
			}
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

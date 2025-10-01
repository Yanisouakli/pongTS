package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"pongServer/internal/handlers"
	"pongServer/internal/models"
	"pongServer/internal/services"
	"time"
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserID string
	GameID string
	Gm     *handlers.GameManager
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
				if err := services.CheckConnectedUser(initEvent.Params.GameID, initEvent.Params.PlayerInit.PlayerID, initEvent.Params.PlayerInit.XPos, initEvent.Params.PlayerInit.YPos, c.Gm); err != nil {
					log.Printf("error while joining the game %v", err)
					c.Send <- jsonErrorHandshake
				}
				succesHandshake := models.WsEvent[models.SuccesInitEvent]{
					Type: "succes-handshake",
					Params: models.SuccesInitEvent{
						Message: "succes handshake",
					},
				}

				jsonSuccesHandshake, handshakeErr := json.Marshal(succesHandshake)

				if handshakeErr != nil {
					log.Printf("error marshaling succes event %v", handshakeErr)
					return
				}
				c.Send <- jsonSuccesHandshake

				isAuthenticated = true
			}

		}
		//
		switch msgType {
		case "ack-start-game":
			tick := time.Tick(16 * time.Millisecond)

			go func() {
				for {
					select {
					case <-tick:
						var UpdatesBody models.UpdatesBody
						UpdatesBody = models.UpdatesBody{
							Update: "update body",
						}
						jsonUpdate, err := json.MarshalIndent(UpdatesBody, "", " ")
						if err != nil {
							log.Fatalf("error while marshalling json %v", err)
							return
						}
						c.Send <- jsonUpdate
					}
				}
			}()

		case "input":
			//acccepts input and updates gameState
			var InputEvent models.WsEvent[models.InputEvent]

			err := json.Unmarshal(msg, &InputEvent)
			if err != nil {
				log.Fatalf("error while marshaling json %v", err)
			}
      
      
			//based on the input update the player y position and calculate the ball position

		case "game_over":

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

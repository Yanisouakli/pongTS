package websocket

import (
	"encoding/json"
	"log"
	"pongServer/internal/handlers"
	"pongServer/internal/models"
	"pongServer/internal/services"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserID string
	GameID string
	Gm     *handlers.GameManager
	Ticker *time.Ticker
}

func (c *Client) Read() {
	defer func() {
		if c.Ticker != nil {
			c.Ticker.Stop()
			c.Ticker = nil
		}
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
				c.GameID = initEvent.Params.GameID

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
			if c.Ticker == nil {
				c.Ticker = time.NewTicker(16 * time.Millisecond)

				go func() {
					for range c.Ticker.C {
						game, ok := c.Gm.GetGame(c.GameID)
						if !ok {
							log.Printf("error while getting the game data %v", err)
							continue
						}
						UpdatesBody := models.UpdatesBody{
							Type:    "updates",
							Players: game.Players,
							Ball:    game.State.Ball,
							Score:   game.State.Score,
							Timer:   game.State.Timer,
							Canvas:  game.State.Canvas,
						}
						jsonUpdate, err := json.MarshalIndent(UpdatesBody, "", " ")
						if err != nil {
							log.Printf("error while marshalling json %v", err)
							continue
						}
						c.Send <- jsonUpdate
					}
				}()
			}

		case "input":
			//acccepts input and updates gameState
			var InputEvent models.WsEvent[models.InputEvent]

			err := json.Unmarshal(msg, &InputEvent)
			if err != nil {
				log.Printf("error while marshaling json %v", err)
				continue
			}
			if err := c.Gm.UpdateGame(InputEvent); err != nil {
				log.Printf("error while updating the game data %v", err)
				continue
			}

		case "game_over":
			if c.Ticker != nil {
				c.Ticker.Stop()
				c.Ticker = nil
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

package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"pongServer/internal/handlers"
	"pongServer/internal/models"
	"pongServer/internal/utils"
	"time"
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
				if err := c.Gm.InitGameState(initEvent.Params.GameID, initEvent.Params.PlayerInit, initEvent.Params.BallInit, initEvent.Params.CanvasInit); err != nil {
					log.Printf("error while joining the game %v", err)
					c.Send <- jsonErrorHandshake
				}
				succesHandshake := models.WsEvent[models.SuccesInitEvent]{
					Type: "succes-handshake",
					Params: models.SuccesInitEvent{
						Message:  "succes handshake",
						PlayerID: initEvent.Params.PlayerInit.PlayerID,
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
							log.Printf("game not found for ID: %s", c.GameID)
							continue
						}

						//update racket position
						speed := int64(20)
						for i := range game.Players {
							game.Players[i].PreviousY = game.Players[i].YPos

							if game.Players[i].Direction == "up" && game.Players[i].YPos > 0 {
								if game.Players[i].YPos-speed < 0 {
									game.Players[i].YPos = 0
								} else {
									game.Players[i].YPos -= speed
								}
							} else if game.Players[i].Direction == "down" && game.Players[i].YPos+game.Players[i].Height < game.State.Canvas.CanvasHeight {
								if game.Players[i].YPos+game.Players[i].Height+speed > game.State.Canvas.CanvasHeight {
									game.Players[i].YPos = game.State.Canvas.CanvasHeight - game.Players[i].Height
								} else {
									game.Players[i].YPos += speed
								}
							}
							game.Players[i].VelocityY = (game.Players[i].YPos - game.Players[i].PreviousY) * 10
						}

						game.State.Ball.XPos = game.State.Ball.XPos + game.State.Ball.VelocityX
						game.State.Ball.YPos = game.State.Ball.YPos + game.State.Ball.VelocityY

						if game.State.Ball.YPos <= 0 || game.State.Ball.YPos+game.State.Ball.Height >= game.State.Canvas.CanvasHeight {
							game.State.Ball.VelocityY = game.State.Ball.VelocityY * -1
						}

						for i := range game.Players {
							if utils.CollisionWithRacket(game.State.Ball, game.Players[i]) {
								game.State.Ball.VelocityX = game.State.Ball.VelocityX * -1

								game.State.Ball.VelocityY += game.Players[i].VelocityY / 10

								speedIncrease := float64(1.05)
								game.State.Ball.VelocityX = int64(float64(game.State.Ball.VelocityX) * speedIncrease)
								game.State.Ball.VelocityY = int64(float64(game.State.Ball.VelocityY) * speedIncrease)

								break
							}
						}

						goalResult := utils.GoalHandler(game.State.Ball, game.State.Canvas)
						if goalResult.Goal {
							if goalResult.Player == utils.PlayerMe {
								game.State.Score++
							} else if goalResult.Player == utils.PlayerOpp {
								game.State.Score--
							}

							game.State.Ball.XPos = game.State.Canvas.CanvasWidth / 2
							game.State.Ball.YPos = game.State.Canvas.CanvasHeight / 2

							game.State.Ball.VelocityX = -game.State.Ball.VelocityX
							if game.State.Ball.VelocityX > 0 {
								game.State.Ball.VelocityX = 5
							} else {
								game.State.Ball.VelocityX = -5
							}
							game.State.Ball.VelocityY = 0

							goalEvent := models.WsEvent[models.GoalEvent]{
								Type: "goal",
								Params: models.GoalEvent{
									Player: string(goalResult.Player),
									Score:  game.State.Score,
								},
							}
							jsonGoal, _ := json.Marshal(goalEvent)
							c.Send <- jsonGoal
						}

						c.Gm.SetGame(c.GameID, game)

						UpdatesBody := models.UpdatesBody{
							Type:    "updates",
							Players: game.Players,
							Ball:    game.State.Ball,
							Score:   game.State.Score,
							Timer:   game.State.Timer,
							Canvas:  game.State.Canvas,
						}
						jsonUpdate, err := json.Marshal(UpdatesBody)
						if err != nil {
							log.Printf("error while marshalling json %v", err)
							continue
						}
						c.Send <- jsonUpdate
					}
				}()
			}

		case "input":
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

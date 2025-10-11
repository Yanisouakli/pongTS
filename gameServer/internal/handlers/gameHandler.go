package handlers

import (
	"fmt"
	"net/http"
	"pongServer/internal/models"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GameManager struct {
	games   map[string]models.Game
	gamesMu sync.RWMutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[string]models.Game),
	}
}

func (gm *GameManager) GetGameUrlHandler(c *gin.Context) {
	GameID := uuid.New().String()
	players := make([]models.Player, 2)

	gm.gamesMu.Lock()

	gm.games[GameID] = models.Game{
		GameID:    GameID,
		Players:   players,
		CreatedAt: time.Now(),
		State:     models.GameState{},
	}

	gm.gamesMu.Unlock()
	c.JSON(http.StatusOK, gin.H{"GameID": GameID})
}

func (gm *GameManager) DoesGameExist(gameID string) bool {
	gm.gamesMu.Lock()
	defer gm.gamesMu.Unlock()
	_, ok := gm.games[gameID]
	if !ok {
		return false
	}
	return true
}

func (gm *GameManager) PlayerInGame(gameID string, userID string, x_pos int64, y_pos int64) error {
	gm.gamesMu.Lock()
	defer gm.gamesMu.Unlock()

	game, ok := gm.games[gameID]
	if !ok {
		return fmt.Errorf("game not found")
	}

	for _, player := range game.Players {
		if player.PlayerID == userID {
			return nil
		}
	}

	newPlayer := models.Player{
		PlayerID:  userID,
		XPos:      x_pos,
		YPos:      y_pos,
		Direction: "stop",
		VelocityY: 0,
		Width:     20,
		Height:    100,
	}

	game.Players = append(game.Players, newPlayer)

	gm.games[gameID] = game

	return nil
}

func (gm *GameManager) InitGameState(gameID string, ball models.BallState, players models.Player, canvas models.Canvas) error {
	gm.gamesMu.Lock()
	defer gm.gamesMu.Unlock()

	game, ok := gm.games[gameID]
	if !ok {
		return fmt.Errorf("geme not found")
	}
	game.State.Ball = ball
	game.State.Canvas = canvas

	gm.games[gameID] = game

	return nil
}

func (gm *GameManager) UpdateGame(input models.WsEvent[models.InputEvent]) error {
	gm.gamesMu.Lock()
	defer gm.gamesMu.Unlock()

	game, ok := gm.games[input.Params.GameID]
	if !ok {
		return fmt.Errorf("game not found")
	}

	speed := int64(20)
	canvasH := game.State.Canvas.CanvasHeight

	for i := range game.Players {
		if game.Players[i].PlayerID == input.Params.PlayerID {
			p := &game.Players[i]

			// Direction from key
			switch input.Params.Key {
			case "up", "z":
				p.Direction = "up"
			case "down", "s":
				p.Direction = "down"
			case "stop", "none", "":
				p.Direction = "stop"
			}

			// Move with bounds
			p.PreviousY = p.YPos
			if p.Direction == "up" && p.YPos > 0 {
				if p.YPos-speed < 0 {
					p.YPos = 0
				} else {
					p.YPos -= speed
				}
			} else if p.Direction == "down" && p.YPos+p.Height < canvasH {
				if p.YPos+p.Height+speed > canvasH {
					p.YPos = canvasH - p.Height
				} else {
					p.YPos += speed
				}
			}

			// Velocity
			p.VelocityY = (p.YPos - p.PreviousY) * 10

			gm.games[input.Params.GameID] = game
			return nil
		}
	}

	return fmt.Errorf("player not found")
}

func (gm *GameManager) GetGame(gameID string) (models.Game, bool) {
	gm.gamesMu.RLock()
	defer gm.gamesMu.RUnlock()
	g, ok := gm.games[gameID]
	return g, ok
}

func (gm *GameManager) SetGame(gameID string, g models.Game) {
	gm.gamesMu.Lock()
	defer gm.gamesMu.Unlock()
	gm.games[gameID] = g
}

func (gm *GameManager) GetGamesSnapshot() map[string]models.Game {
	gm.gamesMu.RLock()
	defer gm.gamesMu.RUnlock()
	out := make(map[string]models.Game, len(gm.games))
	for k, v := range gm.games {
		out[k] = v
	}
	return out
}

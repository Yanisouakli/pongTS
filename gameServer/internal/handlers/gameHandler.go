package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"pongServer/internal/models"
	"sync"
	"time"
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
	players := make([]models.Player, 1)

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

func (gm *GameManager) InitGameState(gameID string, player models.Player, ball models.BallState, canvas models.Canvas) error {
	gm.gamesMu.Lock()
	defer gm.gamesMu.Unlock()

	game, ok := gm.games[gameID]
	if !ok {
		return fmt.Errorf("geme not found")
	}

	for _, p := range game.Players {
		if p.PlayerID == player.PlayerID {
			return nil
		}
	}

	newPlayer := models.Player{
		PlayerID:  player.PlayerID,
		XPos:      player.XPos,
		YPos:      player.YPos,
		Direction: "stop",
		VelocityY: 0,
		Width:     20,
		Height:    100,
	}

	game.Players = append(game.Players, newPlayer)
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

	for i := range game.Players {
		if game.Players[i].PlayerID == input.Params.PlayerID {
			switch input.Params.Key {
			case "up", "z":
				game.Players[i].Direction = "up"
			case "down", "s":
				game.Players[i].Direction = "down"
			case "stop", "none", "":
				game.Players[i].Direction = "stop"
			}
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

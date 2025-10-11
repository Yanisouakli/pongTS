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
		PlayerID: userID,
		XPos:     x_pos,
		YPos:     y_pos,
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

	return nil
}

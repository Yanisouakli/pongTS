package handlers

import (
	"fmt"
	"pongServer/internal/models"
	"sync"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


var (
	games   = make(map[string]models.Game)
	gamesMu sync.RWMutex
)

func GetGameUrlHandler(c *gin.Context) {
	GameID := uuid.New().String()
	players := make([]models.Player, 2)

	gamesMu.Lock()

	games[GameID] = models.Game{
		GameID:    GameID,
		Players:   players,
		CreatedAt: time.Now(),
		State:     models.GameState{},
	}
	gamesMu.Unlock()
	c.JSON(http.StatusOK, gin.H{"GameID": GameID})
}

func DoesGameExist(gameID string) bool {
	gamesMu.Lock()
  defer gamesMu.Unlock()
	_, ok := games[gameID]
	if !ok {
		return false
	}
	return true
}

func PlayerInGame(gameID string,userID string,x_pos int64, y_pos int64) error {

  gamesMu.Lock()
	defer gamesMu.Unlock()

	game, ok := games[gameID]
	if !ok {
		return fmt.Errorf("game not found")
	}


  for _,player:= range game.Players {
    if player.PlayerID == userID  {
      return nil
    }
  }

  newPlayer:= models.Player{
    PlayerID:userID, 
    Score : 0,
    XPos  :x_pos,
    YPos  :y_pos,    
  }
   
  game.Players = append(game.Players,newPlayer)

  games[gameID] = game

  return nil
}


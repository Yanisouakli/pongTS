package handlers

import (
	"pongServer/internal/models"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
  "net/http"
)


var(
  games = make(map[string]models.Game)
  gamesMu sync.RWMutex
)

func GetGameUrlHandler(c *gin.Context){
  GameID:= uuid.New().String()
  players := make([]models.Player,2)

  gamesMu.Lock()

  games[GameID] = models.Game{
    GameID: GameID,
    Players:  players, 
    CreatedAt: time.Now(),
    State:models.GameState{},
  }  
  gamesMu.Unlock()
  c.JSON(http.StatusOK,gin.H{"GameID":GameID})  
}

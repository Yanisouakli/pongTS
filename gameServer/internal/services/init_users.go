package services

import (
	"pongServer/internal/handlers"

	"fmt"
)

func CheckConnectedUser(gameID string, userID string, x_pos int64, y_pos int64) error {
	res := handlers.DoesGameExist(gameID)
	if !res {
		return fmt.Errorf("{{this game doesnt exist}}")
	} else {
		if err := handlers.PlayerInGame(gameID, userID, x_pos, y_pos); err != nil {
			return err
		}

	}

	return nil
}

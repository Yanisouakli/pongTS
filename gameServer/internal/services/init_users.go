package services

import (
	"pongServer/internal/handlers"
	"fmt"
)

func CheckConnectedUser(gameID string, userID string, x_pos int64, y_pos int64,h *handlers.GameManager) error {
	res := h.DoesGameExist(gameID)
	if !res {
		return fmt.Errorf("{{this game doesnt exist}}")
	} else {
		if err := h.PlayerInGame(gameID, userID, x_pos, y_pos); err != nil {
			return err
		}

	}

	return nil
}



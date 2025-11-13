package utils

import "pongServer/internal/models"

type Player string

const (
	PlayerMe   Player = "me"
	PlayerOpp  Player = "opp"
	PlayerNone Player = "none"
)

type GoalReturn struct {
	Goal   bool
	Player Player
}

// CollisionWithRacket checks if the ball collides with a player's racket
func CollisionWithRacket(ball models.BallState, racket models.Player) bool {
	return ball.XPos <= racket.XPos+racket.Width &&
		ball.XPos+ball.Width > racket.XPos &&
		ball.YPos <= racket.YPos+racket.Height &&
		ball.YPos+ball.Height > racket.YPos
}

// GoalHandler checks if a goal has been scored
func GoalHandler(ball models.BallState, canvas models.Canvas) GoalReturn {
	if ball.XPos < 0 {
		return GoalReturn{
			Goal:   true,
			Player: PlayerOpp, // Right player scored (ball went off left side)
		}
	} else if ball.XPos >= canvas.CanvasWidth {
		return GoalReturn{
			Goal:   true,
			Player: PlayerMe, // Left player scored (ball went off right side)
		}
	}
	return GoalReturn{
		Goal:   false,
		Player: PlayerNone,
	}
}

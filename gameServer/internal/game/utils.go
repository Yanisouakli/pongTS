package game

import (
	"pongServer/internal/models"
)


func collisionWithRacket(ball models.BallState,racket models.Player) bool {
  return (
    ball.XPos <= racket.XPos + racket.Width &&
    ball.XPos + ball.Width > racket.XPos &&
    ball.YPos <= racket.YPos + racket.Height &&
    ball.YPos + ball.Height > racket.YPos )
}


func goalHandler(ball models.BallState, width int64  ) models.GoalReturn {
  if (ball.XPos < 0) {
    return models.GoalReturn{
      Goal: true,
      Player: "me",
    }
  } else if (ball.XPos >= width) {
    return models.GoalReturn{
      Goal: true,
      Player: "opp",
    }
  }
  return models.GoalReturn{
    Goal: false,
    Player: "none",
  }
}




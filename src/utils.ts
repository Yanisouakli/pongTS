import { Ball } from "./ball";
import { Racket } from "./racket";

type Player = "me" | "opp" | "none"

interface GoalReturnI {
  goal: boolean;
  player: Player;

}

export function collisionWithRacket(ball: Ball, racket: Racket): boolean {
  return (
    ball.x <= racket.x + racket.width &&
    ball.x + ball.width > racket.x &&
    ball.y <= racket.y + racket.height &&
    ball.y + ball.height > racket.y
  );
}

export function goalHandler(ball: Ball, canvas: HTMLCanvasElement): GoalReturnI {
  if (ball.x < 0) {
    return {
      goal: true,
      player: "me"
    }
  } else if (ball.x >= canvas.width) {
    return {
      goal: true,
      player: "opp"
    }
  }
  return {
    goal: false,
    player: "none"
  }

}

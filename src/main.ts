import './style.css'
import { Ball } from './ball'
import { Racket } from './racket'
import { collisionWithRacket, goalHandler } from './utils'
import WsConnection from './connection'


const canvas = document.getElementById("game") as HTMLCanvasElement

const scoreboard = document.getElementById("score")

const MyRacket = new Racket(false, canvas)
const OpRacket = new Racket(true, canvas)
const ball = new Ball(canvas)
const Ws = new WsConnection("0.0.0.0.8080")


const ctx = canvas.getContext("2d")

if (!ctx) throw new Error("canves not supported")

let lastTime = 0
let pausedUntil = 0


function gameLoop(time: number) {
  if (!ctx || !Ws) return;
  const deltaTime = (time - lastTime) / 1000;
  lastTime = time;
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  MyRacket.move(canvas);
  MyRacket.draw(ctx);
  OpRacket.draw(ctx);

  if (time >= pausedUntil) {
    ball.move(canvas, deltaTime);
    ball.draw(ctx);

    if (collisionWithRacket(ball, MyRacket)) {
      ball.velocityX *= -1;
      ball.velocityY += MyRacket.velocityY * 2
    }

    if (collisionWithRacket(ball, OpRacket)) {
      ball.velocityX *= -1;
      ball.velocityY += OpRacket.velocityY * 2
    }

    const { goal, player } = goalHandler(ball, canvas);
    if (goal === true) {
      if (player === "opp") {
        console.log("goal against you");
      } else if (player === "me") {
        console.log("goal for you");
      }
      ball.resetBall(canvas);
      pausedUntil = time + 2000;
    }

  } else {
    ball.draw(ctx);
  }

  requestAnimationFrame(gameLoop);
}

pausedUntil = performance.now() + 2000;

requestAnimationFrame(gameLoop)

import './style.css'
import { Racket } from './racket'


const canvas = document.getElementById("game") as HTMLCanvasElement

const scoreboard = document.getElementById("score")  

const MyRacket = new Racket(false,canvas)
const OpRacket = new Racket(true,canvas)


const ctx=  canvas.getContext("2d")
let lastTime = 0
const speed = 20


if (!ctx) throw new Error("canves not supported")

function gameLoop(time:number){
  const deltaTime = (time - lastTime)/1000
  if (deltaTime > 1 /speed){
    ctx?.clearRect(0, 0, canvas.width, canvas.height);
    lastTime = time;
    MyRacket.move(canvas) 
    MyRacket.draw(ctx!!)
    OpRacket.draw(ctx!!) 
  }

requestAnimationFrame(gameLoop)
  

}

requestAnimationFrame(gameLoop)


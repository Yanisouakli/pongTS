import { startGameLoop } from "./gameLoop"
import WsConnection from "./connection"


const ws = new WsConnection("http://localhost:8080/ws")

export function renderGamePage(root:HTMLElement,gameID:string){
  const canvas = document.createElement("canvas")
  canvas.id= "game" 
  canvas.width= 1200
  canvas.height= 600
  root.appendChild(canvas)
  //connect to the websocket 
  startGameLoop(canvas)
  
}


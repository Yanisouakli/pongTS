import { startGameLoop } from "./gameLoop"
import WsConnection from "./connection"



export function renderGamePage(root: HTMLElement, gameID: string) {
  const playerId = "player-1234";
  const ws = new WsConnection(`ws://localhost:8080/ws?id=${playerId}`);

  ws.onOpen(() => {
    console.log("connectionn established")

    ws.send({
      type: "init",
      params: {
        game_id: gameID,
        player_init: {
          player_id: "string-id",
          score: 0,
          x_pos: 5,
          y_pos: 20,
        }
      },
    });
  })


  ws.onMessage((msg) => {
    console.log("msg", msg)
  })

  const canvas = document.createElement("canvas")
  canvas.id = "game"
  canvas.width = 1200
  canvas.height = 600
  root.appendChild(canvas)
  //connect to the websocket 
  startGameLoop(canvas)
}


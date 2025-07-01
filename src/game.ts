import { startGameLoop } from "./gameLoop"

export function renderGamePage(root: HTMLElement) {

  const canvas = document.createElement("canvas")
  canvas.id = "game"
  canvas.width = 1200
  canvas.height = 600
  root.appendChild(canvas)
  startGameLoop(canvas)
}


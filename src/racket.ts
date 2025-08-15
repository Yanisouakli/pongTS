type Direction = "up" | "down" | "stop";
import WsConnection from "./connection";

export class Racket {
  x: number;
  y: number;
  width: number;
  height: number;
  direction: Direction;
  velocityY: number;
  private keyPressed: Set<string>;
  private previousY: number = 0;

  constructor(isOpp: boolean, canvas: HTMLCanvasElement, ws: WsConnection) {
    this.width = 20;
    this.height = 100;

    this.x = isOpp ? canvas.width - this.width - 20 : 20;
    this.y = canvas.height / 2 - this.height / 2;

    this.direction = "stop";
    this.keyPressed = new Set();
    this.velocityY = 0;
    this.changeDirection(ws);
  }

  draw(ctx: CanvasRenderingContext2D) {
    ctx.fillStyle = "#fff";
    ctx.fillRect(this.x, this.y, this.width, this.height);
  }

  move(canvas: HTMLCanvasElement) {
    const speed = 20;
    this.previousY = this.y;

    if (this.direction === "up" && this.y > 0) {
      this.y -= speed;
    } else if (this.direction === "down" && this.y + this.height < canvas.height) {
      this.y += speed;
    }
    this.velocityY = (this.y - this.previousY) * 10
  }

  changeDirection(ws:WsConnection) {
    window.addEventListener("keydown", (e: KeyboardEvent) => {
      this.keyPressed.add(e.key);
      if (this.keyPressed.has("z") && !this.keyPressed.has("s")) {
        ws.send({
          type:"input",
          params:{
            key:"up"
          }
        })
        this.direction = "up";
      } else if (this.keyPressed.has("s") && !this.keyPressed.has("z")) {
        ws.send({
          type:"input",
          params:{
            key:"down"
          }
        })
        this.direction = "down";
      }
    });

    window.addEventListener("keyup", (e: KeyboardEvent) => {
      this.keyPressed.delete(e.key);
      if (!this.keyPressed.has("z") && !this.keyPressed.has("s")) {
        this.direction = "stop";
      } else if (this.keyPressed.has("z")) {
        this.direction = "up";
      } else if (this.keyPressed.has("s")) {
        this.direction = "down";
      }
    });
  }
}

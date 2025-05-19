type Direction = "up" | "down" | "stop"
export class Racket {
  body: number[];
  direction: Direction;
  private keyPressed: Set<string>;
  constructor(isOpp:boolean,canvas:HTMLCanvasElement) {
    const xOpp = (canvas.width / 20) - 2 
    this.body = isOpp? [xOpp,5] :[1, 5]
    this.direction = "stop"
    this.keyPressed = new Set();
    this.changeDirection()
  }

  draw(ctx: CanvasRenderingContext2D) {
    ctx.fillStyle = "#fff";
    ctx.fillRect(this.body[0] * 20, this.body[1] * 20, 20, 100)
  }
  move(canvas:HTMLCanvasElement) {
    const maxY = (canvas.height - 100) / 20;
    if (this.direction === "up" && this.body[1] > 0) {
      this.body[1] -= 1;
    } else if (this.direction === "down" && this.body[1] < maxY) {
      this.body[1] += 1;
    }

  }

  changeDirection() {
    window.addEventListener("keydown", (e: KeyboardEvent) => {
      this.keyPressed.add(e.key)
      if (this.keyPressed.has("z") && !this.keyPressed.has("s")) {
        this.direction = "up";

      } else if (this.keyPressed.has("s") && !this.keyPressed.has("z")) {
        this.direction = "down"
      }

    })
    window.addEventListener("keyup", (e: KeyboardEvent) => {
      this.keyPressed.delete(e.key)
      if (!this.keyPressed.has("z") && !this.keyPressed.has("s")) {
        this.direction = "stop"
      } else if (this.keyPressed.has("z")) {
        this.direction = "up"
      } else if (this.keyPressed.has("s")) {
        this.direction = "down"
      }
    })
  }
}









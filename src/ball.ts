export class Ball {
  x: number;
  y: number;
  velocityX: number;
  velocityY: number;
  radius: number;
  width: number;
  height: number;

  constructor(canvas: HTMLCanvasElement) {
    this.radius = 10;
    this.width = this.radius * 2;
    this.height = this.radius * 2;

    this.x = canvas.width / 2 - this.radius;
    this.y = canvas.height / 2 - this.radius;
    this.velocityX = -600;
    this.velocityY = 0;
  }

  move(canvas: HTMLCanvasElement, deltaTime: number) {
    this.x += this.velocityX * deltaTime;
    this.y += this.velocityY * deltaTime;

    if (this.y <= 0 || this.y + this.height >= canvas.height) {
      this.velocityY *= -1;
    }
  }

  resetBall(canvas:HTMLCanvasElement){
    this.x = canvas.width / 2 - this.radius;
    this.y = canvas.height / 2 - this.radius;
    this.velocityX = -600;
    this.velocityY = 0;
  }

  draw(ctx: CanvasRenderingContext2D) {
    ctx.beginPath();
    ctx.arc(this.x + this.radius, this.y + this.radius, this.radius, 0, 2 * Math.PI);
    ctx.strokeStyle = "black";
    ctx.stroke();
    ctx.fillStyle = "#fff";
    ctx.fill();
  }
}

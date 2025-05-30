import { renderGamePage } from './game'
import { renderHomePage } from './home'

export function router() {
  const path = window.location.pathname;
  const app = document.getElementById("app")!

  app.innerHTML = ""
  if (path === "/" || path === '/index.html') {
    renderHomePage(app)
  } else if (path.includes("game")) {
    renderGamePage(app,"id")
  } else {
    app.innerHTML = `<h1>404 - Not Found</h1>`

  }
}



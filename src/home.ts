export function renderHomePage(root: HTMLElement) {
  const title = document.createElement("h1")
  title.setAttribute("class", "title")
  title.innerText = "Pong TS"

  const createGameButton = document.createElement("button")

  createGameButton.setAttribute("class", "create_button")
  createGameButton.innerText = "Create Game"

  createGameButton.onclick = async () => {
    const url = await generateGameUrl()
    if (url) {
      const linkToGame = document.createElement("div")
      linkToGame.innerText = `share this link with your mate ${linkToGame}`
      linkToGame.setAttribute("class","game_url")
      root.appendChild(linkToGame)
    } else {
      alert("server failed to generate link")
    }
  }

  root.appendChild(title)
  root.appendChild(createGameButton)
}


async function generateGameUrl(): Promise<string | null> {
  try {
    const res = await fetch("url")
    const { url } = await res.json()
    return url
  } catch (error) {
    console.log("error while createing the game", error)
    return null
  }
}

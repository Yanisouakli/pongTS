export function renderHomePage(root: HTMLElement) {
  const title = document.createElement("h1")
  title.setAttribute("class", "title")
  title.innerText = "Pong TS"

  const createGameButton = document.createElement("button")

  createGameButton.setAttribute("class", "create_button")
  createGameButton.innerText = "Create Game"

  createGameButton.onclick = async () => {
    const url = await generateGameUrl()
    root.innerHTML = ""

    root.appendChild(title)
    root.appendChild(createGameButton)
    
    if (url) {
      const linkToGame = document.createElement("div")
      linkToGame.setAttribute("class","link_container")

      const label = document.createElement("p")
      label.innerText = "this is your link"

      const urlElem = document.createElement("p")
      urlElem.setAttribute("class","link")
      urlElem.innerText = "http://localhost:5173/gameID="+url 

      linkToGame.appendChild(label)
      linkToGame.appendChild(urlElem)

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
    const res = await fetch("http://localhost:8080/generate_game_url")
    const data  = await res.json()
    console.log("GFameID",data.GameID)
    return data.GameID
  } catch (error) {
    console.log("error while createing the game", error)
    return null
  }
}

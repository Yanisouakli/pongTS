package router

import (
	"pongServer/internal/handlers"
	"pongServer/internal/websockets"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
  router := gin.Default()

  hub:= websocket.NewHub()
  go hub.Run()
   
  router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:5173"} ,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE","OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length","Set-Cookie"},
		AllowCredentials: true,
	}))


  router.GET("/generate_game_url", handlers.GetGameUrlHandler)

  router.GET("/ws",func(c *gin.Context){
    websocket.WebsocketHandler(hub,c)
  })

  return router
}

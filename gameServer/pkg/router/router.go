package router

import (
  "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
  "pongServer/internal/handlers"
)



func SetupRouter() *gin.Engine {
  router := gin.Default()
   
  router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:5173"} ,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE","OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length","Set-Cookie"},
		AllowCredentials: true,
	}))


  router.GET("/generate_game_url", handlers.GetGameUrlHandler)

  return router
}



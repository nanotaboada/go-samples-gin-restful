package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/controllers"
)

func GetEngine() *gin.Engine {

	engine := gin.Default()

	engine.GET("/players", controllers.GetPlayers)
	engine.GET("/players/:id", controllers.GetPlayerByID)

	return engine
}

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/controllers"
)

func Setup() *gin.Engine {

	router := gin.Default()

	router.GET("/players/", controllers.GetPlayers)
	router.GET("/players/:id", controllers.GetPlayerByID)
	router.POST("/players/", controllers.CreatePlayer)
	router.PUT("/players/:id", controllers.UpdatePlayer)
	router.DELETE("/players/:id", controllers.DeletePlayer)

	return router
}

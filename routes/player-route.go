/* -----------------------------------------------------------------------------
 * Routes
 * -------------------------------------------------------------------------- */

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/controllers"
)

// Setup configures the router Engine connecting URL paths with controller handlers
func Setup() *gin.Engine {
	router := gin.Default()

	router.GET("/players/", controllers.GetAll)
	router.GET("/players/:id", controllers.GetByID)
	router.POST("/players/", controllers.Post)
	router.PUT("/players/:id", controllers.Put)
	router.DELETE("/players/:id", controllers.Delete)

	return router
}

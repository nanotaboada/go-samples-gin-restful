/* -----------------------------------------------------------------------------
 * Routes
 * -------------------------------------------------------------------------- */

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/controllers"
	"github.com/nanotaboada/go-samples-gin-restful/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup configures the router Engine connecting URL paths with controller handlers
func Setup() *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/players/", controllers.GetAll)
	router.GET("/players/:id", controllers.GetByID)
	router.POST("/players/", controllers.Post)
	router.PUT("/players/:id", controllers.Put)
	router.DELETE("/players/:id", controllers.Delete)

	return router
}

func SetSwaggerInfo() {
	docs.SwaggerInfo.Title = "go-samples-gin-restful"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Description = "🧪 Proof of Concept for a RESTful API made with Go and Gin"
	docs.SwaggerInfo.Schemes = []string{"http"}
}

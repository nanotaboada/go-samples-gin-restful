/* -----------------------------------------------------------------------------
 * Routes
 * -------------------------------------------------------------------------- */

package routes

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/controllers"
	"github.com/nanotaboada/go-samples-gin-restful/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup configures the router Engine connecting URL paths with controller handlers
func Setup() *gin.Engine {
	store := persistence.NewInMemoryStore(time.Hour)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/players/", cache.CachePage(store, time.Hour, controllers.GetAll))
	router.GET("/players/:id", cache.CachePage(store, time.Hour, controllers.GetByID))
	router.GET("/players/squadnumber/:squadnumber", cache.CachePage(store, time.Hour, controllers.GetBySquadNumber))
	router.POST("/players/", ClearCache(store, controllers.Post))
	router.PUT("/players/:id", ClearCache(store, controllers.Put))
	router.DELETE("/players/:id", ClearCache(store, controllers.Delete))

	return router
}

// SetSwaggerInfo configures basic Swagger details
func SetSwaggerInfo() {
	docs.SwaggerInfo.Title = "go-samples-gin-restful"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Description = "🧪 Proof of Concept for a RESTful API made with Go and Gin"
	docs.SwaggerInfo.Schemes = []string{"http"}
}

// ClearCache resets the cache when the collection is modified (POST, PUT, DELETE)
func ClearCache(store *persistence.InMemoryStore, handler gin.HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		keys := []string{
			"players",
			fmt.Sprintf("players/%s", context.Param("id")),
			fmt.Sprintf("players/squadnumber/%s", context.Param("squadnumber")),
		}
		for _, key := range keys {
			store.Delete(key)
		}
		handler(context)
	}
}

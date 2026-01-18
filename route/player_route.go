// Package route sets up the routing and middleware for Player-related endpoints.
package route

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/controller"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup configures the router Engine connecting URL paths with controller handlers
func Setup(controller *controller.PlayerController) *gin.Engine {
	store := persistence.NewInMemoryStore(time.Hour)

	router := gin.Default()

	router.GET(SwaggerPath, ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register routes for /players (without trailing slash)
	router.GET(GetAllPath, cache.CachePage(store, time.Hour, controller.GetAll))
	router.POST(GetAllPath, ClearCache(store, controller.Post))

	// Register alias routes for /players/ (with trailing slash)
	router.GET(GetAllPathTrailingSlash, cache.CachePage(store, time.Hour, controller.GetAll))
	router.POST(GetAllPathTrailingSlash, ClearCache(store, controller.Post))

	router.GET(GetByIDPath, cache.CachePage(store, time.Hour, controller.GetByID))
	router.GET(GetBySquadNumberPath, cache.CachePage(store, time.Hour, controller.GetBySquadNumber))
	router.PUT(GetByIDPath, ClearCache(store, controller.Put))
	router.DELETE(GetByIDPath, ClearCache(store, controller.Delete))

	router.GET(HealthPath, func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}

// ClearCache resets the cache when the collection is modified (POST, PUT, DELETE)
func ClearCache(store persistence.CacheStore, handler gin.HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		id := context.Param(IDParam)

		keys := []string{
			cache.CreateKey(PlayersPath),
			cache.CreateKey(PlayersPathTrailingSlash),
		}
		if id != "" {
			keys = append(keys, cache.CreateKey(fmt.Sprintf("%s/%s", PlayersPath, id)))
		}
		for _, key := range keys {
			_ = store.Delete(key)
		}
		handler(context)
	}
}

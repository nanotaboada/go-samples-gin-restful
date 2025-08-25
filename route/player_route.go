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
func Setup() *gin.Engine {
	store := persistence.NewInMemoryStore(time.Hour)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET(GetAllPath, cache.CachePage(store, time.Hour, controller.GetAll))
	router.GET(GetByIDPath, cache.CachePage(store, time.Hour, controller.GetByID))
	router.GET(GetBySquadNumberPath, cache.CachePage(store, time.Hour, controller.GetBySquadNumber))
	router.POST(GetAllPath, ClearCache(store, controller.Post))
	router.PUT(GetByIDPath, ClearCache(store, controller.Put))
	router.DELETE(GetByIDPath, ClearCache(store, controller.Delete))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}

// ClearCache resets the cache when the collection is modified (POST, PUT, DELETE)
func ClearCache(store *persistence.InMemoryStore, handler gin.HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		keys := []string{
			PlayersKey,
			fmt.Sprintf("%s%s", PlayersKey, "/"+context.Param("id")),
			fmt.Sprintf("%s/squadnumber/%s", PlayersKey, context.Param("squadnumber")),
		}
		for _, key := range keys {
			store.Delete(key)
		}
		handler(context)
	}
}

// Package main initializes and runs the RESTful API server.
//
// It connects to the SQLite3 database, configures routes, and starts the
// Gin HTTP server with Swagger docs enabled.
package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/route"
	"github.com/nanotaboada/go-samples-gin-restful/swagger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	dsn := os.Getenv("STORAGE_PATH")
	// If STORAGE_PATH is not set by Docker Compose,
	if dsn == "" {
		// then the app is running locally in Debug mode.
		dsn = "./storage/players-sqlite3.db"
	}
	data.Connect(dsn)

	store := persistence.NewInMemoryStore(time.Hour)
	router := gin.Default()

	route.RegisterPlayerRoutes(router, store)

	router.GET(route.SwaggerPath, ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET(route.HealthPath, func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	swagger.Setup()
	router.Run(":9000")
}

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
	"github.com/nanotaboada/go-samples-gin-restful/controller"
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/route"
	"github.com/nanotaboada/go-samples-gin-restful/service"
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
	db := data.Connect(dsn)
	playerService := service.NewPlayerService(db)
	playerController := controller.NewPlayerController(playerService)

	store := persistence.NewInMemoryStore(time.Hour)
	app := gin.Default()

	route.RegisterPlayerRoutes(app, playerController, store)

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	swagger.Setup()
	app.Run(":9000")
}

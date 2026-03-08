// Package main initializes and runs the RESTful API server.
//
// It connects to the SQLite3 database, configures routes, and starts the
// Gin HTTP server with Swagger docs enabled.
package main

import (
	"log"
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
	// STORAGE_PATH is injected by Docker Compose (see compose.yaml).
	// When running locally without Docker the variable is empty, so we fall
	// back to the pre-seeded file in the repository.
	dsn := os.Getenv("STORAGE_PATH")
	if dsn == "" {
		// Running locally in debug mode; use the bundled SQLite file.
		dsn = "./storage/players-sqlite3.db"
	}

	// Dependency injection chain: data → service → controller.
	// Each layer depends only on the abstraction of the layer below it:
	//   data.Connect  returns *gorm.DB   (concrete, no interface needed at this level)
	//   NewPlayerService wraps *gorm.DB  and exposes PlayerService (interface)
	//   NewPlayerController wraps PlayerService (interface) — easy to mock in tests
	db := data.Connect(dsn)
	playerService := service.NewPlayerService(db)
	playerController := controller.NewPlayerController(playerService)

	// InMemoryStore is the in-process cache used by gin-contrib/cache.
	// The TTL passed here is the default; individual routes may override it.
	store := persistence.NewInMemoryStore(time.Hour)

	// gin.Default() creates a router pre-configured with two middleware:
	//   Logger  — logs every request (method, path, status, latency) to stdout
	//   Recovery — catches panics, logs the stack trace, and returns 500
	// Use gin.New() if you want a bare router with no middleware.
	app := gin.Default()

	route.RegisterPlayerRoutes(app, playerController, store)

	// The Swagger UI is served at /swagger/index.html.
	// ginSwagger.WrapHandler adapts the swaggerFiles.Handler (an http.Handler)
	// to Gin's handler type.
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Minimal liveness probe — returns {"status":"ok"} with no DB dependency.
	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	swagger.Setup()
	// app.Run blocks until the process exits.  The port is fixed at 9000 to
	// match the Docker EXPOSE directive and the compose.yaml port mapping.
	if err := app.Run(":9000"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// Package main initializes and runs the RESTful API server.
//
// It connects to the SQLite3 database, configures routes, and starts the
// Gin HTTP server with Swagger docs enabled.
package main

import (
	"os"

	"github.com/nanotaboada/go-samples-gin-restful/controller"
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/route"
	"github.com/nanotaboada/go-samples-gin-restful/service"
	"github.com/nanotaboada/go-samples-gin-restful/swagger"
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
	app := route.Setup(playerController)
	swagger.Setup()
	app.Run(":9000")
}

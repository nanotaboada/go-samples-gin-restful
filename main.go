// Package main initializes and runs the RESTful API server.
//
// It connects to the SQLite3 database, configures routes, and starts the
// Gin HTTP server with Swagger docs enabled.
package main

import (
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/route"
	"github.com/nanotaboada/go-samples-gin-restful/swagger"
)

func main() {
	dsn := "./data/players_sqlite3.db"
	data.Connect(dsn)
	app := route.Setup()
	swagger.Setup()
	app.Run(":9000")
}

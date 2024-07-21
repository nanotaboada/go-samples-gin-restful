/* -----------------------------------------------------------------------------
 * Main
 * -------------------------------------------------------------------------- */

package main

import (
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/routes"
)

func main() {
	dsn := "./data/players-sqlite3.db"
	data.Connect(dsn)
	app := routes.Setup()
	routes.SetSwaggerInfo()
	app.Run(":9000")
}

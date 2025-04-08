/* -----------------------------------------------------------------------------
 * Main
 * -------------------------------------------------------------------------- */

package main

import (
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/routes"
	"github.com/nanotaboada/go-samples-gin-restful/swagger"
)

func main() {
	dsn := "./data/players-sqlite3.db"
	data.Connect(dsn)
	app := routes.Setup()
	swagger.Setup()
	app.Run(":9000")
}

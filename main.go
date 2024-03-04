package main

import (
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/routes"
)

func main() {

	path := "./data/players-sqlite3.db"
	data.Connect(path)
	app := routes.Setup()
	app.Run()
}

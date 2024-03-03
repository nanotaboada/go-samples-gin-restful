package main

import (
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/routes"
)

func main() {

	data.ConnectToSqlite()

	app := routes.GetEngine()

	app.Run()
}

package main

import "github.com/nanotaboada/go-samples-gin-restful/routes"

func main() {

	app := routes.GetEngine()

	app.Run()
}

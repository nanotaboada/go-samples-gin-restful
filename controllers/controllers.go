package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/data"
)

var players = data.Seed()

func GetPlayers(context *gin.Context) {

	context.IndentedJSON(http.StatusOK, players)
}

func GetPlayerByID(context *gin.Context) {

	id := context.Param("id")

	for _, player := range players {
		if player.ID == id {
			context.IndentedJSON(http.StatusOK, player)
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not Found."})
}

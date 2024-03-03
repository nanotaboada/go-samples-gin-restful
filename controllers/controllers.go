package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/models"
)

func GetPlayers(context *gin.Context) {

	var players []models.Player
	database := data.Database
	database.Find(&players)
	context.IndentedJSON(http.StatusOK, players)
}

func GetPlayerByID(context *gin.Context) {

	id := context.Param("id")
	database := data.Database
	var player models.Player
	database.Find(&player, id)

	if player.ID == id {
		context.IndentedJSON(http.StatusOK, player)
		return
	}

	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not Found."})
}

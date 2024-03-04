package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/models"
)

func GetPlayers(context *gin.Context) {

	var players []models.Player

	db := data.DB
	db.Find(&players)

	context.IndentedJSON(http.StatusOK, players)
}

func GetPlayerByID(context *gin.Context) {

	var player models.Player
	id := context.Param("id")

	db := data.DB
	db.Find(&player, id)

	if player.ID != "" {
		context.IndentedJSON(http.StatusOK, player)
		return
	}

	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not Found."})
}

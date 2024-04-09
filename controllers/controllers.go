/* -----------------------------------------------------------------------------
 * Controllers
 * -------------------------------------------------------------------------- */

package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/models"
	"gorm.io/gorm"
)

func GetPlayers(context *gin.Context) {
	var players []models.Player
	db := data.DB
	// https://gorm.io/docs/query.html
	db.Find(&players)
	context.IndentedJSON(http.StatusOK, players)
}

func GetPlayerByID(context *gin.Context) {
	id := context.Param("id")
	var player models.Player
	db := data.DB
	// https://gorm.io/docs/query.html
	result := db.First(&player, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			context.Status(http.StatusNotFound)
			return
		}
	}
	context.IndentedJSON(http.StatusOK, player)
}

func CreatePlayer(context *gin.Context) {
	id := context.Param("id")
	var player models.Player
	if err := context.BindJSON(&player); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	db := data.DB
	// https://gorm.io/docs/query.html
	result := db.First(&player, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// https://gorm.io/docs/create.html
			db.Create(&player)
			context.Status(http.StatusCreated)
			return
		}
	}
	context.Status(http.StatusConflict)
}

func UpdatePlayer(context *gin.Context) {
	id := context.Param("id")
	var player models.Player
	var update models.Player
	if err := context.BindJSON(&update); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	db := data.DB
	result := db.First(&player, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			context.Status(http.StatusNotFound)
			return
		}
	}
	// https://gorm.io/docs/update.html
	db.Save(&update)
	context.Status(http.StatusNoContent)
}

func DeletePlayer(context *gin.Context) {
	id := context.Param("id")
	var player models.Player
	db := data.DB
	// https://gorm.io/docs/delete.html
	result := db.Delete(&player, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			context.Status(http.StatusNotFound)
			return
		}
	}
	context.Status(http.StatusNoContent)
}

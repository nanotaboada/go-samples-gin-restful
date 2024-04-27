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

// Creates a new Player
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

// Retrieves all Players
func GetPlayers(context *gin.Context) {
	var players []models.Player
	db := data.DB
	// https://gorm.io/docs/query.html
	db.Find(&players)
	context.IndentedJSON(http.StatusOK, players)
}

// Retrieves a Player by Id
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

// Updates a Player by Id
func UpdatePlayer(context *gin.Context) {
	id := context.Param("id")
	var player models.Player
	var update models.Player
	if err := context.BindJSON(&update); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	db := data.DB
	// https://gorm.io/docs/query.html
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

// Deletes a Player by Id
func DeletePlayer(context *gin.Context) {
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
	// https://gorm.io/docs/delete.html
	db.Delete(&models.Player{}, id)
	context.Status(http.StatusNoContent)
}

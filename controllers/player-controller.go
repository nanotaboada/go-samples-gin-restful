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

// Post creates a new Player
func Post(context *gin.Context) {
	var player models.Player
	if err := context.BindJSON(&player); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	db := data.DB
	// https://gorm.io/docs/query.html
	result := db.First(&player, player.ID)
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

// GetAll retrieves all players
func GetAll(context *gin.Context) {
	var players []models.Player
	db := data.DB
	// https://gorm.io/docs/query.html
	db.Find(&players)
	context.IndentedJSON(http.StatusOK, players)
}

// GetByID retrieves a Player by ID
func GetByID(context *gin.Context) {
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

// Put updates (entirely) a Player by ID
func Put(context *gin.Context) {
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

// Delete removes a Player by ID
func Delete(context *gin.Context) {
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

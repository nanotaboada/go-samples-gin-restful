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

// Post creates a Player
//
// @Summary Creates a Player
// @Tags players
// @Accept application/json
// @Param player body models.Player true "Player"
// @Success 201 {object} models.Player "Created"
// @Failure 400 "Bad Request"
// @Failure 409 "Conflict"
// @Router /players [post]
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
//
// @Summary Retrieves all players
// @Tags players
// @Produce application/json
// @Success 200 {array} models.Player "OK"
// @Router /players [get]
func GetAll(context *gin.Context) {
	var players []models.Player
	db := data.DB
	// https://gorm.io/docs/query.html
	db.Find(&players)
	context.IndentedJSON(http.StatusOK, players)
}

// GetByID retrieves a Player by its ID
//
// @Summary Retrieves a Player by its ID
// @Tags players
// @Produce application/json
// @Param id path string true "Player.ID"
// @Success 200 {object} models.Player "OK"
// @Failure 404 "Not Found"
// @Router /players/{id} [get]
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

// Put updates (entirely) a Player by its ID
//
// @Summary Updates (entirely) a Player by its ID
// @Tags players
// @Accept application/json
// @Param id path string true "Player.ID"
// @Param player body models.Player true "Player"
// @Success 204 "No Content"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Router /players/{id} [put]
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

// Delete deletes a Player by its ID
//
// @Summary Deletes a Player by its ID
// @Tags players
// @Param id path string true "Player.ID"
// @Success 204 "No Content"
// @Failure 404 "Not Found"
// @Router /players/{id} [delete]
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

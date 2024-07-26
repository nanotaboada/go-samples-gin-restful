/* -----------------------------------------------------------------------------
 * Controllers
 * -------------------------------------------------------------------------- */

package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/models"
	"github.com/nanotaboada/go-samples-gin-restful/services"
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
	_, err := services.RetrieveByID(player.ID)
	if err == nil {
		context.Status(http.StatusConflict)
		return
	}
	if err := services.Create(&player); err == nil {
		context.Status(http.StatusCreated)
		return
	}
}

// GetAll retrieves all players
//
// @Summary Retrieves all players
// @Tags players
// @Produce application/json
// @Success 200 {array} models.Player "OK"
// @Router /players [get]
func GetAll(context *gin.Context) {
	players, _ := services.RetrieveAll()
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
	id, _ := strconv.Atoi(context.Param("id"))
	player, err := services.RetrieveByID(id)
	if err != nil {
		context.Status(http.StatusNotFound)
		return
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
	id, _ := strconv.Atoi(context.Param("id"))
	_, err := services.RetrieveByID(id)
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}
	var player models.Player
	if err := context.BindJSON(&player); err != nil || player.ID != id {
		context.Status(http.StatusBadRequest)
		return
	}
	if err := services.Update(&player); err == nil {
		context.Status(http.StatusNoContent)
		return
	}
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
	id, _ := strconv.Atoi(context.Param("id"))
	_, err := services.RetrieveByID(id)
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}
	if err := services.Delete(id); err == nil {
		context.Status(http.StatusNoContent)
		return
	}
}

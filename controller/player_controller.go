// Package controller defines the HTTP handlers for Player-related endpoints.
package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/model"
	"github.com/nanotaboada/go-samples-gin-restful/service"
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
	var player model.Player
	if err := context.BindJSON(&player); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	_, err := service.RetrieveByID(player.ID)
	if err == nil {
		context.Status(http.StatusConflict)
		return
	}
	if err := service.Create(&player); err == nil {
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
	players, _ := service.RetrieveAll()
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
	player, err := service.RetrieveByID(id)
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}
	context.IndentedJSON(http.StatusOK, player)
}

// GetByID retrieves a Player by its Squad Number
//
// @Summary Retrieves a Player by its Squad Number
// @Tags players
// @Produce application/json
// @Param squadnumber path string true "Player.SquadNumber"
// @Success 200 {object} models.Player "OK"
// @Failure 404 "Not Found"
// @Router /players/squadnumber/{squadnumber} [get]
func GetBySquadNumber(context *gin.Context) {
	squadNumber, _ := strconv.Atoi(context.Param("squadnumber"))
	player, err := service.RetrieveBySquadNumber(squadNumber)
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
	_, err := service.RetrieveByID(id)
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}
	var player model.Player
	if err := context.BindJSON(&player); err != nil || player.ID != id {
		context.Status(http.StatusBadRequest)
		return
	}
	if err := service.Update(&player); err == nil {
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
	_, err := service.RetrieveByID(id)
	if err != nil {
		context.Status(http.StatusNotFound)
		return
	}
	if err := service.Delete(id); err == nil {
		context.Status(http.StatusNoContent)
		return
	}
}

// Package controller defines the HTTP handlers for Player-related endpoints.
//
// In Gin, each handler is a function that receives a *gin.Context, which
// bundles the HTTP request (path params, query params, body, headers) and the
// HTTP response writer into a single object.  Handlers are registered on a
// *gin.Engine (the router) and called by the framework when a matching request
// arrives.
//
// This package follows the constructor-injection pattern: PlayerController
// holds its service dependency as a field, set once by NewPlayerController.
// Using an interface type (PlayerService) rather than the concrete struct
// keeps the controller testable — tests can supply a mock that implements the
// same interface without touching the database.
package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nanotaboada/go-samples-gin-restful/model"
	"github.com/nanotaboada/go-samples-gin-restful/service"
	"gorm.io/gorm"
)

// PlayerController holds dependencies for player handlers.
// Exported so the route package can receive it as a parameter; the service
// field is unexported because nothing outside this package needs it.
type PlayerController struct {
	service service.PlayerService
}

// NewPlayerController returns a PlayerController wired to the given service.
// Callers (main.go, tests) pass different implementations of the interface,
// enabling dependency injection without a DI framework.
func NewPlayerController(service service.PlayerService) *PlayerController {
	return &PlayerController{service: service}
}

// Post creates a Player
//
// @Summary Creates a Player
// @Tags players
// @Accept application/json
// @Param player body model.Player true "Player"
// @Success 201 "Created"
// @Failure 400 "Bad Request"
// @Failure 409 "Conflict"
// @Failure 500 "Internal Server Error"
// @Router /players [post]
func (c *PlayerController) Post(context *gin.Context) {
	var player model.Player
	// BindJSON deserialises the request body into the struct and validates
	// required fields declared with the `binding:"required"` tag.  On failure
	// it writes a 400 response automatically; we still return to stop execution.
	if err := context.BindJSON(&player); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	// UUID is always generated server-side; any client-provided ID is overwritten.
	// uuid.NewString() returns a random UUID v4 string (e.g. "6ba7b810-...").
	player.ID = uuid.NewString()
	// Conflict is checked by squadNumber (the user-facing unique identifier).
	// If RetrieveBySquadNumber returns nil error, the squad number is taken → 409.
	_, err := c.service.RetrieveBySquadNumber(player.SquadNumber)
	if err == nil {
		context.Status(http.StatusConflict)
		return
	}
	// errors.Is unwraps error chains, so it works even if the service wraps
	// gorm.ErrRecordNotFound in another error.  Any error other than "not found"
	// is an unexpected DB failure → 500.
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		context.Status(http.StatusInternalServerError)
		return
	}
	if err := c.service.Create(&player); err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	context.Status(http.StatusCreated)
}

// GetAll retrieves all players
//
// @Summary Retrieves all players
// @Tags players
// @Produce application/json
// @Success 200 {array} model.Player "OK"
// @Failure 500 "Internal Server Error"
// @Router /players [get]
func (c *PlayerController) GetAll(context *gin.Context) {
	players, err := c.service.RetrieveAll()
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	// IndentedJSON writes a pretty-printed JSON body with the given status code.
	// Use context.JSON for compact output in production if payload size matters.
	context.IndentedJSON(http.StatusOK, players)
}

// GetByID retrieves a Player by its internal UUID (admin use only)
//
// @Summary Retrieves a Player by its internal UUID (admin only)
// @Tags players
// @Produce application/json
// @Param id path string true "Player.ID (UUID)"
// @Success 200 {object} model.Player "OK"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /players/uuid/{id} [get]
func (c *PlayerController) GetByID(context *gin.Context) {
	// context.Param reads a named route parameter defined with ":name" syntax.
	// For the route "/players/uuid/:id", context.Param("id") returns the
	// value captured from the URL.
	id := context.Param("id")
	player, err := c.service.RetrieveByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.Status(http.StatusNotFound)
		} else {
			context.Status(http.StatusInternalServerError)
		}
		return
	}
	context.IndentedJSON(http.StatusOK, player)
}

// GetBySquadNumber retrieves a Player by its Squad Number
//
// @Summary Retrieves a Player by its Squad Number
// @Tags players
// @Produce application/json
// @Param squadnumber path string true "Player.SquadNumber"
// @Success 200 {object} model.Player "OK"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /players/squadnumber/{squadnumber} [get]
func (c *PlayerController) GetBySquadNumber(context *gin.Context) {
	// Route parameters are always strings; strconv.Atoi converts to int.
	// A non-numeric value (e.g. "/players/squadnumber/abc") returns an error → 400.
	squadNumber, err := strconv.Atoi(context.Param("squadnumber"))
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	player, err := c.service.RetrieveBySquadNumber(squadNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.Status(http.StatusNotFound)
		} else {
			context.Status(http.StatusInternalServerError)
		}
		return
	}
	context.IndentedJSON(http.StatusOK, player)
}

// Put updates (entirely) a Player by its Squad Number
//
// @Summary Updates (entirely) a Player by its Squad Number
// @Tags players
// @Accept application/json
// @Param squadnumber path int true "Player.SquadNumber"
// @Param player body model.Player true "Player"
// @Success 204 "No Content"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /players/{squadnumber} [put]
func (c *PlayerController) Put(context *gin.Context) {
	squadNumber, err := strconv.Atoi(context.Param("squadnumber"))
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	var player model.Player
	// Guard against mismatched URL and body: the squad number in the URL must
	// equal the one in the JSON body, otherwise the request is ambiguous → 400.
	if err := context.BindJSON(&player); err != nil || player.SquadNumber != squadNumber {
		context.Status(http.StatusBadRequest)
		return
	}
	existing, err := c.service.RetrieveBySquadNumber(squadNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.Status(http.StatusNotFound)
		} else {
			context.Status(http.StatusInternalServerError)
		}
		return
	}
	// Preserve the internal UUID — clients identify players by squadNumber, not UUID.
	// Without this, Save would try to zero out the primary key, causing a DB error.
	player.ID = existing.ID
	if err := c.service.Update(&player); err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	// 204 No Content is conventional for a successful PUT with no response body.
	context.Status(http.StatusNoContent)
}

// Delete deletes a Player by its Squad Number
//
// @Summary Deletes a Player by its Squad Number
// @Tags players
// @Param squadnumber path int true "Player.SquadNumber"
// @Success 204 "No Content"
// @Failure 400 "Bad Request"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
// @Router /players/{squadnumber} [delete]
func (c *PlayerController) Delete(context *gin.Context) {
	squadNumber, err := strconv.Atoi(context.Param("squadnumber"))
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	// Fetch first so GORM has a populated struct (including the primary key)
	// before issuing the DELETE statement; deleting by struct avoids an
	// unintended "DELETE FROM players WHERE id = 0" on a zero-value struct.
	existing, err := c.service.RetrieveBySquadNumber(squadNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.Status(http.StatusNotFound)
		} else {
			context.Status(http.StatusInternalServerError)
		}
		return
	}
	if err := c.service.Delete(&existing); err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	context.Status(http.StatusNoContent)
}

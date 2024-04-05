package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/models"
	"github.com/nanotaboada/go-samples-gin-restful/routes"
	"github.com/stretchr/testify/assert"
)

func TestMain(main *testing.M) {

	gin.SetMode(gin.TestMode)
	path := "../data/players-sqlite3.db"
	data.Connect(path)
	os.Exit(main.Run())
}

const path = "/players/"

// Given GET, when request to /players (no trailing slash), then response status should be 301 (Moved Permanently)
func TestGetPlayersNoTrailingSlashStatusMovedPermanently(test *testing.T) {

	// Arrange
	router := routes.Setup()
	request, _ := http.NewRequest("GET", "/players", nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusMovedPermanently, recorder.Code)
}

// Given GET, when request has no parameters, then response status should be 200 (OK).
func TestRequestGetPlayersResponseStatusOK(test *testing.T) {

	// Arrange
	router := routes.Setup()
	request, _ := http.NewRequest("GET", path, nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

// Given GET, when request has no parameters, then response body should be the collection of players.
func TestRequestGetPlayersResponsePlayers(test *testing.T) {

	// Arrange
	router := routes.Setup()
	request, _ := http.NewRequest("GET", path, nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)
	var players []models.Player
	json.Unmarshal(recorder.Body.Bytes(), &players)

	// Assert
	assert.NotEmpty(test, players)
}

// Given GET, when request parameter identifies existing player, then response status should be 200 (OK).
func TestRequestGetPlayersIdResponseStatusOK(test *testing.T) {

	// Arrange
	id := "10"
	router := routes.Setup()
	request, _ := http.NewRequest("GET", path+id, nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

// Given GET, when request parameter identifies existing player, then response body should be matching Player.
func TestRequestGetPlayersIdResponsePlayer(test *testing.T) {

	// Arrange
	id := "10"
	router := routes.Setup()
	request, _ := http.NewRequest("GET", path+id, nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)
	var player models.Player
	json.Unmarshal(recorder.Body.Bytes(), &player)

	// Assert
	assert.NotEmpty(test, player)
	assert.Equal(test, 10, player.SquadNumber)
	assert.Equal(test, "Lionel", player.FirstName)
	assert.Equal(test, "Messi", player.LastName)
}

// Given GET, when request parameter does not identify a player, then response status should be 404 (Not Found).
func TestRequestGetPlayersIdResponseStatusNotFound(test *testing.T) {

	// Arrange
	id := "99"
	router := routes.Setup()
	request, _ := http.NewRequest("GET", path+id, nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

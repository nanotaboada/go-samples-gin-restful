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

func Test_GivenHTTPGET_WhenRequestHasNoParameter_ThenResponseBodyShouldBeAllPlayers(test *testing.T) {

	// Arrange
	router := routes.Setup()
	request, _ := http.NewRequest("GET", "/players", nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)
	var players []models.Player
	json.Unmarshal(recorder.Body.Bytes(), &players)

	// Assert
	assert.NotEmpty(test, players)
	assert.Equal(test, http.StatusOK, recorder.Code)
}

func Test_GivenHTTPGET_WhenRequestParameterIdentifiesExistingPlayer_ThenResponseCodeShouldBeStatusOK(test *testing.T) {

	// Arrange
	id := "10"
	router := routes.Setup()
	request, _ := http.NewRequest("GET", "/players/"+id, nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

func Test_GivenHTTPGET_WhenRequestParameterIdentifiesExistingPlayer_ThenResponseBodyShouldBeThePlayer(test *testing.T) {

	// Arrange
	id := "10"
	router := routes.Setup()
	request, _ := http.NewRequest("GET", "/players/"+id, nil)
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

func Test_GivenHTTPGET_WhenRequestParameterDoesNotIdentifyExistingPlayer_ThenResponseCodeShouldBeStatusNotFound(test *testing.T) {

	// Arrange
	id := "99"
	router := routes.Setup()
	request, _ := http.NewRequest("GET", "/players/"+id, nil)
	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

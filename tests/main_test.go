/* -----------------------------------------------------------------------------
 * Tests
 * -------------------------------------------------------------------------- */

package tests

import (
	"bytes"
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
	data.Connect("../data/players-sqlite3.db")
	os.Exit(main.Run())
}

const (
	Url             = "/players/"
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
)

// Given GET, when request to /players (no trailing slash), then response status should be 301 (Moved Permanently)
func TestRequestGetNoTrailingSlashResponseStatusMovedPermanently(test *testing.T) {
	// Arrange
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/players", nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusMovedPermanently, recorder.Code)
}

// Given GET, when request has no parameters, then response status should be 200 (OK).
func TestRequestGetNoParamResponseStatusOK(test *testing.T) {
	// Arrange
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, Url, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

// Given GET, when request has no parameters, then response body should be the collection of Players.
func TestRequestGetNoParamResponsePlayers(test *testing.T) {
	// Arrange
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, Url, nil)

	// Act
	router.ServeHTTP(recorder, request)
	var players []models.Player
	json.Unmarshal(recorder.Body.Bytes(), &players)

	// Assert
	assert.NotEmpty(test, players)
}

// Given GET, when request parameter identifies existing Player, then response status should be 200 (OK).
func TestRequestGETIdExistingPlayerResponseStatusOK(test *testing.T) {
	// Arrange
	id := "10"
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, Url+id, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

// Given GET, when request parameter identifies existing Player, then response body should be matching Player.
func TestRequestGETIdExistingPlayerResponsePlayer(test *testing.T) {
	// Arrange
	id := "10"
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, Url+id, nil)

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

// Given GET, when request parameter does not identify a Player, then response status should be 404 (Not Found).
func TestRequestGETIdNonExistingPlayerResponseStatusNotFound(test *testing.T) {
	// Arrange
	id := "999"
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, Url+id, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// Given POST, when request body is empty, then response status should be 400 (Bad Request)
func TestRequestPOSTBodyEmptyResponseStatusBadRequest(test *testing.T) {
	// Arrange
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, Url, nil)
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// Given POST, when request is existing Player, then response status should be 409 (Conflict)
func TestRequestPOSTBodyExistingPlayerResponseStatusConflict(test *testing.T) {
	// Arrange
	player := GetExistingPlayer()
	body, _ := json.Marshal(player)
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, Url, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusConflict, recorder.Code)
}

// Given POST, when request is non-existing Player, then response status should be 201 (Created)
func TestRequestPOSTBodyNonExistingPlayerResponseStatusCreated(test *testing.T) {
	// Arrange
	player := GetNonExistingPlayer()
	body, _ := json.Marshal(player)
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, Url, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusCreated, recorder.Code)
}

// Given PUT, when request body is empty, then response status should be 400 (Bad Request)
func TestRequestPUTBodyEmptyResponseStatusBadRequest(test *testing.T) {
	// Arrange
	id := "10"
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, Url+id, nil)
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// Given PUT, when request is non-existing Player, then response status should be 404 (Not Found)
func TestRequestPUTBodyNonExistingPlayerResponseStatusNotFound(test *testing.T) {
	// Arrange
	id := "999"
	player := models.Player{
		FirstName: "John",
		LastName:  "Doe",
	}
	body, _ := json.Marshal(player)
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, Url+id, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// Given PUT, when request is existing Player, then response status should be 204 (No Content)
func TestRequestPUTPBodyExistingPlayerResponseStatusNoContent(test *testing.T) {
	// Arrange
	id := "1"
	player := GetExistingPlayer()
	player.FirstName = "Emiliano"
	player.MiddleName = ""
	body, _ := json.Marshal(player)
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, Url+id, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNoContent, recorder.Code)
}

// Given DELETE, when request is non-existing Player, then response status should be 404 (Not Found)
func TestRequestDELETEIdNonExistingPlayerResponseStatusNotFound(test *testing.T) {
	// Arrange
	id := "999"
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, Url+id, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// Given DELETE, when request is existing Player, then response status should be 204 (No Content)
func TestRequestDELETEIdExistingPlayerResponseStatusNoContent(test *testing.T) {
	// Arrange
	id := "12"
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, Url+id, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNoContent, recorder.Code)
}

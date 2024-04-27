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
	URL             = "/players/"
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
)

/* POST /players/ ----------------------------------------------------------- */

// Given POST
// When request body is empty
// Then response status should be 400 (Bad Request)
func TestRequestPOSTBodyEmptyResponseStatusBadRequest(test *testing.T) {
	// Arrange
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, URL, nil)
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// Given POST
// When request body is existing Player
// Then response status should be 409 (Conflict)
func TestRequestPOSTBodyExistingPlayerResponseStatusConflict(test *testing.T) {
	// Arrange
	player := GetExistingPlayer()
	body, _ := json.Marshal(player)
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusConflict, recorder.Code)
}

// Given POST
// when request body is non-existing Player
// Then response status should be 201 (Created)
func TestRequestPOSTBodyNonExistingPlayerResponseStatusCreated(test *testing.T) {
	// Arrange
	player := GetNonExistingPlayer()
	body, _ := json.Marshal(player)
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusCreated, recorder.Code)
}

/* GET /players/ ------------------------------------------------------------ */

// Given GET
// When request to /players (no trailing slash)
// Then response status should be 301 (Moved Permanently)
func TestRequestGETNoTrailingSlashResponseStatusMovedPermanently(test *testing.T) {
	// Arrange
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/players", nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusMovedPermanently, recorder.Code)
}

// Given GET
// When request path has no Id
// Then response status should be 200 (OK)
func TestRequestGETNoParamResponseStatusOK(test *testing.T) {
	// Arrange
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, URL, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

// Given GET
// When request path has no Id
// Then response body should be collection of Players
func TestRequestGETNoParamResponsePlayers(test *testing.T) {
	// Arrange
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, URL, nil)

	// Act
	router.ServeHTTP(recorder, request)
	var players []models.Player
	json.Unmarshal(recorder.Body.Bytes(), &players)

	// Assert
	assert.NotEmpty(test, players)
}

/* GET /players/:id --------------------------------------------------------- */

// Given GET
// When request path is non-existing Id
// Then response status should be 404 (Not Found)
func TestRequestGETIdNonExistingResponseStatusNotFound(test *testing.T) {
	// Arrange
	id := "999"
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, URL+id, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// Given GET
// When request path is existing Id
// Then response status should be 200 (OK)
func TestRequestGETIdExistingResponseStatusOK(test *testing.T) {
	// Arrange
	id := "10"
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, URL+id, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

// Given GET
// When request path is existing Id
// Then response body should be matching Player
func TestRequestGETIdExistingResponsePlayer(test *testing.T) {
	// Arrange
	id := "10"
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, URL+id, nil)

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

/* PUT /players/:id --------------------------------------------------------- */

// Given PUT
// When request body is empty
// Then response status should be 400 (Bad Request)
func TestRequestPUTBodyEmptyResponseStatusBadRequest(test *testing.T) {
	// Arrange
	id := "10"
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, URL+id, nil)
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// Given PUT
// When request body is unknown Player
// Then response status should be 404 (Not Found)
func TestRequestPUTBodyUnknownPlayerResponseStatusNotFound(test *testing.T) {
	// Arrange
	id := "999"
	player := models.Player{
		FirstName: "John",
		LastName:  "Doe",
	}
	body, _ := json.Marshal(player)
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, URL+id, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// Given PUT
// When request body is existing Player
// Then response status should be 204 (No Content)
func TestRequestPUTPBodyExistingPlayerResponseStatusNoContent(test *testing.T) {
	// Arrange
	id := "1"
	player := GetExistingPlayer()
	player.FirstName = "Emiliano"
	player.MiddleName = ""
	body, _ := json.Marshal(player)
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, URL+id, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNoContent, recorder.Code)
}

/* DELETE /players/:id ------------------------------------------------------ */

// Given DELETE
// When request path is non-existing Id
// Then response status should be 404 (Not Found)
func TestRequestDELETEIdNonExistingIdResponseStatusNotFound(test *testing.T) {
	// Arrange
	id := "999"
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, URL+id, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// Given DELETE
// when request path is existing Id
// Then response status should be  204 (No Content)
func TestRequestDELETEIdExistingIdResponseStatusNoContent(test *testing.T) {
	// Arrange
	id := "12"
	router := routes.Setup()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, URL+id, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNoContent, recorder.Code)
}

// Package tests provides integration and utility code to support automated
// testing of the application.
package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/model"
	"github.com/nanotaboada/go-samples-gin-restful/route"
	"github.com/stretchr/testify/assert"
)

func TestMain(main *testing.M) {
	gin.SetMode(gin.TestMode)
	data.Connect("file::memory:?cache=shared")
	if err := data.DB.AutoMigrate(&model.Player{}); err != nil {
		log.Fatal(err)
	}
	players, err := MakePlayersFromJSON()
	if err != nil {
		log.Fatal(err)
	}
	if err := data.DB.Create(&players).Error; err != nil {
		log.Fatal(err)
	}
	os.Exit(main.Run())
}

// setupRouter creates a router with all routes for testing
func setupRouter() *gin.Engine {
	store := persistence.NewInMemoryStore(time.Hour)
	router := gin.Default()

	route.RegisterPlayerRoutes(router, store)

	router.GET(route.HealthPath, func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}

const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
)

/* GET /health -------------------------------------------------------------- */

// Given GET
// When request
// Then response status should be 200 (OK)
func TestRequestGETHealthResponseStatusOK(test *testing.T) {
	// Arrange
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/health", nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

/* POST /players/ ----------------------------------------------------------- */

// Given POST
// When request body is empty
// Then response status should be 400 (Bad Request)
func TestRequestPOSTBodyEmptyResponseStatusBadRequest(test *testing.T) {
	// Arrange
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, route.GetAllPath, nil)
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
	player := MakeExistingPlayer()
	body, _ := json.Marshal(player)
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, route.GetAllPath, bytes.NewBuffer(body))
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
	player := MakeNonExistingPlayer()
	body, _ := json.Marshal(player)
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, route.GetAllPath, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusCreated, recorder.Code)
}

// Given POST
// When request to /players/ (with trailing slash) and body is empty
// Then response status should be 400 (Bad Request)
func TestRequestPOSTTrailingSlashBodyEmptyResponseStatusBadRequest(test *testing.T) {
	// Arrange
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, route.GetAllPathTrailingSlash, nil)
	if err != nil {
		test.Fatalf("failed to create POST request: %v", err)
	}
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

/* GET /players/ ------------------------------------------------------------ */

// Given GET
// When request to /players/ (with trailing slash)
// Then response status should be 200 (OK) - both routes are handled directly
func TestRequestGETTrailingSlashResponseStatusOK(test *testing.T) {
	// Arrange
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, route.GetAllPathTrailingSlash, nil)
	if err != nil {
		test.Fatalf("failed to create GET request: %v", err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

// Given GET
// When request path has no Id
// Then response status should be 200 (OK)
func TestRequestGETNoParamResponseStatusOK(test *testing.T) {
	// Arrange
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.GetAllPath, nil)

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
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.GetAllPath, nil)

	// Act
	router.ServeHTTP(recorder, request)
	var players []model.Player
	json.Unmarshal(recorder.Body.Bytes(), &players)

	// Assert
	assert.NotEmpty(test, players)
}

/* GET /players/:id --------------------------------------------------------- */

// Given GET
// When request path is non-existing Squad Number
// Then response status should be 404 (Not Found)
func TestRequestGETSquadNumberNonExistingResponseStatusNotFound(test *testing.T) {
	// Arrange
	squadNumber := "999"
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+route.SquadNumberParam+"/"+squadNumber, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// Given GET
// When request path is existing Squad Number
// Then response status should be 200 (OK)
func TestRequestGETSquadNumberExistingResponseStatusOK(test *testing.T) {
	// Arrange
	squadNumber := "11"
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+route.SquadNumberParam+"/"+squadNumber, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

// Given GET
// When request path is existing Squad Number
// Then response body should be matching Player
func TestRequestGETSquadNumberExistingResponsePlayer(test *testing.T) {
	// Arrange
	squadNumber := "11"
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+route.SquadNumberParam+"/"+squadNumber, nil)

	// Act
	router.ServeHTTP(recorder, request)
	var player model.Player
	json.Unmarshal(recorder.Body.Bytes(), &player)

	// Assert
	assert.NotEmpty(test, player)
	assert.Equal(test, 11, player.SquadNumber)
	assert.Equal(test, "Ángel", player.FirstName)
	assert.Equal(test, "Di María", player.LastName)
}

/* GET /players/:id --------------------------------------------------------- */

// Given GET
// When request path is non-existing Id
// Then response status should be 404 (Not Found)
func TestRequestGETIdNonExistingResponseStatusNotFound(test *testing.T) {
	// Arrange
	id := "999"
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+id, nil)

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
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+id, nil)

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
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+id, nil)

	// Act
	router.ServeHTTP(recorder, request)
	var player model.Player
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
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, route.PlayersPath+"/"+id, nil)
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
	player := model.Player{
		FirstName: "John",
		LastName:  "Doe",
	}
	body, _ := json.Marshal(player)
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, route.PlayersPath+"/"+id, bytes.NewBuffer(body))
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
	player := MakeExistingPlayer()
	player.FirstName = "Emiliano"
	player.MiddleName = ""
	body, _ := json.Marshal(player)
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, route.PlayersPath+"/"+id, bytes.NewBuffer(body))
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
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, route.PlayersPath+"/"+id, nil)

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
	router := setupRouter()
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, route.PlayersPath+"/"+id, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNoContent, recorder.Code)
}

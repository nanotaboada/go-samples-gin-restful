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
	"github.com/nanotaboada/go-samples-gin-restful/controller"
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/model"
	"github.com/nanotaboada/go-samples-gin-restful/route"
	"github.com/nanotaboada/go-samples-gin-restful/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	// testDB is global because tests share the same in-memory database instance
	testDB *gorm.DB
	// playerController is global because integration tests use it via setupRouter()
	playerController *controller.PlayerController
	// Note: playerService is local in TestMain since it's only needed to construct
	// the controller and is never referenced directly by tests
)

func TestMain(main *testing.M) {
	gin.SetMode(gin.TestMode)
	testDB = data.Connect("file::memory:?cache=shared")
	if err := testDB.AutoMigrate(&model.Player{}); err != nil {
		log.Fatal(err)
	}
	players, err := MakePlayersFromJSON()
	if err != nil {
		log.Fatal(err)
	}
	if err := testDB.Create(&players).Error; err != nil {
		log.Fatal(err)
	}
	// playerService is local - only used to initialize playerController,
	// then garbage collected
	playerService := service.NewPlayerService(testDB)
	playerController = controller.NewPlayerController(playerService)
	os.Exit(main.Run())
}

func setupRouter(controller *controller.PlayerController) *gin.Engine {
	store := persistence.NewInMemoryStore(time.Hour)
	app := gin.Default()
	route.RegisterPlayerRoutes(app, controller, store)
	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	return app
}

const (
	ContentType        = "Content-Type"
	ApplicationJSON    = "application/json"
	InvalidID          = "invalid-id"
	InvalidSquadNumber = "invalid-squadnumber"
)

/* GET /health -------------------------------------------------------------- */

// TestRequestGETHealthResponseStatusOK tests that a
// GET request to /health
// returns a 200 OK status, indicating the server is running and healthy.
func TestRequestGETHealthResponseStatusOK(test *testing.T) {
	// Arrange
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/health", nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

/* POST /players/ ----------------------------------------------------------- */

// TestRequestPOSTPlayersEmptyBodyResponseStatusBadRequest tests that a
// POST request to /players/ with an empty body
// returns a 400 Bad Request status.
func TestRequestPOSTPlayersEmptyBodyResponseStatusBadRequest(test *testing.T) {
	// Arrange
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, route.GetAllPath, nil)
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// TestRequestPOSTPlayersExistingPlayerResponseStatusConflict tests that a
// POST request to /players/ with an existing player
// returns a 409 Conflict status.
func TestRequestPOSTPlayersExistingPlayerResponseStatusConflict(test *testing.T) {
	// Arrange
	player := MakeExistingPlayer()
	body, _ := json.Marshal(player)
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, route.GetAllPath, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusConflict, recorder.Code)
}

// TestRequestPOSTPlayersNonExistingPlayerResponseStatusCreated tests that a
// POST request to /players/ with a non-existing player
// returns a 201 Created status.
func TestRequestPOSTPlayersNonExistingPlayerResponseStatusCreated(test *testing.T) {
	// Arrange
	player := MakeNonExistingPlayer()
	body, _ := json.Marshal(player)
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, route.GetAllPath, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusCreated, recorder.Code)
}

// TestRequestPOSTPlayersTrailingSlashEmptyBodyResponseStatusBadRequest tests that a
// POST request to /players/ (with trailing slash) and an empty body
// returns a 400 Bad Request status.
func TestRequestPOSTPlayersTrailingSlashEmptyBodyResponseStatusBadRequest(test *testing.T) {
	// Arrange
	router := setupRouter(playerController)
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

// TestRequestPOSTPlayersRetrieveErrorResponseStatusInternalServerError tests that a
// POST request to /players/ when service.RetrieveByID() returns an unexpected error
// returns a 500 Internal Server Error status.
func TestRequestPOSTPlayersRetrieveErrorResponseStatusInternalServerError(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveByIDFunc: func(id int) (model.Player, error) {
			return model.Player{}, ErrGenericError
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	player := MakeNonExistingPlayer()
	body, _ := json.Marshal(player)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, route.GetAllPath, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

// TestRequestPOSTPlayersCreateErrorResponseStatusInternalServerError tests that a
// POST request to /players/ when service.Create() returns an error
// returns a 500 Internal Server Error status.
func TestRequestPOSTPlayersCreateErrorResponseStatusInternalServerError(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveByIDFunc: func(id int) (model.Player, error) {
			return model.Player{}, gorm.ErrRecordNotFound
		},
		CreateFunc: func(player *model.Player) error {
			return ErrDatabaseFailure
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	player := MakeNonExistingPlayer()
	body, _ := json.Marshal(player)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, route.GetAllPath, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

/* GET /players/ ------------------------------------------------------------ */

// TestRequestGETPlayersTrailingSlashResponseStatusOK tests that a
// GET request to /players/ (with trailing slash)
// returns a 200 OK status.
func TestRequestGETPlayersTrailingSlashResponseStatusOK(test *testing.T) {
	// Arrange
	router := setupRouter(playerController)
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

// TestRequestGETPlayersResponseStatusOK tests that a
// GET request to /players/
// returns a 200 OK status.
func TestRequestGETPlayersResponseStatusOK(test *testing.T) {
	// Arrange
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.GetAllPath, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

// TestRequestGETPlayersResponsePlayers tests that a
// GET request to /players/
// returns a collection of Players.
func TestRequestGETPlayersResponsePlayers(test *testing.T) {
	// Arrange
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.GetAllPath, nil)

	// Act
	router.ServeHTTP(recorder, request)
	var players []model.Player
	json.Unmarshal(recorder.Body.Bytes(), &players)

	// Assert
	assert.NotEmpty(test, players)
}

// TestRequestGETPlayersRetrieveErrorResponseStatusInternalServerError tests that a
// GET request to /players/ when service.RetrieveAll() returns an unexpected error
// returns a 500 Internal Server Error status.
func TestRequestGETPlayersRetrieveErrorResponseStatusInternalServerError(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveAllFunc: func() ([]model.Player, error) {
			return nil, ErrDatabaseFailure
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.GetAllPath, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

/* GET /players/:id --------------------------------------------------------- */

// TestRequestGETPlayerByIDNonExistingResponseStatusNotFound tests that a
// GET request to /players/:id when the id does not exist
// returns a 404 Not Found status.
func TestRequestGETPlayerByIDNonExistingResponseStatusNotFound(test *testing.T) {
	// Arrange
	id := "999"
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+id, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// TestRequestGETPlayerByIDInvalidParamResponseStatusBadRequest tests that a
// GET request to /players/:id when the id is invalid (non-numeric)
// returns a 400 Bad Request status.
func TestRequestGETPlayerByIDInvalidParamResponseStatusBadRequest(test *testing.T) {
	// Arrange
	id := InvalidID
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+id, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// TestRequestGETPlayerByIDExistingResponseStatusOK tests that a
// GET request to /players/:id when the id exists
// returns a 200 OK status.
func TestRequestGETPlayerByIDExistingResponseStatusOK(test *testing.T) {
	// Arrange
	id := "10"
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+id, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

// TestRequestGETPlayerByIDExistingResponsePlayer tests that a
// GET request to /players/:id when the id exists
// returns a matching Player.
func TestRequestGETPlayerByIDExistingResponsePlayer(test *testing.T) {
	// Arrange
	id := "10"
	router := setupRouter(playerController)
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

// TestRequestGETPlayerByIDRetrieveErrorResponseStatusInternalServerError tests that a
// GET request to /players/:id when service.RetrieveByID() returns an unexpected error
// returns a 500 Internal Server Error status.
func TestRequestGETPlayerByIDRetrieveErrorResponseStatusInternalServerError(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveByIDFunc: func(id int) (model.Player, error) {
			return model.Player{}, ErrGenericError
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.PlayersPath+"/10", nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

/* GET /players/squadNumber/:squadNumber ------------------------------------ */

// TestRequestGETPlayerBySquadNumberNonExistingResponseStatusNotFound tests that a
// GET request to /players/squadNumber/:squadNumber when the squad number does not exist
// returns a 404 Not Found status.
func TestRequestGETPlayerBySquadNumberNonExistingResponseStatusNotFound(test *testing.T) {
	// Arrange
	squadNumber := "999"
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+route.SquadNumberParam+"/"+squadNumber, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// TestRequestGETPlayerBySquadNumberInvalidParamResponseStatusBadRequest tests that a
// GET request to /players/squadNumber/:squadNumber when the squad number is invalid (non-numeric)
// returns a 400 Bad Request status.
func TestRequestGETPlayerBySquadNumberInvalidParamResponseStatusBadRequest(test *testing.T) {
	// Arrange
	squadNumber := InvalidSquadNumber
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+route.SquadNumberParam+"/"+squadNumber, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// TestRequestGETPlayerBySquadNumberExistingResponseStatusOK tests that a
// GET request to /players/squadNumber/:squadNumber when the squad number exists
// returns a 200 OK status.
func TestRequestGETPlayerBySquadNumberExistingResponseStatusOK(test *testing.T) {
	// Arrange
	squadNumber := "11"
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+route.SquadNumberParam+"/"+squadNumber, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

// TestRequestGETPlayerBySquadNumberExistingResponsePlayer tests that a
// GET request to /players/squadNumber/:squadNumber when the squad number exists
// returns a matching Player.
func TestRequestGETPlayerBySquadNumberExistingResponsePlayer(test *testing.T) {
	// Arrange
	squadNumber := "11"
	router := setupRouter(playerController)
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

// TestRequestGETPlayerBySquadNumberRetrieveErrorResponseStatusInternalServerError tests that a
// GET request to /players/squadNumber/:squadNumber when service.RetrieveBySquadNumber() returns an unexpected error
// returns a 500 Internal Server Error status.
func TestRequestGETPlayerBySquadNumberRetrieveErrorResponseStatusInternalServerError(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveBySquadNumberFunc: func(squadNumber int) (model.Player, error) {
			return model.Player{}, ErrGenericError
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+route.SquadNumberParam+"/10", nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

/* PUT /players/:id --------------------------------------------------------- */

// TestRequestPUTPlayerByIDEmptyBodyResponseStatusBadRequest tests that a
// PUT request to /players/:id with an empty body
// returns a 400 Bad Request status.
func TestRequestPUTPlayerByIDEmptyBodyResponseStatusBadRequest(test *testing.T) {
	// Arrange
	id := "10"
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, route.PlayersPath+"/"+id, nil)
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// TestRequestPUTPlayerByIDUnknownPlayerResponseStatusNotFound tests that a
// PUT request to /players/:id when the player does not exist
// returns a 404 Not Found status.
func TestRequestPUTPlayerByIDUnknownPlayerResponseStatusNotFound(test *testing.T) {
	// Arrange
	id := "999"
	player := model.Player{
		FirstName: "John",
		LastName:  "Doe",
	}
	body, _ := json.Marshal(player)
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, route.PlayersPath+"/"+id, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// TestRequestPUTPlayerByIDInvalidParamResponseStatusBadRequest tests that a
// PUT request to /players/:id when the id is invalid (non-numeric)
// returns a 400 Bad Request status.
func TestRequestPUTPlayerByIDInvalidParamResponseStatusBadRequest(test *testing.T) {
	// Arrange
	id := InvalidID
	player := MakeExistingPlayer()
	body, _ := json.Marshal(player)
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, route.PlayersPath+"/"+id, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// TestRequestPUTPlayerByIDExistingPlayerResponseStatusNoContent tests that a
// PUT request to /players/:id with an existing player
// returns a 204 No Content status.
func TestRequestPUTPlayerByIDExistingPlayerResponseStatusNoContent(test *testing.T) {
	// Arrange
	id := "1"
	player := MakeExistingPlayer()
	player.FirstName = "Emiliano"
	player.MiddleName = ""
	body, _ := json.Marshal(player)
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, route.PlayersPath+"/"+id, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNoContent, recorder.Code)
}

// TestRequestPUTPlayerByIDRetrieveErrorResponseStatusInternalServerError tests that a
// PUT request to /players/:id when service.RetrieveByID() returns an unexpected error
// returns a 500 Internal Server Error status.
func TestRequestPUTPlayerByIDRetrieveErrorResponseStatusInternalServerError(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveByIDFunc: func(id int) (model.Player, error) {
			return model.Player{}, ErrGenericError
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	player := MakeExistingPlayer()
	body, _ := json.Marshal(player)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, route.PlayersPath+"/1", bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

// TestRequestPUTPlayerByIDUpdateErrorResponseStatusInternalServerError tests that a
// PUT request to /players/:id when service.Update() returns an error
// returns a 500 Internal Server Error status.
func TestRequestPUTPlayerByIDUpdateErrorResponseStatusInternalServerError(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveByIDFunc: func(id int) (model.Player, error) {
			return MakeExistingPlayer(), nil
		},
		UpdateFunc: func(player *model.Player) error {
			return ErrDatabaseFailure
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	player := MakeExistingPlayer()
	body, _ := json.Marshal(player)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, route.PlayersPath+"/1", bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

/* DELETE /players/:id ------------------------------------------------------ */

// TestRequestDELETEPlayerByIDNonExistingResponseStatusNotFound tests that a
// DELETE request to /players/:id when the id does not exist
// returns a 404 Not Found status.
func TestRequestDELETEPlayerByIDNonExistingResponseStatusNotFound(test *testing.T) {
	// Arrange
	id := "999"
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, route.PlayersPath+"/"+id, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// TestRequestDELETEPlayerByIDInvalidParamResponseStatusBadRequest tests that a
// DELETE request to /players/:id when the id is invalid (non-numeric)
// returns a 400 Bad Request status.
func TestRequestDELETEPlayerByIDInvalidParamResponseStatusBadRequest(test *testing.T) {
	// Arrange
	id := InvalidID
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, route.PlayersPath+"/"+id, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// TestRequestDELETEPlayerByIDExistingResponseStatusNoContent tests that a
// DELETE request to /players/:id when the id exists
// returns a 204 No Content status.
func TestRequestDELETEPlayerByIDExistingResponseStatusNoContent(test *testing.T) {
	// Arrange
	id := "12"
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, route.PlayersPath+"/"+id, nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNoContent, recorder.Code)
}

// TestRequestDELETEPlayerByIDRetrieveErrorResponseStatusInternalServerError tests that a
// DELETE request to /players/:id when service.RetrieveByID() returns an unexpected error
// returns a 500 Internal Server Error status.
func TestRequestDELETEPlayerByIDRetrieveErrorResponseStatusInternalServerError(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveByIDFunc: func(id int) (model.Player, error) {
			return model.Player{}, ErrGenericError
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, route.PlayersPath+"/1", nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

// TestRequestDELETEPlayerByIDDeleteErrorResponseStatusInternalServerError tests that a
// DELETE request to /players/:id when service.Delete() returns an error
// returns a 500 Internal Server Error status.
func TestRequestDELETEPlayerByIDDeleteErrorResponseStatusInternalServerError(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveByIDFunc: func(id int) (model.Player, error) {
			return MakeExistingPlayer(), nil
		},
		DeleteFunc: func(id int) error {
			return ErrDatabaseFailure
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, route.PlayersPath+"/1", nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

// Package tests provides integration and utility code to support automated
// testing of the application.
//
// # Test database seeding
//
// Most tests in this package are integration tests: they exercise the full
// route → controller → service → data stack using a real in-memory SQLite
// database seeded in [TestMain].
//
// The seeding flow is:
//
//	players.json → MakePlayersFromJSON() → testDB.Create() → in-memory SQLite
//
// After [TestMain] completes every test hits SQLite directly — the JSON file
// is used exactly once as a human-readable fixture source and is then
// discarded. This means schema constraints (e.g. the unique index on
// squadNumber) are enforced for real, not mocked.
//
// # Mock-assisted tests
//
// A subset of tests uses [MockPlayerService] to inject controlled error
// conditions (e.g. unexpected DB failures) that cannot be triggered naturally
// with a healthy in-memory SQLite. These are controller-level tests in
// disguise: the mock replaces only the service layer so the controller's error
// handling branches can be reached and verified.
//
// Once a dedicated player_service_test.go exists, the mock-assisted
// tests should be moved there and rewritten as pure unit tests.
package tests

import (
	"bytes"
	"encoding/json"
	"errors"
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
	ErrNewRequest      = "failed to create request: %v"
	ErrMarshal         = "failed to marshal player: %v"
	ErrUnmarshal       = "failed to unmarshal response body: %v"
)

/* GET /health -------------------------------------------------------------- */

// TestRequestGETHealthResponseStatusOK tests that a
// GET request to /health
// returns a 200 OK status, indicating the server is running and healthy.
func TestRequestGETHealthResponseStatusOK(test *testing.T) {
	// Arrange
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

/* POST /players/ ----------------------------------------------------------- */

// TestRequestPOSTPlayersEmptyBodyResponseStatusBadRequest tests that a
// POST request to /players with an empty body
// returns a 400 Bad Request status.
func TestRequestPOSTPlayersEmptyBodyResponseStatusBadRequest(test *testing.T) {
	// Arrange
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, route.GetAllPath, nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// TestRequestPOSTPlayersExistingResponseStatusConflict tests that a
// POST request to /players with an existing player
// returns a 409 Conflict status.
func TestRequestPOSTPlayersExistingResponseStatusConflict(test *testing.T) {
	// Arrange
	player := MakeExistingPlayer()
	body, err := json.Marshal(player)
	if err != nil {
		test.Fatalf(ErrMarshal, err)
	}
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, route.GetAllPath, bytes.NewBuffer(body))
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusConflict, recorder.Code)
}

// TestRequestPOSTPlayersNonExistingResponseStatusCreated tests that a
// POST request to /players with a non-existing player
// returns a 201 Created status.
func TestRequestPOSTPlayersNonExistingResponseStatusCreated(test *testing.T) {
	// Arrange
	player := MakeNonExistingPlayer()
	body, err := json.Marshal(player)
	if err != nil {
		test.Fatalf(ErrMarshal, err)
	}
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, route.GetAllPath, bytes.NewBuffer(body))
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}
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
		test.Fatalf(ErrNewRequest, err)
	}
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// TestRequestPOSTPlayersRetrieveErrorResponseStatusInternalServerError tests that a
// POST request to /players when service.RetrieveBySquadNumber() returns an unexpected error
// returns a 500 Internal Server Error status.
func TestRequestPOSTPlayersRetrieveErrorResponseStatusInternalServerError(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveBySquadNumberFunc: func(squadNumber int) (model.Player, error) {
			return model.Player{}, ErrGenericError
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	player := MakeNonExistingPlayer()
	body, err := json.Marshal(player)
	if err != nil {
		test.Fatalf(ErrMarshal, err)
	}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, route.GetAllPath, bytes.NewBuffer(body))
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

// TestRequestPOSTPlayersCreateErrorResponseStatusInternalServerError tests that a
// POST request to /players when service.Create() returns an error
// returns a 500 Internal Server Error status.
func TestRequestPOSTPlayersCreateErrorResponseStatusInternalServerError(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveBySquadNumberFunc: func(squadNumber int) (model.Player, error) {
			return model.Player{}, gorm.ErrRecordNotFound
		},
		CreateFunc: func(player *model.Player) error {
			return ErrDatabaseFailure
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	player := MakeNonExistingPlayer()
	body, err := json.Marshal(player)
	if err != nil {
		test.Fatalf(ErrMarshal, err)
	}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, route.GetAllPath, bytes.NewBuffer(body))
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

// TestRequestPOSTPlayersCreateErrorResponseStatusConflict tests that a
// POST request to /players when service.Create() returns a unique constraint
// error (concurrent insert race) returns a 409 Conflict status.
func TestRequestPOSTPlayersCreateErrorResponseStatusConflict(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveBySquadNumberFunc: func(squadNumber int) (model.Player, error) {
			// Preflight check passes (squad number not found), simulating the
			// window between the read and the write where a concurrent request
			// inserts the same squadNumber, causing the subsequent Create to
			// violate the UNIQUE constraint.
			return model.Player{}, gorm.ErrRecordNotFound
		},
		CreateFunc: func(player *model.Player) error {
			return errors.New("UNIQUE constraint failed: players.squad_number")
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	player := MakeNonExistingPlayer()
	body, err := json.Marshal(player)
	if err != nil {
		test.Fatalf(ErrMarshal, err)
	}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPost, route.GetAllPath, bytes.NewBuffer(body))
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusConflict, recorder.Code)
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
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

// TestRequestGETPlayersResponseStatusOK tests that a
// GET request to /players
// returns a 200 OK status.
func TestRequestGETPlayersResponseStatusOK(test *testing.T) {
	// Arrange
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, route.GetAllPath, nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

// TestRequestGETPlayersResponsePlayers tests that a
// GET request to /players
// returns a collection of Players.
func TestRequestGETPlayersResponsePlayers(test *testing.T) {
	// Arrange
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, route.GetAllPath, nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)
	var players []model.Player
	if err := json.Unmarshal(recorder.Body.Bytes(), &players); err != nil {
		test.Fatalf(ErrUnmarshal, err)
	}

	// Assert
	assert.NotEmpty(test, players)
}

// TestRequestGETPlayersRetrieveErrorResponseStatusInternalServerError tests that a
// GET request to /players when service.RetrieveAll() returns an unexpected error
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
	request, err := http.NewRequest(http.MethodGet, route.GetAllPath, nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

/* GET /players/:id --------------------------------------------------------- */

// TestRequestGETPlayerByIDNonExistingResponseStatusNotFound tests that a
// GET request to /players/:id when the UUID does not exist
// returns a 404 Not Found status.
func TestRequestGETPlayerByIDNonExistingResponseStatusNotFound(test *testing.T) {
	// Arrange
	// Any valid UUID v4 that doesn't exist in the seeded data
	id := "00000000-0000-4000-8000-000000000000"
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+id, nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// TestRequestGETPlayerByIDExistingResponseStatusOK tests that a
// GET request to /players/:id when the UUID exists
// returns a 200 OK status.
func TestRequestGETPlayerByIDExistingResponseStatusOK(test *testing.T) {
	// Arrange
	// Squad #10 = Lionel Messi → UUID v5 derived from squad number 10
	id := "9a5fa2e4-9c9e-58e5-aeb3-8b1b46e87e03"
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+id, nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

// TestRequestGETPlayerByIDExistingResponsePlayer tests that a
// GET request to /players/:id when the UUID exists
// returns a matching Player.
func TestRequestGETPlayerByIDExistingResponsePlayer(test *testing.T) {
	// Arrange
	// Squad #10 = Lionel Messi → UUID v5 derived from squad number 10
	id := "9a5fa2e4-9c9e-58e5-aeb3-8b1b46e87e03"
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+id, nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)
	var player model.Player
	if err := json.Unmarshal(recorder.Body.Bytes(), &player); err != nil {
		test.Fatalf(ErrUnmarshal, err)
	}

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
		RetrieveByIDFunc: func(id string) (model.Player, error) {
			return model.Player{}, ErrGenericError
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, route.PlayersPath+"/9a5fa2e4-9c9e-58e5-aeb3-8b1b46e87e03", nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

/* GET /players/squadnumber/:squadnumber ------------------------------------ */

// TestRequestGETPlayerBySquadNumberNonExistingResponseStatusNotFound tests that a
// GET request to /players/squadnumber/:squadnumber when the squad number does not exist
// returns a 404 Not Found status.
func TestRequestGETPlayerBySquadNumberNonExistingResponseStatusNotFound(test *testing.T) {
	// Arrange
	squadNumber := "999"
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+route.SquadNumberParam+"/"+squadNumber, nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// TestRequestGETPlayerBySquadNumberInvalidParamResponseStatusBadRequest tests that a
// GET request to /players/squadnumber/:squadnumber when the squad number is invalid (non-numeric)
// returns a 400 Bad Request status.
func TestRequestGETPlayerBySquadNumberInvalidParamResponseStatusBadRequest(test *testing.T) {
	// Arrange
	squadNumber := InvalidSquadNumber
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+route.SquadNumberParam+"/"+squadNumber, nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// TestRequestGETPlayerBySquadNumberExistingResponseStatusOK tests that a
// GET request to /players/squadnumber/:squadnumber when the squad number exists
// returns a 200 OK status.
func TestRequestGETPlayerBySquadNumberExistingResponseStatusOK(test *testing.T) {
	// Arrange
	squadNumber := "11"
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+route.SquadNumberParam+"/"+squadNumber, nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusOK, recorder.Code)
}

// TestRequestGETPlayerBySquadNumberExistingResponsePlayer tests that a
// GET request to /players/squadnumber/:squadnumber when the squad number exists
// returns a matching Player.
func TestRequestGETPlayerBySquadNumberExistingResponsePlayer(test *testing.T) {
	// Arrange
	squadNumber := "11"
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+route.SquadNumberParam+"/"+squadNumber, nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)
	var player model.Player
	if err := json.Unmarshal(recorder.Body.Bytes(), &player); err != nil {
		test.Fatalf(ErrUnmarshal, err)
	}

	// Assert
	assert.NotEmpty(test, player)
	assert.Equal(test, 11, player.SquadNumber)
	assert.Equal(test, "Ángel", player.FirstName)
	assert.Equal(test, "Di María", player.LastName)
}

// TestRequestGETPlayerBySquadNumberRetrieveErrorResponseStatusInternalServerError tests that a
// GET request to /players/squadnumber/:squadnumber when service.RetrieveBySquadNumber() returns an unexpected error
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
	request, err := http.NewRequest(http.MethodGet, route.PlayersPath+"/"+route.SquadNumberParam+"/10", nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

/* PUT /players/:squadnumber ------------------------------------------------ */

// TestRequestPUTPlayerBySquadNumberEmptyBodyResponseStatusBadRequest tests that a
// PUT request to /players/:squadnumber with an empty body
// returns a 400 Bad Request status.
func TestRequestPUTPlayerBySquadNumberEmptyBodyResponseStatusBadRequest(test *testing.T) {
	// Arrange
	squadNumber := "23"
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPut, route.PlayersPath+"/"+squadNumber, nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// TestRequestPUTPlayerBySquadNumberNonExistingResponseStatusNotFound tests that a
// PUT request to /players/:squadnumber when the squad number does not exist
// returns a 404 Not Found status.
func TestRequestPUTPlayerBySquadNumberNonExistingResponseStatusNotFound(test *testing.T) {
	// Arrange
	squadNumber := "999"
	player := model.Player{
		SquadNumber: 999,
		FirstName:   "John",
		LastName:    "Doe",
	}
	body, err := json.Marshal(player)
	if err != nil {
		test.Fatalf(ErrMarshal, err)
	}
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPut, route.PlayersPath+"/"+squadNumber, bytes.NewBuffer(body))
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// TestRequestPUTPlayerBySquadNumberInvalidParamResponseStatusBadRequest tests that a
// PUT request to /players/:squadnumber when the squad number is invalid (non-numeric)
// returns a 400 Bad Request status.
func TestRequestPUTPlayerBySquadNumberInvalidParamResponseStatusBadRequest(test *testing.T) {
	// Arrange
	squadNumber := InvalidSquadNumber
	player := MakeExistingPlayer()
	body, err := json.Marshal(player)
	if err != nil {
		test.Fatalf(ErrMarshal, err)
	}
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPut, route.PlayersPath+"/"+squadNumber, bytes.NewBuffer(body))
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// TestRequestPUTPlayerBySquadNumberExistingResponseStatusNoContent tests that a
// PUT request to /players/:squadnumber with an existing player
// returns a 204 No Content status.
func TestRequestPUTPlayerBySquadNumberExistingResponseStatusNoContent(test *testing.T) {
	// Arrange
	squadNumber := "23"
	player := MakeExistingPlayer()
	player.FirstName = "Emiliano"
	player.MiddleName = ""
	body, err := json.Marshal(player)
	if err != nil {
		test.Fatalf(ErrMarshal, err)
	}
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPut, route.PlayersPath+"/"+squadNumber, bytes.NewBuffer(body))
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNoContent, recorder.Code)
}

// TestRequestPUTPlayerBySquadNumberMismatchSquadNumberResponseStatusBadRequest tests that a
// PUT request to /players/:squadnumber when the squad number in the URL
// does not match the one in the request body returns a 400 Bad Request status.
func TestRequestPUTPlayerBySquadNumberMismatchSquadNumberResponseStatusBadRequest(test *testing.T) {
	// Arrange
	player := MakeExistingPlayer() // SquadNumber == 23
	player.SquadNumber = 99        // mismatch: URL targets /players/23, body carries 99
	body, err := json.Marshal(player)
	if err != nil {
		test.Fatalf(ErrMarshal, err)
	}
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPut, route.PlayersPath+"/23", bytes.NewBuffer(body))
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// TestRequestPUTPlayerBySquadNumberRetrieveErrorResponseStatusInternalServerError tests that a
// PUT request to /players/:squadnumber when service.RetrieveBySquadNumber() returns an unexpected error
// returns a 500 Internal Server Error status.
func TestRequestPUTPlayerBySquadNumberRetrieveErrorResponseStatusInternalServerError(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveBySquadNumberFunc: func(squadNumber int) (model.Player, error) {
			return model.Player{}, ErrGenericError
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	player := MakeExistingPlayer()
	body, err := json.Marshal(player)
	if err != nil {
		test.Fatalf(ErrMarshal, err)
	}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPut, route.PlayersPath+"/23", bytes.NewBuffer(body))
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

// TestRequestPUTPlayerBySquadNumberUpdateErrorResponseStatusInternalServerError tests that a
// PUT request to /players/:squadnumber when service.Update() returns an error
// returns a 500 Internal Server Error status.
func TestRequestPUTPlayerBySquadNumberUpdateErrorResponseStatusInternalServerError(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveBySquadNumberFunc: func(squadNumber int) (model.Player, error) {
			return MakeExistingPlayer(), nil
		},
		UpdateFunc: func(player *model.Player) error {
			return ErrDatabaseFailure
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	player := MakeExistingPlayer()
	body, err := json.Marshal(player)
	if err != nil {
		test.Fatalf(ErrMarshal, err)
	}
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodPut, route.PlayersPath+"/23", bytes.NewBuffer(body))
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

/* DELETE /players/:squadnumber --------------------------------------------- */

// TestRequestDELETEPlayerBySquadNumberNonExistingResponseStatusNotFound tests that a
// DELETE request to /players/:squadnumber when the squad number does not exist
// returns a 404 Not Found status.
func TestRequestDELETEPlayerBySquadNumberNonExistingResponseStatusNotFound(test *testing.T) {
	// Arrange
	squadNumber := "999"
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodDelete, route.PlayersPath+"/"+squadNumber, nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// TestRequestDELETEPlayerBySquadNumberInvalidParamResponseStatusBadRequest tests that a
// DELETE request to /players/:squadnumber when the squad number is invalid (non-numeric)
// returns a 400 Bad Request status.
func TestRequestDELETEPlayerBySquadNumberInvalidParamResponseStatusBadRequest(test *testing.T) {
	// Arrange
	squadNumber := InvalidSquadNumber
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodDelete, route.PlayersPath+"/"+squadNumber, nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}

// TestRequestDELETEPlayerBySquadNumberExistingResponseStatusNoContent tests that a
// DELETE request to /players/:squadnumber when the squad number exists
// returns a 204 No Content status.
func TestRequestDELETEPlayerBySquadNumberExistingResponseStatusNoContent(test *testing.T) {
	// Arrange
	squadNumber := "1" // Franco Armani
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodDelete, route.PlayersPath+"/"+squadNumber, nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNoContent, recorder.Code)
}

// TestRequestDELETEPlayerBySquadNumberRetrieveErrorResponseStatusInternalServerError tests that a
// DELETE request to /players/:squadnumber when service.RetrieveBySquadNumber() returns an unexpected error
// returns a 500 Internal Server Error status.
func TestRequestDELETEPlayerBySquadNumberRetrieveErrorResponseStatusInternalServerError(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveBySquadNumberFunc: func(squadNumber int) (model.Player, error) {
			return model.Player{}, ErrGenericError
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodDelete, route.PlayersPath+"/23", nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

// TestRequestDELETEPlayerBySquadNumberDeleteErrorResponseStatusInternalServerError tests that a
// DELETE request to /players/:squadnumber when service.Delete() returns an error
// returns a 500 Internal Server Error status.
func TestRequestDELETEPlayerBySquadNumberDeleteErrorResponseStatusInternalServerError(test *testing.T) {
	// Arrange
	mockService := &MockPlayerService{
		RetrieveBySquadNumberFunc: func(squadNumber int) (model.Player, error) {
			return MakeExistingPlayer(), nil
		},
		DeleteFunc: func(player *model.Player) error {
			return ErrDatabaseFailure
		},
	}
	controller := controller.NewPlayerController(mockService)
	router := setupRouter(controller)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodDelete, route.PlayersPath+"/23", nil)
	if err != nil {
		test.Fatalf(ErrNewRequest, err)
	}

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusInternalServerError, recorder.Code)
}

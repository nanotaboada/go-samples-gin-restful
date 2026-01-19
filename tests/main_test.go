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
	// playerService is local - only used to initialize playerController, then garbage collected
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

// Given GET
// When request
// Then response status should be 200 (OK)
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

// Given POST
// When request body is empty
// Then response status should be 400 (Bad Request)
func TestRequestPOSTBodyEmptyResponseStatusBadRequest(test *testing.T) {
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

// Given POST
// When request body is existing Player
// Then response status should be 409 (Conflict)
func TestRequestPOSTBodyExistingPlayerResponseStatusConflict(test *testing.T) {
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

// Given POST
// when request body is non-existing Player
// Then response status should be 201 (Created)
func TestRequestPOSTBodyNonExistingPlayerResponseStatusCreated(test *testing.T) {
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

// Given POST
// When request to /players/ (with trailing slash) and body is empty
// Then response status should be 400 (Bad Request)
func TestRequestPOSTTrailingSlashBodyEmptyResponseStatusBadRequest(test *testing.T) {
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

// Given POST
// When service.RetrieveByID() fails with non-NotFound error
// Then response status should be 500 (Internal Server Error)
func TestRequestPOSTServiceRetrieveByIDGenericErrorResponseStatusInternalServerError(test *testing.T) {
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

// Given POST
// When service.Create() fails
// Then response status should be 500 (Internal Server Error)
func TestRequestPOSTServiceCreateFailureResponseStatusInternalServerError(test *testing.T) {
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

// Given GET
// When request to /players/ (with trailing slash)
// Then response status should be 200 (OK) - both routes are handled directly
func TestRequestGETTrailingSlashResponseStatusOK(test *testing.T) {
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

// Given GET
// When request path has no Id
// Then response status should be 200 (OK)
func TestRequestGETNoParamResponseStatusOK(test *testing.T) {
	// Arrange
	router := setupRouter(playerController)
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

// Given GET
// When service.RetrieveAll() fails
// Then response status should be 500 (Internal Server Error)
func TestRequestGETAllServiceRetrieveAllFailureResponseStatusInternalServerError(test *testing.T) {
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

// Given GET by Squad Number
// When service.RetrieveBySquadNumber() fails with non-NotFound error
// Then response status should be 500 (Internal Server Error)
func TestRequestGETBySquadNumberServiceGenericErrorResponseStatusInternalServerError(test *testing.T) {
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

/* GET /players/:id --------------------------------------------------------- */

// Given GET
// When request path is non-existing Squad Number
// Then response status should be 404 (Not Found)
func TestRequestGETSquadNumberNonExistingResponseStatusNotFound(test *testing.T) {
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

// Given GET
// When request path is invalid Squad Number (non-numeric)
// Then response status should be 400 (Bad Request)
func TestRequestGETSquadNumberInvalidParamResponseStatusBadRequest(test *testing.T) {
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

// Given GET
// When request path is existing Squad Number
// Then response status should be 200 (OK)
func TestRequestGETSquadNumberExistingResponseStatusOK(test *testing.T) {
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

// Given GET
// When request path is existing Squad Number
// Then response body should be matching Player
func TestRequestGETSquadNumberExistingResponsePlayer(test *testing.T) {
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

/* GET /players/:id --------------------------------------------------------- */

// Given GET
// When request path is non-existing Id
// Then response status should be 404 (Not Found)
func TestRequestGETIdNonExistingResponseStatusNotFound(test *testing.T) {
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

// Given GET
// When request path is invalid Id (non-numeric)
// Then response status should be 400 (Bad Request)
func TestRequestGETIdInvalidParamResponseStatusBadRequest(test *testing.T) {
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

// Given GET
// When request path is existing Id
// Then response status should be 200 (OK)
func TestRequestGETIdExistingResponseStatusOK(test *testing.T) {
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

// Given GET
// When request path is existing Id
// Then response body should be matching Player
func TestRequestGETIdExistingResponsePlayer(test *testing.T) {
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

// Given GET by ID
// When service.RetrieveByID() fails with non-NotFound error
// Then response status should be 500 (Internal Server Error)
func TestRequestGETByIDServiceGenericErrorResponseStatusInternalServerError(test *testing.T) {
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

/* PUT /players/:id --------------------------------------------------------- */

// Given PUT
// When request body is empty
// Then response status should be 400 (Bad Request)
func TestRequestPUTBodyEmptyResponseStatusBadRequest(test *testing.T) {
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
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, route.PlayersPath+"/"+id, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNotFound, recorder.Code)
}

// Given PUT
// When request path is invalid Id (non-numeric)
// Then response status should be 400 (Bad Request)
func TestRequestPUTIdInvalidParamResponseStatusBadRequest(test *testing.T) {
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
	router := setupRouter(playerController)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, route.PlayersPath+"/"+id, bytes.NewBuffer(body))
	request.Header.Set(ContentType, ApplicationJSON)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(test, http.StatusNoContent, recorder.Code)
}

// Given PUT
// When first service.RetrieveByID() fails with non-NotFound error
// Then response status should be 500 (Internal Server Error)
func TestRequestPUTServiceRetrieveByIDGenericErrorResponseStatusInternalServerError(test *testing.T) {
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

// Given PUT
// When service.Update() fails
// Then response status should be 500 (Internal Server Error)
func TestRequestPUTServiceUpdateFailureResponseStatusInternalServerError(test *testing.T) {
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

// Given DELETE
// When request path is non-existing Id
// Then response status should be 404 (Not Found)
func TestRequestDELETEIdNonExistingIdResponseStatusNotFound(test *testing.T) {
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

// Given DELETE
// When request path is invalid Id (non-numeric)
// Then response status should be 400 (Bad Request)
func TestRequestDELETEIdInvalidParamResponseStatusBadRequest(test *testing.T) {
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

// Given DELETE
// when request path is existing Id
// Then response status should be  204 (No Content)
func TestRequestDELETEIdExistingIdResponseStatusNoContent(test *testing.T) {
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

// Given DELETE
// When first service.RetrieveByID() fails with non-NotFound error
// Then response status should be 500 (Internal Server Error)
func TestRequestDELETEServiceRetrieveByIDGenericErrorResponseStatusInternalServerError(test *testing.T) {
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

// Given DELETE
// When service.Delete() fails
// Then response status should be 500 (Internal Server Error)
func TestRequestDELETEServiceDeleteFailureResponseStatusInternalServerError(test *testing.T) {
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

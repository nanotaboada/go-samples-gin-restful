package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/models"
	"github.com/nanotaboada/go-samples-gin-restful/routes"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMain(main *testing.M) {

	data.Database, data.Error = gorm.Open(sqlite.Open("../data/players-sqlite3.db"), &gorm.Config{})

	if data.Error != nil {
		log.Fatal(data.Error)
	}

	data.Database.AutoMigrate(&models.Player{})

	os.Exit(main.Run())
}

func TestGetPlayers(test *testing.T) {

	// Arrange
	engine := routes.GetEngine()
	request, _ := http.NewRequest("GET", "/players", nil)
	recorder := httptest.NewRecorder()

	// Act
	engine.ServeHTTP(recorder, request)
	var players []models.Player
	json.Unmarshal(recorder.Body.Bytes(), &players)

	// Assert
	assert.NotEmpty(test, players)
	assert.Equal(test, http.StatusOK, recorder.Code)
}

func TestGetPlayerByID(test *testing.T) {

	// Arrange
	engine := routes.GetEngine()
	request, _ := http.NewRequest("GET", "/players/10", nil)
	recorder := httptest.NewRecorder()

	// Act
	engine.ServeHTTP(recorder, request)
	var player models.Player
	json.Unmarshal(recorder.Body.Bytes(), &player)

	// Assert
	assert.NotEmpty(test, player)
	assert.Equal(test, http.StatusOK, recorder.Code)
	assert.Equal(test, 10, player.SquadNumber)
	assert.Equal(test, "Lionel", player.FirstName)
	assert.Equal(test, "Messi", player.LastName)
}

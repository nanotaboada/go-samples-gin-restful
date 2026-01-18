// Package tests provides integration and utility code to support automated
// testing of the application.
package tests

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/nanotaboada/go-samples-gin-restful/model"
)

// MakeExistingPlayer returns a Player that already exists in the original collection.
func MakeExistingPlayer() model.Player {
	return model.Player{
		ID:           1,
		FirstName:    "Damián",
		MiddleName:   "Emiliano",
		LastName:     "Martínez",
		DateOfBirth:  "1992-09-02T00:00:00.000Z",
		SquadNumber:  23,
		Position:     "Goalkeeper",
		AbbrPosition: "GK",
		Team:         "Aston Villa FC",
		League:       "Premier League",
		Starting11:   true,
	}
}

// MakeNonExistingPlayer returns a Player that does not exist in the original collection.
func MakeNonExistingPlayer() model.Player {
	return model.Player{
		ID:           19,
		FirstName:    "Leandro",
		MiddleName:   "Daniel",
		LastName:     "Paredes",
		DateOfBirth:  "1994-06-29T00:00:00.000Z",
		SquadNumber:  5,
		Position:     "Defensive Midfield",
		AbbrPosition: "DM",
		Team:         "AS Roma",
		League:       "Serie A",
		Starting11:   false,
	}
}

// MakePlayersFromJSON loads test players from a JSON file.
func MakePlayersFromJSON() ([]model.Player, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	log.Println("Current directory:", wd)
	// Load test fixture from the tests directory
	name := filepath.Join(wd, "tests", "players.json")
	file, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var players []model.Player
	if err := json.Unmarshal(file, &players); err != nil {
		return nil, err
	}
	return players, nil
}

// Package tests provides integration and utility code to support automated
// testing of the application.
package tests

import "github.com/nanotaboada/go-samples-gin-restful/model"

// Gets a Player that already exists in the original collection
func GetExistingPlayer() model.Player {
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

// Gets a Player that does not exist in the original collection
func GetNonExistingPlayer() model.Player {
	return model.Player{
		ID:           12,
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

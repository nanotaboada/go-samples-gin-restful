// Package tests provides integration and utility code to support automated
// testing of the application.
//
// # UUID strategy for test fixtures
//
// Production code generates player IDs using UUID v4 (random) via
// uuid.NewString() — a different value every time, guaranteed unique.
//
// Test fixtures use UUID v5 (name-based / deterministic) instead:
//
//	uuid5(namespace, name) = SHA-1(namespace + name) → stable 128-bit UUID
//
// Given the same namespace and the same name, uuid5 always returns the same
// UUID. This means fixture IDs are:
//   - Stable across test runs — no random regeneration needed
//   - Derived from meaningful data — the squad number is the name, so the
//     UUID for squad #10 is always the same value
//   - Scoped to this project — the namespace UUID below is project-specific,
//     preventing accidental collisions with other uuid5 users
//
// Namespace: a7d5e3b2-1c4f-5a8d-9e6b-3f2c0d1e4a7b (arbitrary, project-fixed)
// Name:      squad number as a decimal string (e.g. "10" for Messi)
package tests

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/nanotaboada/go-samples-gin-restful/model"
)

// MakeExistingPlayer returns a Player that already exists in the original collection.
// ID is UUID v5 derived from squad number 23 using the project namespace.
func MakeExistingPlayer() model.Player {
	return model.Player{
		ID:           "45ef18c4-a919-50c8-8003-279845045804",
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
// ID is UUID v5 derived from squad number 5 using the project namespace.
func MakeNonExistingPlayer() model.Player {
	return model.Player{
		ID:           "f9897dec-7d3e-568a-9f7f-03d2739c5a7c",
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
	// Important: tests are run from the project root, so the path is relative
	// to the root (not the tests folder).
	name := filepath.Join(wd, "players.json")
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

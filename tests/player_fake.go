// Package tests provides integration and utility code to support automated
// testing of the application.
//
// # Data-state vocabulary for fake factory functions
//
// Factory functions follow a shared three-term vocabulary across all comparison repos:
//
//   - existing    — player is present in the database
//   - nonexistent — player is absent, valid shape for creation (POST scenarios)
//   - unknown     — valid ID format, absent from database (404-by-lookup scenarios)
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
//   - Derived from meaningful data — the player name is the input, so the
//     UUID for Lionel Messi is always the same value
//   - Scoped to this project — the namespace UUID below is the canonical
//     cross-project standard for the 2022 FIFA World Cup Argentina squad
//
// Namespace: f201b13e-c670-473d-885d-e2be219f74c8 (FIFA_WORLD_CUP_QATAR_2022_ARGENTINA_SQUAD)
// Formula:   uuidv5("{firstName}-{lastName}", namespace) — UTF-8
package tests

import (
	"github.com/nanotaboada/go-samples-gin-restful/model"
)

// MakeExistingPlayer returns a Player that already exists in the original collection.
// ID is UUID v5 derived from "Damián-Martínez" using the canonical namespace.
func MakeExistingPlayer() model.Player {
	return model.Player{
		ID:           "01772c59-43f0-5d85-b913-c78e4e281452",
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

// MakeNonexistentPlayer returns a Player that does not exist in the original collection.
// Giovani Lo Celso was selected for the preliminary squad but ruled out by a hamstring
// injury during training camp and replaced by Thiago Almada. Squad 27 is outside the
// seeded range, so POST never conflicts with seeded data and a failed cleanup never
// corrupts seeded records.
// ID is UUID v5 derived from "Giovani-Lo Celso" using the canonical namespace.
func MakeNonexistentPlayer() model.Player {
	return model.Player{
		ID:           "f8d13028-0d22-5513-8774-08a2332b5814",
		FirstName:    "Giovani",
		MiddleName:   "",
		LastName:     "Lo Celso",
		DateOfBirth:  "1996-04-09T00:00:00.000Z",
		SquadNumber:  27,
		Position:     "Central Midfield",
		AbbrPosition: "CM",
		Team:         "Villarreal CF",
		League:       "La Liga",
		Starting11:   false,
	}
}

// MakeUnknownPlayer returns a Player with a valid UUID that does not exist in the
// database. Used for 404-by-lookup scenarios (GET by ID, PUT/DELETE by squad number).
// Squad 99 is chosen to be outside the seeded range and distinct from the nonexistent
// fixture (squad 27), so the two terms remain unambiguous in test output.
func MakeUnknownPlayer() model.Player {
	return model.Player{
		ID:          "00000000-0000-4000-8000-000000000000",
		SquadNumber: 99,
	}
}

// MakeUpdatePlayer returns the update body for Damián Emiliano Martínez (squad 23).
// The canonical Update scenario changes FirstName from "Damián" to "Emiliano" and
// clears MiddleName, reflecting that he goes by his middle name.
func MakeUpdatePlayer() model.Player {
	return model.Player{
		FirstName:    "Emiliano",
		MiddleName:   "",
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

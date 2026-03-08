//go:build ignore

// Seed 001 — Starting Eleven
//
// Drops and recreates storage/players-sqlite3.db, then seeds the 11 players
// in Argentina's starting eleven for the 2022 FIFA World Cup Final.
//
// This file is excluded from normal builds by the //go:build ignore constraint
// above. It is a standalone program intended for developer use only:
//
//	go run ./tools/seed_001_starting_eleven.go
//
// Run seed_002_substitutes.go afterwards to complete the full squad.
// These two seeds map to migration 001 and 002 in the future Goose sequence.
package main

import (
	"log"
	"os"

	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/model"
)

func main() {
	const dbPath = "./storage/players-sqlite3.db"

	// Drop the existing DB so AutoMigrate (called inside data.Connect) starts
	// with a clean slate and applies the current schema from the Player struct.
	if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to remove existing DB: %v", err)
	}
	log.Println("Removed existing DB (or none existed)")

	db := data.Connect(dbPath)
	log.Println("Created DB and applied schema via AutoMigrate")

	players := []model.Player{
		{
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
		},
		{
			ID:           "2b1a505b-8350-510f-8aaf-208a413a2268",
			FirstName:    "Nahuel",
			LastName:     "Molina",
			DateOfBirth:  "1998-04-06T00:00:00.000Z",
			SquadNumber:  26,
			Position:     "Right-Back",
			AbbrPosition: "RB",
			Team:         "Atlético Madrid",
			League:       "La Liga",
			Starting11:   true,
		},
		{
			ID:           "e115fb02-0f87-52f0-91a6-4e5d4ff21a6c",
			FirstName:    "Cristian",
			MiddleName:   "Gabriel",
			LastName:     "Romero",
			DateOfBirth:  "1998-04-27T00:00:00.000Z",
			SquadNumber:  13,
			Position:     "Centre-Back",
			AbbrPosition: "CB",
			Team:         "Tottenham Hotspur",
			League:       "Premier League",
			Starting11:   true,
		},
		{
			ID:           "46d02909-2e4f-52e3-92f0-142d9fbbb280",
			FirstName:    "Nicolás",
			MiddleName:   "Hernán Gonzalo",
			LastName:     "Otamendi",
			DateOfBirth:  "1988-02-12T00:00:00.000Z",
			SquadNumber:  19,
			Position:     "Centre-Back",
			AbbrPosition: "CB",
			Team:         "SL Benfica",
			League:       "Liga Portugal",
			Starting11:   true,
		},
		{
			ID:           "9cb75237-89ca-5722-a9b1-b5aee822fe03",
			FirstName:    "Nicolás",
			MiddleName:   "Alejandro",
			LastName:     "Tagliafico",
			DateOfBirth:  "1992-08-31T00:00:00.000Z",
			SquadNumber:  3,
			Position:     "Left-Back",
			AbbrPosition: "LB",
			Team:         "Olympique Lyon",
			League:       "Ligue 1",
			Starting11:   true,
		},
		{
			ID:           "972d7357-6109-5d7a-b0f2-9f9797610a54",
			FirstName:    "Ángel",
			MiddleName:   "Fabián",
			LastName:     "Di María",
			DateOfBirth:  "1988-02-14T00:00:00.000Z",
			SquadNumber:  11,
			Position:     "Right Winger",
			AbbrPosition: "LW",
			Team:         "SL Benfica",
			League:       "Liga Portugal",
			Starting11:   true,
		},
		{
			ID:           "d5b1810e-a40c-5d4b-9692-c9702a6f4699",
			FirstName:    "Rodrigo",
			MiddleName:   "Javier",
			LastName:     "de Paul",
			DateOfBirth:  "1994-05-24T00:00:00.000Z",
			SquadNumber:  7,
			Position:     "Central Midfield",
			AbbrPosition: "CM",
			Team:         "Atlético Madrid",
			League:       "La Liga",
			Starting11:   true,
		},
		{
			ID:           "7148b61b-a894-5628-bc8d-ca0e55d3a495",
			FirstName:    "Enzo",
			MiddleName:   "Jeremías",
			LastName:     "Fernández",
			DateOfBirth:  "2001-01-17T00:00:00.000Z",
			SquadNumber:  24,
			Position:     "Central Midfield",
			AbbrPosition: "CM",
			Team:         "Chelsea FC",
			League:       "Premier League",
			Starting11:   true,
		},
		{
			ID:           "0b759f4d-3e64-56fd-9cbe-19e3d4c04923",
			FirstName:    "Alexis",
			LastName:     "Mac Allister",
			DateOfBirth:  "1998-12-24T00:00:00.000Z",
			SquadNumber:  20,
			Position:     "Central Midfield",
			AbbrPosition: "CM",
			Team:         "Liverpool FC",
			League:       "Premier League",
			Starting11:   true,
		},
		{
			ID:           "9a5fa2e4-9c9e-58e5-aeb3-8b1b46e87e03",
			FirstName:    "Lionel",
			MiddleName:   "Andrés",
			LastName:     "Messi",
			DateOfBirth:  "1987-06-24T00:00:00.000Z",
			SquadNumber:  10,
			Position:     "Right Winger",
			AbbrPosition: "RW",
			Team:         "Inter Miami CF",
			League:       "Major League Soccer",
			Starting11:   true,
		},
		{
			ID:           "0345510c-bb7a-5ad0-9493-ea4dae3bef49",
			FirstName:    "Julián",
			LastName:     "Álvarez",
			DateOfBirth:  "2000-01-31T00:00:00.000Z",
			SquadNumber:  9,
			Position:     "Centre-Forward",
			AbbrPosition: "CF",
			Team:         "Manchester City",
			League:       "Premier League",
			Starting11:   true,
		},
	}

	if err := db.Create(&players).Error; err != nil {
		log.Fatalf("failed to seed starting eleven: %v", err)
	}
	log.Printf("Seeded %d players (starting eleven) into %s", len(players), dbPath)
}

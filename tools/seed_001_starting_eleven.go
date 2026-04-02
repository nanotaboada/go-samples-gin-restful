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
		},
		{
			ID:           "da31293b-4c7e-5e0f-a168-469ee29ecbc4",
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
			ID:           "c096c69e-762b-5281-9290-bb9c167a24a0",
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
			ID:           "d5f7dd7a-1dcb-5960-ba27-e34865b63358",
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
			ID:           "2f6f90a0-9b9d-5023-96d2-a2aaf03143a6",
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
			ID:           "b5b46e79-929e-5ed2-949d-0d167109c022",
			FirstName:    "Ángel",
			MiddleName:   "Fabián",
			LastName:     "Di María",
			DateOfBirth:  "1988-02-14T00:00:00.000Z",
			SquadNumber:  11,
			Position:     "Right Winger",
			AbbrPosition: "RW",
			Team:         "SL Benfica",
			League:       "Liga Portugal",
			Starting11:   true,
		},
		{
			ID:           "0293b282-1da8-562e-998e-83849b417a42",
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
			ID:           "d3ba552a-dac3-588a-b961-1ea7224017fd",
			FirstName:    "Enzo",
			MiddleName:   "Jeremías",
			LastName:     "Fernández",
			DateOfBirth:  "2001-01-17T00:00:00.000Z",
			SquadNumber:  24,
			Position:     "Central Midfield",
			AbbrPosition: "CM",
			Team:         "SL Benfica",
			League:       "Liga Portugal",
			Starting11:   true,
		},
		{
			ID:           "9613cae9-16ab-5b54-937e-3135123b9e0d",
			FirstName:    "Alexis",
			LastName:     "Mac Allister",
			DateOfBirth:  "1998-12-24T00:00:00.000Z",
			SquadNumber:  20,
			Position:     "Central Midfield",
			AbbrPosition: "CM",
			Team:         "Brighton & Hove Albion",
			League:       "Premier League",
			Starting11:   true,
		},
		{
			ID:           "acc433bf-d505-51fe-831e-45eb44c4d43c",
			FirstName:    "Lionel",
			MiddleName:   "Andrés",
			LastName:     "Messi",
			DateOfBirth:  "1987-06-24T00:00:00.000Z",
			SquadNumber:  10,
			Position:     "Right Winger",
			AbbrPosition: "RW",
			Team:         "Paris Saint-Germain",
			League:       "Ligue 1",
			Starting11:   true,
		},
		{
			ID:           "38bae91d-8519-55a2-b30a-b9fe38849bfb",
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

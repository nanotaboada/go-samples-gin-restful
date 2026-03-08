// Seed 002 — Substitutes
//
// Seeds the 14 substitute players from Argentina's 2022 FIFA World Cup squad.
// Requires seed_001_starting_eleven.go to have been run first (it creates the
// DB and the schema).
//
// This file is excluded from normal builds by the //go:build ignore constraint
// below. It is a standalone program intended for developer use only:
//
//	go run ./tools/seed_002_substitutes.go
//
// These two seeds map to migration 001 and 002 in the future Goose sequence.

//go:build ignore

package main

import (
	"log"

	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/model"
)

func main() {
	const dbPath = "./storage/players-sqlite3.db"

	// data.Connect calls AutoMigrate, which is idempotent — safe to call even
	// if the schema was already created by seed_001.
	db := data.Connect(dbPath)

	players := []model.Player{
		{
			ID:           "f8e66ec2-008f-50ee-937f-8f4180165216",
			FirstName:    "Franco",
			MiddleName:   "Daniel",
			LastName:     "Armani",
			DateOfBirth:  "1986-10-16T00:00:00.000Z",
			SquadNumber:  1,
			Position:     "Goalkeeper",
			AbbrPosition: "GK",
			Team:         "River Plate",
			League:       "Copa de la Liga",
			Starting11:   false,
		},
		{
			ID:           "5d1cf792-9fb8-57b1-b8af-ecb6be526315",
			FirstName:    "Juan",
			MiddleName:   "Marcos",
			LastName:     "Foyth",
			DateOfBirth:  "1998-01-12T00:00:00.000Z",
			SquadNumber:  2,
			Position:     "Right-Back",
			AbbrPosition: "RB",
			Team:         "Villareal",
			League:       "La Liga",
			Starting11:   false,
		},
		{
			ID:           "4507092d-3454-5664-922f-dc718ed10025",
			FirstName:    "Gonzalo",
			MiddleName:   "Ariel",
			LastName:     "Montiel",
			DateOfBirth:  "1997-01-01T00:00:00.000Z",
			SquadNumber:  4,
			Position:     "Right-Back",
			AbbrPosition: "RB",
			Team:         "Nottingham Forrest",
			League:       "Premier League",
			Starting11:   false,
		},
		{
			ID:           "5f3099ad-ba4d-5b52-bb80-c30383061231",
			FirstName:    "Germán",
			MiddleName:   "Alejo",
			LastName:     "Pezzella",
			DateOfBirth:  "1991-06-27T00:00:00.000Z",
			SquadNumber:  6,
			Position:     "Centre-Back",
			AbbrPosition: "CB",
			Team:         "Real Betis Balompié",
			League:       "La Liga",
			Starting11:   false,
		},
		{
			ID:           "813181ca-1c37-5fe5-b0f5-0a99cb7445a2",
			FirstName:    "Marcos",
			MiddleName:   "Javier",
			LastName:     "Acuña",
			DateOfBirth:  "1991-10-28T00:00:00.000Z",
			SquadNumber:  8,
			Position:     "Left-Back",
			AbbrPosition: "LB",
			Team:         "Sevilla FC",
			League:       "La Liga",
			Starting11:   false,
		},
		{
			ID:           "d273b527-ff15-5133-9f4e-81cd2531503c",
			FirstName:    "Gerónimo",
			LastName:     "Rulli",
			DateOfBirth:  "1992-05-20T00:00:00.000Z",
			SquadNumber:  12,
			Position:     "Goalkeeper",
			AbbrPosition: "GK",
			Team:         "Ajax Amsterdam",
			League:       "Eredivisie",
			Starting11:   false,
		},
		{
			ID:           "ca9bca3f-31a3-58af-be6e-9ea11b0e0c0a",
			FirstName:    "Exequiel",
			MiddleName:   "Alejandro",
			LastName:     "Palacios",
			DateOfBirth:  "1998-10-05T00:00:00.000Z",
			SquadNumber:  14,
			Position:     "Central Midfield",
			AbbrPosition: "CM",
			Team:         "Bayer 04 Leverkusen",
			League:       "Bundesliga",
			Starting11:   false,
		},
		{
			ID:           "a1f44b3b-64de-5796-9616-be11a1fa1ac8",
			FirstName:    "Ángel",
			MiddleName:   "Martín",
			LastName:     "Correa",
			DateOfBirth:  "1995-03-09T00:00:00.000Z",
			SquadNumber:  15,
			Position:     "Right Winger",
			AbbrPosition: "RW",
			Team:         "Atlético Madrid",
			League:       "La Liga",
			Starting11:   false,
		},
		{
			ID:           "3e7bc8e0-0523-5495-9f77-2455422dadb5",
			FirstName:    "Thiago",
			MiddleName:   "Ezequiel",
			LastName:     "Almada",
			DateOfBirth:  "2001-04-26T00:00:00.000Z",
			SquadNumber:  16,
			Position:     "Attacking Midfield",
			AbbrPosition: "AM",
			Team:         "Atlanta United FC",
			League:       "Major League Soccer",
			Starting11:   false,
		},
		{
			ID:           "dc503d71-b6d7-51e2-90d0-932075a74048",
			FirstName:    "Alejandro",
			MiddleName:   "Darío",
			LastName:     "Gómez",
			DateOfBirth:  "1988-02-15T00:00:00.000Z",
			SquadNumber:  17,
			Position:     "Left Winger",
			AbbrPosition: "LW",
			Team:         "AC Monza",
			League:       "Serie A",
			Starting11:   false,
		},
		{
			ID:           "0091ab48-00bb-5744-89f3-200d8582a28a",
			FirstName:    "Guido",
			LastName:     "Rodríguez",
			DateOfBirth:  "1994-04-12T00:00:00.000Z",
			SquadNumber:  18,
			Position:     "Defensive Midfield",
			AbbrPosition: "DM",
			Team:         "Real Betis Balompié",
			League:       "La Liga",
			Starting11:   false,
		},
		{
			ID:           "d13a8488-6373-5265-9c16-98a32c462562",
			FirstName:    "Paulo",
			MiddleName:   "Exequiel",
			LastName:     "Dybala",
			DateOfBirth:  "1993-11-15T00:00:00.000Z",
			SquadNumber:  21,
			Position:     "Second Striker",
			AbbrPosition: "SS",
			Team:         "AS Roma",
			League:       "Serie A",
			Starting11:   false,
		},
		{
			ID:           "551991be-7313-505c-be42-7cfa63ec214c",
			FirstName:    "Lautaro",
			MiddleName:   "Javier",
			LastName:     "Martínez",
			DateOfBirth:  "1997-08-22T00:00:00.000Z",
			SquadNumber:  22,
			Position:     "Centre-Forward",
			AbbrPosition: "CF",
			Team:         "Inter Milan",
			League:       "Serie A",
			Starting11:   false,
		},
		{
			ID:           "76db51c9-4d21-50d9-9da0-1939ff821ac1",
			FirstName:    "Lisandro",
			LastName:     "Martínez",
			DateOfBirth:  "1998-01-18T00:00:00.000Z",
			SquadNumber:  25,
			Position:     "Centre-Back",
			AbbrPosition: "CB",
			Team:         "Manchester United",
			League:       "Premier League",
			Starting11:   false,
		},
	}

	if err := db.Create(&players).Error; err != nil {
		log.Fatalf("failed to seed substitutes: %v", err)
	}
	log.Printf("Seeded %d players (substitutes) into %s", len(players), dbPath)
}

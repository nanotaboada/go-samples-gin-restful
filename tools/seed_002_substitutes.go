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
			ID:           "5a9cd988-95e6-54c1-bc34-9aa08acca8d0",
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
			ID:           "5fdb10e8-38c0-5084-9a3f-b369a960b9c2",
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
			ID:           "bbd441f7-fcfb-5834-8468-2a9004b64c8c",
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
			ID:           "d8bfea25-f189-5d5e-b3a5-ed89329b9f7c",
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
			ID:           "dca343a8-12e5-53d6-89a8-916b120a5ee4",
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
			ID:           "c62f2ac1-41e8-5d34-b073-2ba0913d0e31",
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
			ID:           "d3b0e8e8-2c34-531a-b608-b24fed0ef986",
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
			ID:           "b1306b7b-a3a4-5f7c-90fd-dd5bdbed57ba",
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
			ID:           "ecec27e8-487b-5622-b116-0855020477ed",
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
			ID:           "7cc8d527-56a2-58bd-9528-2618fc139d30",
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
			ID:           "191c82af-0c51-526a-b903-c3600b61b506",
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
			ID:           "7941cd7c-4df1-5952-97e8-1e7f5d08e8aa",
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
			ID:           "79c96f29-c59f-5f98-96b8-3a5946246624",
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
			ID:           "98306555-a466-5d18-804e-dc82175e697b",
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

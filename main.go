package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
--------------------------------------------------------------------------------
Model
--------------------------------------------------------------------------------
*/

// https://go.dev/tour/basics/11
type player struct {
	ID           string `json:"id"`
	FirstName    string `json:"firstName"`
	MiddleName   string `json:"middleName"`
	LastName     string `json:"lastName"`
	DateOfBirth  string `json:"dateOfBirth"`
	SquadNumber  int    `json:"squadNumber"`
	Position     string `json:"position"`
	AbbrPosition string `json:"abbrPosition"`
	Team         string `json:"team"`
	League       string `json:"league"`
	Starting11   bool   `json:"starting11"`
}

/*
--------------------------------------------------------------------------------
Data
--------------------------------------------------------------------------------
*/

var players = []player{
	{
		ID:           "1",
		FirstName:    "Damián",
		MiddleName:   "Emiliano",
		LastName:     "Martínez",
		DateOfBirth:  "1992-09-01T21:00:00-00:00",
		SquadNumber:  23,
		Position:     "Goalkeeper",
		AbbrPosition: "GK",
		Team:         "Aston Villa FC",
		League:       "Premier League",
		Starting11:   true,
	},
	{
		ID:           "2",
		FirstName:    "Nahuel",
		MiddleName:   "",
		LastName:     "Molina",
		DateOfBirth:  "1998-04-05T21:00:00-00:00",
		SquadNumber:  26,
		Position:     "Right-Back",
		AbbrPosition: "RB",
		Team:         "Atlético Madrid",
		League:       "La Liga",
		Starting11:   true,
	},
	{
		ID:           "3",
		FirstName:    "Cristian",
		MiddleName:   "Gabriel",
		LastName:     "Romero",
		DateOfBirth:  "1998-04-26T21:00:00-00:00",
		SquadNumber:  13,
		Position:     "Centre-Back",
		AbbrPosition: "CB",
		Team:         "Tottenham Hotspur",
		League:       "Premier League",
		Starting11:   true,
	},
	{
		ID:           "4",
		FirstName:    "Nicolás",
		MiddleName:   "Hernán Gonzalo",
		LastName:     "Otamendi",
		DateOfBirth:  "1988-02-11T21:00:00-00:00",
		SquadNumber:  19,
		Position:     "Centre-Back",
		AbbrPosition: "CB",
		Team:         "SL Benfica",
		League:       "Liga Portugal",
		Starting11:   true,
	},
	{
		ID:           "5",
		FirstName:    "Nicolás",
		MiddleName:   "Alejandro",
		LastName:     "Tagliafico",
		DateOfBirth:  "1992-08-30T21:00:00-00:00",
		SquadNumber:  3,
		Position:     "Left-Back",
		AbbrPosition: "LB",
		Team:         "Olympique Lyon",
		League:       "Ligue 1",
		Starting11:   true,
	},
	{
		ID:           "6",
		FirstName:    "Ángel",
		MiddleName:   "Fabián",
		LastName:     "Di María",
		DateOfBirth:  "1988-02-13T21:00:00-00:00",
		SquadNumber:  11,
		Position:     "Right Winger",
		AbbrPosition: "LW",
		Team:         "SL Benfica",
		League:       "Liga Portugal",
		Starting11:   true,
	},
	{
		ID:           "7",
		FirstName:    "Rodrigo",
		MiddleName:   "Javier",
		LastName:     "de Paul",
		DateOfBirth:  "1994-05-23T21:00:00-00:00",
		SquadNumber:  7,
		Position:     "Central Midfield",
		AbbrPosition: "CM",
		Team:         "Atlético Madrid",
		League:       "La Liga",
		Starting11:   true,
	},
	{
		ID:           "8",
		FirstName:    "Enzo",
		MiddleName:   "Jeremías",
		LastName:     "Fernández",
		DateOfBirth:  "2001-01-16T21:00:00-00:00",
		SquadNumber:  24,
		Position:     "Central Midfield",
		AbbrPosition: "CM",
		Team:         "Chelsea FC",
		League:       "Premier League",
		Starting11:   true,
	},
	{
		ID:           "9",
		FirstName:    "Alexis",
		MiddleName:   "",
		LastName:     "Mac Allister",
		DateOfBirth:  "1998-12-23T21:00:00-00:00",
		SquadNumber:  20,
		Position:     "Central Midfield",
		AbbrPosition: "CM",
		Team:         "Liverpool FC",
		League:       "Premier League",
		Starting11:   true,
	},
	{
		ID:           "10",
		FirstName:    "Lionel",
		MiddleName:   "Andrés",
		LastName:     "Messi",
		DateOfBirth:  "1987-06-23T21:00:00-00:00",
		SquadNumber:  10,
		Position:     "Right Winger",
		AbbrPosition: "RW",
		Team:         "Inter Miami CF",
		League:       "Major League Soccer",
		Starting11:   true,
	},
	{
		ID:           "11",
		FirstName:    "Julián",
		MiddleName:   "",
		LastName:     "Álvarez",
		DateOfBirth:  "2000-01-30T21:00:00-00:00",
		SquadNumber:  9,
		Position:     "Centre-Forward",
		AbbrPosition: "CF",
		Team:         "Manchester City",
		League:       "Premier League",
		Starting11:   true,
	},
}

/*
--------------------------------------------------------------------------------
Routes
--------------------------------------------------------------------------------
*/

func main() {

	router := gin.Default()
	router.GET("/players", getPlayers)
	router.GET("/player/:id", getPlayerByID)

	router.Run("localhost:8080")
}

/*
--------------------------------------------------------------------------------
Handlers
--------------------------------------------------------------------------------
*/

func getPlayers(context *gin.Context) {

	context.IndentedJSON(http.StatusOK, players)
}

func getPlayerByID(context *gin.Context) {

	id := context.Param("id")

	for _, player := range players {
		if player.ID == id {
			context.IndentedJSON(http.StatusOK, player)
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Not Found."})
}

package models

// https://go.dev/tour/basics/11
// https://go.dev/ref/spec#Exported_identifiers
type Player struct {
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

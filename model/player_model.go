// Package model defines the data structures used throughout the application,
// including Player.
package model

// Player is a footballer, a sportsperson who plays football
type Player struct {
	ID           int    `json:"id" gorm:"primaryKey"`           // The unique identifier for the Player
	FirstName    string `json:"firstName" gorm:"column:firstName"` // The first name of the Player
	MiddleName   string `json:"middleName" gorm:"column:middleName"` // The middle name of the Player, if any
	LastName     string `json:"lastName" gorm:"column:lastName"` // The last name of the Player
	DateOfBirth  string `json:"dateOfBirth" gorm:"column:dateOfBirth"` // The date of birth of the Player
	SquadNumber  int    `json:"squadNumber" gorm:"column:squadNumber"` // The squad number assigned to the Player
	Position     string `json:"position" gorm:"column:position"` // The playing position of the Player
	AbbrPosition string `json:"abbrPosition" gorm:"column:abbrPosition"` // The abbreviated form of the Player's position
	Team         string `json:"team" gorm:"column:team"` // The team to which the Player belongs
	League       string `json:"league" gorm:"column:league"` // The league where the team plays
	Starting11   bool   `json:"starting11" gorm:"column:starting11"` // Indicates whether the Player is in the starting 11
}

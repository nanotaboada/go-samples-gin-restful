package models

// https://go.dev/tour/basics/11
// https://go.dev/ref/spec#Exported_identifiers
type Player struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	FirstName    string `json:"firstName" gorm:"column:firstName"`
	MiddleName   string `json:"middleName" gorm:"column:middleName"`
	LastName     string `json:"lastName" gorm:"column:lastName"`
	DateOfBirth  string `json:"dateOfBirth" gorm:"column:dateOfBirth"`
	SquadNumber  int    `json:"squadNumber" gorm:"column:squadNumber"`
	Position     string `json:"position" gorm:"column:position"`
	AbbrPosition string `json:"abbrPosition" gorm:"column:abbrPosition"`
	Team         string `json:"team" gorm:"column:team"`
	League       string `json:"league" gorm:"column:league"`
	Starting11   bool   `json:"starting11" gorm:"column:starting11"`
}

// Package model defines the data structures used throughout the application,
// including Player.
package model

// Player is a footballer, a sportsperson who plays football.
//
// # Struct tags
//
// Each field carries two sets of struct tags that Go reads at runtime via
// reflection:
//
//   - `json:"..."` — controls marshalling to/from JSON.
//     Field names use camelCase to follow JSON conventions, matching what
//     the API clients (JavaScript, Python, etc.) expect.
//
//   - `gorm:"..."` — controls how GORM maps the struct to the SQLite table.
//     `column:` sets the exact column name in the DB.
//     `primaryKey` marks the primary key (GORM won't auto-increment it because
//     the type is string — UUIDs are assigned by the application, not the DB).
//     `uniqueIndex` creates a unique index in SQLite, enforced at the DB level.
//
// # ID design
//
// ID is a string (not an integer auto-increment) because it stores a UUID v4,
// generated server-side on POST. This keeps the internal key opaque and stable
// across environments.  Clients use squadNumber to identify players in PUT and
// DELETE requests. The UUID is exposed via GET /players/uuid/:id.
type Player struct {
	ID           string `json:"id" gorm:"column:id;primaryKey"`                    // Internal UUID (server-generated, opaque to clients)
	FirstName    string `json:"firstName" gorm:"column:firstName"`                 // The first name of the Player
	MiddleName   string `json:"middleName" gorm:"column:middleName"`               // The middle name of the Player, if any
	LastName     string `json:"lastName" gorm:"column:lastName"`                   // The last name of the Player
	DateOfBirth  string `json:"dateOfBirth" gorm:"column:dateOfBirth"`             // The date of birth of the Player
	SquadNumber  int    `json:"squadNumber" gorm:"column:squadNumber;uniqueIndex"` // User-facing unique identifier; DB-enforced uniqueness
	Position     string `json:"position" gorm:"column:position"`                   // The playing position of the Player
	AbbrPosition string `json:"abbrPosition" gorm:"column:abbrPosition"`           // The abbreviated form of the Player's position
	Team         string `json:"team" gorm:"column:team"`                           // The team to which the Player belongs
	League       string `json:"league" gorm:"column:league"`                       // The league where the team plays
	Starting11   bool   `json:"starting11" gorm:"column:starting11"`               // Indicates whether the Player is in the starting 11
}

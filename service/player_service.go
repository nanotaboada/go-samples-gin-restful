// Package service contains business logic for Player operations, primarily
// interacting with the ORM.
package service

import (
	"github.com/nanotaboada/go-samples-gin-restful/model"
	"gorm.io/gorm"
)

// PlayerService defines the contract for player business logic.
//
// In Go, interfaces are satisfied implicitly: any type that implements all of
// these methods automatically satisfies PlayerService — no "implements"
// keyword is needed. This makes it easy to swap the real implementation for a
// mock in tests without modifying any production code.
type PlayerService interface {
	Create(player *model.Player) error
	RetrieveAll() ([]model.Player, error)
	RetrieveByID(id string) (model.Player, error)
	RetrieveBySquadNumber(squadNumber int) (model.Player, error)
	Update(player *model.Player) error
	Delete(player *model.Player) error
}

// playerService implements PlayerService using GORM.
// It is unexported (lowercase) intentionally: callers interact only through
// the PlayerService interface, never with the concrete struct directly.
type playerService struct {
	db *gorm.DB // The GORM database handle (connection pool, safe for concurrent use)
}

// NewPlayerService returns a PlayerService backed by the given *gorm.DB.
// Returning the interface type (not *playerService) keeps the concrete type
// hidden from callers and allows the mock to substitute it transparently.
func NewPlayerService(db *gorm.DB) PlayerService {
	return &playerService{db: db}
}

// Create inserts a new Player row into the database.
// GORM uses the struct's field values and tags to build the INSERT statement.
// https://gorm.io/docs/create.html
func (s *playerService) Create(player *model.Player) error {
	return s.db.Create(player).Error
}

// RetrieveAll fetches every row from the players table.
// Find populates the slice and never returns gorm.ErrRecordNotFound (it
// returns an empty slice instead), so callers don't need to check for that
// specific error here.
// https://gorm.io/docs/query.html
func (s *playerService) RetrieveAll() ([]model.Player, error) {
	var players []model.Player
	result := s.db.Find(&players)
	return players, result.Error
}

// RetrieveByID fetches a single Player by its internal UUID.
// First adds "LIMIT 1" and returns gorm.ErrRecordNotFound when no row matches,
// which the controller translates into a 404 response.
// https://gorm.io/docs/query.html
func (s *playerService) RetrieveByID(id string) (model.Player, error) {
	var player model.Player
	result := s.db.Where("id = ?", id).First(&player)
	return player, result.Error
}

// RetrieveBySquadNumber fetches a single Player by squad number.
// Like RetrieveByID, First returns gorm.ErrRecordNotFound on miss; the
// controller uses errors.Is to distinguish "not found" from other DB errors.
// https://gorm.io/docs/query.html
func (s *playerService) RetrieveBySquadNumber(squadNumber int) (model.Player, error) {
	var player model.Player
	result := s.db.Where("squadNumber = ?", squadNumber).First(&player)
	return player, result.Error
}

// Update replaces a Player record entirely (full update / HTTP PUT semantics).
// Save issues an UPDATE covering all columns, not just the changed ones.
// Using Save instead of Updates avoids accidentally zeroing fields that the
// caller omitted — the caller must always pass the complete player struct.
// https://gorm.io/docs/update.html
func (s *playerService) Update(player *model.Player) error {
	return s.db.Save(player).Error
}

// Delete removes a Player from the database permanently.
// Because the Player struct has no gorm.DeletedAt (soft-delete) field, GORM
// issues a hard DELETE statement rather than setting a deleted_at timestamp.
// https://gorm.io/docs/delete.html
func (s *playerService) Delete(player *model.Player) error {
	return s.db.Delete(player).Error
}

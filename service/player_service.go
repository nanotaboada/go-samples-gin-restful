// Package service contains business logic for Player operations, primarily
// interacting with the ORM.
package service

import (
	"github.com/nanotaboada/go-samples-gin-restful/model"
	"gorm.io/gorm"
)

// PlayerService defines the contract for player business logic
type PlayerService interface {
	Create(player *model.Player) error
	RetrieveAll() ([]model.Player, error)
	RetrieveByID(id int) (model.Player, error)
	RetrieveBySquadNumber(squadNumber int) (model.Player, error)
	Update(player *model.Player) error
	Delete(id int) error
}

// playerService implements PlayerService using GORM
type playerService struct {
	db *gorm.DB
}

// NewPlayerService creates a new PlayerService with the given database
func NewPlayerService(db *gorm.DB) PlayerService {
	return &playerService{db: db}
}

// Create adds a new Player in the database
func (s *playerService) Create(player *model.Player) error {
	// https://gorm.io/docs/create.html
	return s.db.Create(player).Error
}

// RetrieveAll gets all players from the database
func (s *playerService) RetrieveAll() ([]model.Player, error) {
	var players []model.Player
	// https://gorm.io/docs/query.html
	result := s.db.Find(&players)
	return players, result.Error
}

// RetrieveByID gets a Player by ID from the database
func (s *playerService) RetrieveByID(id int) (model.Player, error) {
	var player model.Player
	// https://gorm.io/docs/query.html
	result := s.db.First(&player, id)
	return player, result.Error
}

// RetrieveBySquadNumber gets a Player by its Squad Number from the database
func (s *playerService) RetrieveBySquadNumber(squadNumber int) (model.Player, error) {
	var player model.Player
	// https://gorm.io/docs/query.html
	result := s.db.Where("squadNumber = ?", squadNumber).First(&player)
	return player, result.Error
}

// Update replaces (completely) a Player in the database
func (s *playerService) Update(player *model.Player) error {
	// https://gorm.io/docs/update.html
	return s.db.Save(player).Error
}

// Delete removes a Player by its ID from the database
func (s *playerService) Delete(id int) error {
	// https://gorm.io/docs/delete.html
	return s.db.Delete(&model.Player{}, id).Error
}

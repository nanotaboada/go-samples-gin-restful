// Package service contains business logic for Player operations, primarily
// interacting with the ORM.
package service

import (
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/model"
)

// Create adds a new Player in the database
func Create(player *model.Player) error {
	// https://gorm.io/docs/create.html
	return data.DB.Create(player).Error
}

// RetrieveAll gets all players from the database
func RetrieveAll() ([]model.Player, error) {
	var players []model.Player
	// https://gorm.io/docs/query.html
	result := data.DB.Find(&players)
	return players, result.Error
}

// RetrieveByID gets a Player by ID from the database
func RetrieveByID(id int) (model.Player, error) {
	var player model.Player
	// https://gorm.io/docs/query.html
	result := data.DB.First(&player, id)
	return player, result.Error
}

// RetrieveBySquadNumber gets a Player by its Squad Number from the database
func RetrieveBySquadNumber(squadNumber int) (model.Player, error) {
	var player model.Player
	// https://gorm.io/docs/query.html
	result := data.DB.Where("squadNumber = ?", squadNumber).First(&player)
	return player, result.Error
}

// Update replaces (completely) a Player in the database
func Update(player *model.Player) error {
	// https://gorm.io/docs/update.html
	return data.DB.Save(player).Error
}

// Delete removes a Player by its ID from the database
func Delete(id int) error {
	// https://gorm.io/docs/delete.html
	return data.DB.Delete(&model.Player{}, id).Error
}

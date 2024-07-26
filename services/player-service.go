package services

import (
	"github.com/nanotaboada/go-samples-gin-restful/data"
	"github.com/nanotaboada/go-samples-gin-restful/models"
)

// Create adds a new player in the database
func Create(player *models.Player) error {
	// https://gorm.io/docs/create.html
	return data.DB.Create(player).Error
}

// RetrieveAll gets all players from the database
func RetrieveAll() ([]models.Player, error) {
	var players []models.Player
	// https://gorm.io/docs/query.html
	result := data.DB.Find(&players)
	return players, result.Error
}

// RetrieveByID gets a player by ID from the database
func RetrieveByID(id int) (models.Player, error) {
	var player models.Player
	// https://gorm.io/docs/query.html
	result := data.DB.First(&player, id)
	return player, result.Error
}

// Update replaces (completely) an existing player in the database
func Update(player *models.Player) error {
	// https://gorm.io/docs/update.html
	return data.DB.Save(player).Error
}

// Delete removes a player by its ID from the database
func Delete(id int) error {
	// https://gorm.io/docs/delete.html
	return data.DB.Delete(&models.Player{}, id).Error
}

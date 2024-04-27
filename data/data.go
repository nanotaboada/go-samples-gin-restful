/* -----------------------------------------------------------------------------
 * Data
 * -------------------------------------------------------------------------- */

package data

import (
	"log"

	"github.com/nanotaboada/go-samples-gin-restful/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

// Connect initializes and returns a global DB connection
func Connect(path string) {
	// https://gorm.io/docs/connecting_to_the_database.html
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = db

	DB.AutoMigrate(&models.Player{})
}

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

func Connect(path string) {

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = db

	DB.AutoMigrate(&models.Player{})
}

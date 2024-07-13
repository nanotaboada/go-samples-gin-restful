/* -----------------------------------------------------------------------------
 * Data
 * -------------------------------------------------------------------------- */

package data

import (
	"log"
	"time"

	"github.com/nanotaboada/go-samples-gin-restful/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

// Connect initializes and returns a global DB connection
func Connect(path string) {
	// https://gorm.io/docs/logger.html
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	// https://gorm.io/docs/connecting_to_the_database.html
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatal(err)
	}

	DB = db

	DB.AutoMigrate(&models.Player{})
}

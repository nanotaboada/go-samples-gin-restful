// Package data provides database connectivity and data access functions for
// Player entities.
package data

import (
	"log"
	"time"

	"github.com/nanotaboada/go-samples-gin-restful/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect initializes and returns a DB connection
func Connect(dataSourceName string) *gorm.DB {
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
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&model.Player{}); err != nil {
		log.Fatal(err)
	}

	return db
}

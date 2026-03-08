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

// Connect initializes and returns a GORM database connection backed by SQLite.
//
// dataSourceName is a SQLite DSN (Data Source Name).  Two forms are used in
// this project:
//
//   - File-based (production): a path such as "./storage/players-sqlite3.db"
//   - In-memory (tests): "file::memory:?cache=shared"
//     The "?cache=shared" query param is required so that all connections in
//     the same process share the same in-memory database; without it each
//     call to gorm.Open would get an empty, isolated database.
//
// AutoMigrate compares the current SQLite schema with the Player struct and
// applies the minimum set of DDL changes (CREATE TABLE if absent, ADD COLUMN
// for new fields, CREATE INDEX for new uniqueIndex tags).  It never drops
// columns or indexes, so it is safe to call on an existing database.
func Connect(dataSourceName string) *gorm.DB {
	// GORM's built-in logger prints slow queries and all SQL statements.
	// SlowThreshold defines when a query is considered "slow" and logged at
	// WARN level; queries above this threshold are highlighted in the output.
	// https://gorm.io/docs/logger.html
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	// gorm.Open returns a *gorm.DB — a connection-pool handle, not a single
	// connection.  The sqlite driver keeps the underlying file/memory handle
	// open for the lifetime of the process.
	// https://gorm.io/docs/connecting_to_the_database.html
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatal(err)
	}

	// AutoMigrate creates or updates the "players" table to match the Player
	// struct.  GORM derives the table name from the struct name by
	// pluralising it ("Player" → "players").
	if err := db.AutoMigrate(&model.Player{}); err != nil {
		log.Fatal(err)
	}

	return db
}

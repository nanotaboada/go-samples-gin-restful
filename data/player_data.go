// Package data provides database connectivity and data access functions for
// Player entities.
package data

import (
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/nanotaboada/go-samples-gin-restful/migrations"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect initializes and returns a GORM database connection backed by SQLite,
// then applies all pending versioned migrations via goose.
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
// Schema and seed migrations live in the /migrations directory and are
// embedded into the binary at compile time.  goose tracks applied migrations
// in a goose_db_version table and is idempotent: already-applied migrations
// are skipped on subsequent startups.
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

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	// SQLite does not support concurrent writes; a single open connection
	// prevents "database is locked" errors under concurrent request load.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(sqlDB, "."); err != nil {
		log.Fatal(err)
	}

	return db
}

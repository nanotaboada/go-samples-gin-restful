package data

import (
	"log"

	"github.com/nanotaboada/go-samples-gin-restful/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Database *gorm.DB
	Error    error
)

func ConnectToSqlite() {

	Database, Error = gorm.Open(sqlite.Open("./data/players-sqlite3.db"), &gorm.Config{})

	if Error != nil {
		log.Fatal(Error)
	}

	Database.AutoMigrate(&models.Player{})
}

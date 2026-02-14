package database

import (
	"url-shortener/src/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB opens an SQLite connection and auto-migrates the schema.
func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.URL{}, &models.Metrics{}); err != nil {
		return nil, err
	}

	return db, nil
}

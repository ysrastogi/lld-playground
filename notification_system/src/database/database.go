package database

import (
	"notification_system/src/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1)

	if err := db.AutoMigrate(
		&models.Notification{},
		&models.DeliveryAttempt{},
		&models.UserPreference{},
	); err != nil {
		return nil, err
	}
	return db, nil
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type Metrics struct {
	gorm.Model
	ShortCode      string `gorm:"index;not null"`
	AccessCount    int    `gorm:"default:0"`
	LastAccessedAt time.Time
}

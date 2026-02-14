package models

import (
	"gorm.io/gorm"
	"time"
)

type URL struct {
	gorm.Model
	ShortCode string `gorm:"uniqueIndex;not null"`
	Url  string `gorm:"not null"`
	IsStale  bool   `gorm:"default:false"`
	Status   Status `gorm:"default:ACTIVE"`
	ExpiresAt time.Time `gorm:"default:null"`
}

type Status string
const (
	ACTIVE Status = "ACTIVE"
	INACTIVE Status = "INACTIVE"
	STALE Status = "STALE"
)
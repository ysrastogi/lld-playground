package models

import "time"

type NotificationCategory string

const (
	CategoryTransaction NotificationCategory = "transaction"
	CategoryMarketing   NotificationCategory = "marketing"
	CategorySecurity    NotificationCategory = "security"
	CategorySystem      NotificationCategory = "system"
)

type NotificationStatus string

const (
	StatusCreated    NotificationStatus = "created"
	StatusQueued     NotificationStatus = "queued"
	StatusProcessing NotificationStatus = "processing"
	StatusSent       NotificationStatus = "sent"
	StatusFailed     NotificationStatus = "failed"
	StatusDelivered  NotificationStatus = "delivered"
	StatusRead       NotificationStatus = "read"
)

type Notification struct {
	ID        uint                 `gorm:"primaryKey;autoIncrement"`
	UserID    string               `gorm:"index;not null"`
	Category  NotificationCategory `gorm:"not null"`
	Title     string               `gorm:"not null;size:255"`
	Content   string               `gorm:"not null"`
	Status    NotificationStatus   `gorm:"not null"`
	Priority  int                  `gorm:"not null;default:1"`
	CreatedAt time.Time            `gorm:"autoCreateTime"`
	UpdatedAt time.Time            `gorm:"autoUpdateTime"`
}

func (n *Notification) UpdateStatus(newStatus NotificationStatus) {
	n.Status = newStatus
	n.UpdatedAt = time.Now()
}

type DeliveryAttempt struct {
	ID             uint        `gorm:"primaryKey;autoIncrement"`
	NotificationID uint        `gorm:"index;not null"`
	Channel        ChannelType `gorm:"not null"`
	AttemptCount   int         `gorm:"not null;default:0"`
	LastAttemptAt  time.Time
	Status         DeliveryStatus `gorm:"not null"`
}

func (d *DeliveryAttempt) IncrementAttempt() {
	d.AttemptCount++
	d.LastAttemptAt = time.Now()
}

type UserPreference struct {
	UserID             string        `gorm:"primaryKey"`
	EnabledChannels    []ChannelType `gorm:"-"`
	EnabledChannelsRaw string        `gorm:"column:enabled_channels"`
	QuietHoursStart    int           `gorm:"default:0"`
	QuietHoursEnd      int           `gorm:"default:0"`
}

func (u *UserPreference) IsChannelAllowed(ch ChannelType) bool {
	for _, c := range u.EnabledChannels {
		if c == ch {
			return true
		}
	}
	return false
}

package repository

import (
	"notification_system/src/models"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Save(n *models.Notification) error
	FindByID(id uint) (*models.Notification, error)
	Update(n *models.Notification) error
}

type notificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepo{db: db}
}

func (r *notificationRepo) Save(n *models.Notification) error {
	return r.db.Create(n).Error
}

func (r *notificationRepo) FindByID(id uint) (*models.Notification, error) {
	var n models.Notification
	if err := r.db.First(&n, id).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *notificationRepo) Update(n *models.Notification) error {
	return r.db.Save(n).Error
}

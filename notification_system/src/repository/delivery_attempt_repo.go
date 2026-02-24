package repository

import (
	"notification_system/src/models"

	"gorm.io/gorm"
)

type DeliveryAttemptRepository interface {
	Save(d *models.DeliveryAttempt) error
	Update(d *models.DeliveryAttempt) error
	FindByNotificationID(notificationID uint) ([]models.DeliveryAttempt, error)
}

type deliveryAttemptRepo struct {
	db *gorm.DB
}

func NewDeliveryAttemptRepository(db *gorm.DB) DeliveryAttemptRepository {
	return &deliveryAttemptRepo{db: db}
}

func (r *deliveryAttemptRepo) Save(d *models.DeliveryAttempt) error {
	return r.db.Create(d).Error
}

func (r *deliveryAttemptRepo) Update(d *models.DeliveryAttempt) error {
	return r.db.Save(d).Error
}

func (r *deliveryAttemptRepo) FindByNotificationID(notificationID uint) ([]models.DeliveryAttempt, error) {
	var attempts []models.DeliveryAttempt
	if err := r.db.Where("notification_id = ?", notificationID).Find(&attempts).Error; err != nil {
		return nil, err
	}
	return attempts, nil
}

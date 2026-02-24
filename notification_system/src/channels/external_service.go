package channels

import "notification_system/src/models"

type ExternalService interface {
	Send(n models.Notification) error
	HealthCheck() error
}

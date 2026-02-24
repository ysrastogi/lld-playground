package channels

import (
	"fmt"
	"math/rand"
	"notification_system/src/models"
)

type EmailService struct{}

func NewEmailService() *EmailService { return &EmailService{} }

func (s *EmailService) Send(n models.Notification) error {
	if rand.Intn(10) == 0 { // ~10% failure rate
		return fmt.Errorf("[email] transient failure")
	}
	return nil
}

func (s *EmailService) HealthCheck() error {
	return nil
}

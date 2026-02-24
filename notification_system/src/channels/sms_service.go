package channels

import (
	"fmt"
	"math/rand"
	"notification_system/src/models"
)

type SMSService struct{}

func NewSMSService() *SMSService { return &SMSService{} }

func (s *SMSService) Send(n models.Notification) error {
	if rand.Intn(10) == 0 { 
		return fmt.Errorf("[sms] transient failure")
	}
	return nil
}

func (s *SMSService) HealthCheck() error {
	return nil
}

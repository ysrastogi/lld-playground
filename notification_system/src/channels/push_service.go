package channels

import (
	"fmt"
	"math/rand"
	"notification_system/src/models"
)

type PushService struct{}

func NewPushService() *PushService { return &PushService{} }

func (s *PushService) Send(n models.Notification) error {
	if rand.Intn(10) == 0 { // ~10% failure rate
		return fmt.Errorf("[push] transient failure")
	}
	return nil
}

func (s *PushService) HealthCheck() error {
	return nil
}

package channels

import (
	"fmt"
	"notification_system/src/models"
)

type InAppService struct{}

func NewInAppService() *InAppService { return &InAppService{} }

func (s *InAppService) Send(n models.Notification) error {
	return nil
}

func (s *InAppService) HealthCheck() error {
	fmt.Println("[inapp] health: OK")
	return nil
}

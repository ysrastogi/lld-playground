package services

import (
	"fmt"
	pubsub "notification_system/src/infrastructure/pub_sub"
	"notification_system/src/models"
	"notification_system/src/repository"
)

// NotificationService orchestrates notification creation and queuing.
type NotificationService struct {
	notifRepo   repository.NotificationRepository
	ruleEngine  *RuleEngine
	rateLimiter *RateLimiter
	idempotency *IdempotencyService
	broker      *pubsub.Broker
}

// NewNotificationService wires up the notification service.
func NewNotificationService(
	notifRepo repository.NotificationRepository,
	ruleEngine *RuleEngine,
	rateLimiter *RateLimiter,
	idempotency *IdempotencyService,
	broker *pubsub.Broker,
) *NotificationService {
	return &NotificationService{
		notifRepo:   notifRepo,
		ruleEngine:  ruleEngine,
		rateLimiter: rateLimiter,
		idempotency: idempotency,
		broker:      broker,
	}
}

func (s *NotificationService) CreateAndQueue(n *models.Notification, prefs *models.UserPreference) error {
	idempotencyKey := fmt.Sprintf("%s:%s:%s", n.UserID, n.Category, n.Title)
	if s.idempotency.IsDuplicate(idempotencyKey) {
		return nil
	}

	channels := s.ruleEngine.ResolveChannels(*n, prefs)

	var allowed []models.ChannelType
	for _, ch := range channels {
		if s.rateLimiter.Allow(n.UserID, string(ch)) {
			allowed = append(allowed, ch)
		}
	}

	if len(allowed) == 0 {
		return nil
	}

	n.Status = models.StatusQueued
	if err := s.notifRepo.Save(n); err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}
	s.idempotency.Store(idempotencyKey)

	for _, ch := range allowed {
		topic := string(ch) + "-notifications"
		event := pubsub.Event{
			ID: fmt.Sprintf("%d:%s", n.ID, ch),
			Payload: models.NotificationEvent{
				Notification: *n,
				Channels:     []models.ChannelType{ch},
			},
		}
		s.broker.Publish(topic, event)
	}

	return nil
}

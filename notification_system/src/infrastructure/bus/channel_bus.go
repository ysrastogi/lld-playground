package bus

import (
	"errors"
	"notification_system/src/models"
)

type ChannelBus struct {
	channel chan models.NotificationEvent
	handler func(models.NotificationEvent)
}

func NewChannelBus(buffer int, handler func(models.NotificationEvent)) *ChannelBus {
	return &ChannelBus{
		channel: make(chan models.NotificationEvent, buffer),
		handler: handler,
	}
}

func (cb *ChannelBus) Publish(event models.NotificationEvent) error {
	select {
	case cb.channel <- event:
		return nil
	default:
		return errors.New("channel bus is full")
	}
}

func (cb *ChannelBus) StartWorker() {
	go func() {
		for event := range cb.channel {
			cb.handler(event)
		}
	}()
}

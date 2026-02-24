package pubsub

import "fmt"

// Publish sends an event to all consumers subscribed to the given topic.
func (b *Broker) Publish(topic string, event Event) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, consumer := range b.consumers[topic] {
		select {
		case consumer.queue <- event:
		default:
			fmt.Printf("[broker] subscriber %s queue full, dropping event %s\n", consumer.name, event.ID)
		}
	}
}

// Subscribe registers a consumer under a topic.
func (b *Broker) Subscribe(topic string, consumer *Consumer) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.consumers[topic] = append(b.consumers[topic], consumer)
}

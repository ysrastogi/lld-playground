package pubsub

type Event struct {
	ID      string
	Payload any
	Retries int
}

type EventHandler func(event Event) error

type AckHandler func(event Event) error

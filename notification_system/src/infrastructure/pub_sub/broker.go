package pubsub

import "sync"

func NewBroker() *Broker {
	return &Broker{
		consumers: make(map[string][]*Consumer),
	}
}

func NewConsumer(name string, bufSize, workerCount, retryLimit int, handler EventHandler) *Consumer {
	return &Consumer{
		name:        name,
		queue:       make(chan Event, bufSize),
		handler:     handler,
		workerCount: workerCount,
		retryLimit:  retryLimit,
		dlq:         make(chan Event, bufSize),
	}
}

type Broker struct {
	mu        sync.Mutex
	consumers map[string][]*Consumer
}

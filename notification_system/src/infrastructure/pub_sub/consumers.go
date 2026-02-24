package pubsub

import (
	"context"
	"sync/atomic"
)

// Consumer reads events from its queue and processes them with a handler.
// On failure it retries up to retryLimit times, then routes to the DLQ.
type Consumer struct {
	name        string
	queue       chan Event
	handler     EventHandler
	workerCount int
	retryLimit  int
	dlq         chan Event
}

// Start launches workerCount goroutines that process events until ctx is done.
func (c *Consumer) Start(ctx context.Context) {
	for i := 0; i < c.workerCount; i++ {
		go c.worker(ctx)
	}
}

func (c *Consumer) worker(ctx context.Context) {
	for {
		select {
		case e := <-c.queue:
			if err := c.handler(e); err != nil {
				c.handleRetry(e)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *Consumer) handleRetry(e Event) {
	if e.Retries < c.retryLimit {
		e.Retries++
		c.queue <- e
	} else {
		c.dlq <- e
	}
}

// StartDLQLogger runs a goroutine that logs and counts DLQ events.
func (c *Consumer) StartDLQLogger(ctx context.Context, dlqCount *int64) {
	go func() {
		for {
			select {
			case <-c.dlq:
				atomic.AddInt64(dlqCount, 1)
				// Occasional summary log or completely silent
			case <-ctx.Done():
				return
			}
		}
	}()
}

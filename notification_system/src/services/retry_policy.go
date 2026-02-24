package services

import "time"

type RetryPolicy struct {
	BaseDelay time.Duration
	MaxDelay  time.Duration
}

func NewRetryPolicy(base, max time.Duration) *RetryPolicy {
	return &RetryPolicy{BaseDelay: base, MaxDelay: max}
}

func (r *RetryPolicy) NextDelay(attempt int) time.Duration {
	delay := r.BaseDelay
	for i := 0; i < attempt; i++ {
		delay *= 2
		if delay > r.MaxDelay {
			return r.MaxDelay
		}
	}
	return delay
}

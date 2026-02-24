package services

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu     sync.Mutex
	counts map[string][]time.Time
	limit  int
	window time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		counts: make(map[string][]time.Time),
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) Allow(userID, channel string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	key := userID + ":" + channel
	now := time.Now()
	cutoff := now.Add(-rl.window)

	ts := rl.counts[key]
	valid := ts[:0]
	for _, t := range ts {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	rl.counts[key] = valid

	if len(valid) >= rl.limit {
		return false
	}
	rl.counts[key] = append(rl.counts[key], now)
	return true
}

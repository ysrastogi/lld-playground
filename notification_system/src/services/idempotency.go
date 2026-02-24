package services

import "sync"

// IdempotencyService prevents processing duplicate notifications.
type IdempotencyService struct {
	mu   sync.RWMutex
	seen map[string]bool
}

func NewIdempotencyService() *IdempotencyService {
	return &IdempotencyService{seen: make(map[string]bool)}
}

// IsDuplicate returns true if this key has been seen before.
func (s *IdempotencyService) IsDuplicate(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.seen[key]
}

// Store records the key as processed.
func (s *IdempotencyService) Store(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seen[key] = true
}

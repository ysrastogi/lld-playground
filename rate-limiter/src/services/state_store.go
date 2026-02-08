package services

import (
	"rate-limiter/src/interfaces"
	"sync"
)

type StateStore struct {
	mu sync.RWMutex
	data map[string]*interfaces.LimiterState
}

func NewStateStore() *StateStore {
	return &StateStore{
		data: make(map[string]*interfaces.LimiterState),
	}
}

func (s *StateStore) GetState(key string) *interfaces.LimiterState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if state, exists := s.data[key]; exists {
		return state
	}
	return nil
}

func (s *StateStore) SetState(key string, state *interfaces.LimiterState){
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = state
}

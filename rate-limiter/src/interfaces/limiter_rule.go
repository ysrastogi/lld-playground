package interfaces

import (
	"rate-limiter/src/models"
	"sync"
	"time"
)

type LimiterRule interface {
	GetKey(ctx models.RequestContext) string
	Evaluate(ctx models.RequestContext, state *LimiterState, policy models.LimitPolicy) (bool, *LimiterState)
}

type TokenBucketRule struct {
	mu          sync.RWMutex
	LimitPolicy models.LimitPolicy
}

func (r *TokenBucketRule) GetKey(ctx models.RequestContext) string {
	return ctx.UserID
}

func (r *TokenBucketRule) Evaluate(ctx models.RequestContext, state *LimiterState, policy models.LimitPolicy) (bool, *LimiterState) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if time.Since(state.TokenBucket.LastRefillTime) > policy.Timeframe {
		state.TokenBucket.Tokens = policy.Requests
		state.TokenBucket.LastRefillTime = time.Now()
	}

	if state.TokenBucket.Tokens < 1 {
		return false, state
	}

	state.TokenBucket.Tokens -= 1
	return true, state
}

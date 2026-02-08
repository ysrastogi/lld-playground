package services

import (
	"rate-limiter/src/interfaces"
	"rate-limiter/src/models"
)

type RateLimiterOrchestrator struct {
	stateStore *StateStore
	rule       interfaces.LimiterRule
	policy     models.LimitPolicy
}

func NewRateLimiterOrchestrator(stateStore *StateStore, rule interfaces.LimiterRule, policy models.LimitPolicy) *RateLimiterOrchestrator {
	return &RateLimiterOrchestrator{
		stateStore: stateStore,
		rule:       rule,
		policy:     policy,
	}
}

func (o *RateLimiterOrchestrator) Allow(ctx models.RequestContext) bool {
	key := o.rule.GetKey(ctx)
	state := o.stateStore.GetState(key)
	if state == nil {
		state = &interfaces.LimiterState{}
	}
	allowed, newState := o.rule.Evaluate(ctx, state, o.policy)
	o.stateStore.SetState(key, newState)

	return allowed
}

func (o *RateLimiterOrchestrator) SetPolicy(policy models.LimitPolicy) {
	o.policy = policy
}

func (o *RateLimiterOrchestrator) SetRule(rule interfaces.LimiterRule) {
	o.rule = rule
}

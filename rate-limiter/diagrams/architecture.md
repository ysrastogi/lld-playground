# Rate Limiter Architecture

## Overview

A flexible, extensible rate limiting system designed using **Strategy Pattern** and **Interface Segregation** principles. The system supports multiple rate limiting algorithms (Token Bucket, Leaky Bucket, Fixed/Sliding/Rolling Window) with centralized state management and algorithm-agnostic orchestration.

## System Architecture

```mermaid
classDiagram

%% ===== INPUTS =====
class RequestContext {
  +UserID string
  +ApiKey string
  +IpAddress string
  +SetUserID(id string)
  +SetAPIKey(key string)
  +SetIPAddress(ip string)
}

class LimitPolicy {
  +Requests int
  +Timeframe Duration
  +MaxBurst int
  +ConcurrentRequests int
  +Entity entityType
  +SetRequests(r int)
  +SetTimeframe(t Duration)
  +SetMaxBurst(b int)
  +SetEntity(e entityType)
}

%% ===== CORE INTERFACE =====
class LimiterRule {
  <<interface>>
  +GetKey(ctx RequestContext) string
  +Evaluate(ctx, state, policy) (bool, *LimiterState)
}

%% ===== STATE STORE =====
class StateStore {
  -mu RWMutex
  -data map[string]*LimiterState
  +GetState(key string) *LimiterState
  +SetState(key string, state *LimiterState)
  +NewStateStore() *StateStore
}

%% ===== STATE STRUCT =====
class LimiterState {
  +TokenBucket TokenBucketState
  +LeakyBucket LeakyBucketState
}

%% ===== STATE IMPLEMENTATIONS =====
class TokenBucketState {
  +Tokens int
  +LastRefillTime Timestamp
}

class LeakyBucketState {
  +Water int
  +LastLeakTime Timestamp
}

class FixedWindowState {
  +WindowStart Timestamp
  +Count int
}

class SlidingWindowState {
  +Windows []Window
}

class RollingWindowState {
  +Timestamps []Timestamp
}

LimiterState *-- TokenBucketState
LimiterState *-- LeakyBucketState
StateStore --> LimiterState : stores

%% ===== LIMITER IMPLEMENTATIONS =====
class TokenBucketRule {
  -mu RWMutex
  +LimitPolicy LimitPolicy
  +GetKey(ctx) string
  +Evaluate(ctx, state, policy) (bool, *LimiterState)
}

class LeakyBucketRule {
  -mu RWMutex
  +LimitPolicy LimitPolicy
  +GetKey(ctx) string
  +Evaluate(ctx, state, policy) (bool, *LimiterState)
}

class FixedWindowRule {
  +GetKey(ctx) string
  +Evaluate(ctx, state, policy) (bool, *LimiterState)
}

class SlidingWindowRule {
  +GetKey(ctx) string
  +Evaluate(ctx, state, policy) (bool, *LimiterState)
}

class RollingWindowRule {
  +GetKey(ctx) string
  +Evaluate(ctx, state, policy) (bool, *LimiterState)
}

LimiterRule <|.. TokenBucketRule : implements
LimiterRule <|.. LeakyBucketRule : implements
LimiterRule <|.. FixedWindowRule : implements
LimiterRule <|.. SlidingWindowRule : implements
LimiterRule <|.. RollingWindowRule : implements

LimiterRule ..> LimiterState : uses
LimiterRule ..> RequestContext : reads
LimiterRule ..> LimitPolicy : enforces

%% ===== ORCHESTRATOR =====
class RateLimiterOrchestrator {
  -stateStore *StateStore
  -rule LimiterRule
  -policy LimitPolicy
  +NewRateLimiterOrchestrator() *RateLimiterOrchestrator
  +Allow(ctx RequestContext) bool
  +SetPolicy(policy LimitPolicy)
  +SetRule(rule LimiterRule)
}

RateLimiterOrchestrator *-- StateStore : contains
RateLimiterOrchestrator *-- LimiterRule : uses
RateLimiterOrchestrator *-- LimitPolicy : applies
RateLimiterOrchestrator ..> RequestContext : processes
RateLimiterOrchestrator ..> LimiterState : manages
```

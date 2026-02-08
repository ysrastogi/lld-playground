# Rate Limiter Components

## Overview

The rate limiter system consists of **6 core components** organized into three layers:

1. **Input Layer**: Request context and rate limit policies
2. **Logic Layer**: Rate limiting algorithms and orchestration
3. **Storage Layer**: State management and persistence

---

## Component Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    INPUT LAYER                          │
├──────────────────────┬──────────────────────────────────┤
│  RequestContext      │      LimitPolicy                 │
│  - UserID            │      - Requests                  │
│  - ApiKey            │      - Timeframe                 │
│  - IpAddress         │      - MaxBurst                  │
└──────────┬───────────┴────────────┬─────────────────────┘
           │                        │
           ▼                        ▼
┌─────────────────────────────────────────────────────────┐
│                    LOGIC LAYER                          │
├─────────────────────────────────────────────────────────┤
│              RateLimiterOrchestrator                    │
│              - Coordinates workflow                     │
│              - Manages state lifecycle                  │
├──────────────────────┬──────────────────────────────────┤
│   LimiterRule        │    Implementations:              │
│   (Interface)        │    - TokenBucketRule             │
│   - GetKey()         │    - LeakyBucketRule             │
│   - Evaluate()       │    - FixedWindowRule             │
│                      │    - SlidingWindowRule           │
└──────────┬───────────┴──────────────┬───────────────────┘
           │                          │
           ▼                          ▼
┌─────────────────────────────────────────────────────────┐
│                   STORAGE LAYER                         │
├──────────────────────┬──────────────────────────────────┤
│   StateStore         │      LimiterState                │
│   - GetState()       │      - TokenBucket               │
│   - SetState()       │      - LeakyBucket               │
│   - Thread-safe      │      - FixedWindow (planned)     │
└─────────────────────────────────────────────────────────┘
```

---

## 1. RequestContext (Input Component)

### Purpose
Encapsulates all information about an incoming request needed for rate limiting decisions.

### Fields
| Field | Type | Description | Usage |
|-------|------|-------------|-------|
| `UserID` | string | Authenticated user identifier | User-level rate limiting |
| `ApiKey` | string | API key for authentication | API key-based limiting |
| `IpAddress` | string | Client IP address | IP-based limiting, DDoS protection |

### Rate Limiting Keys

Requests can be rate-limited based on:

1. **User ID** (authenticated users)
   ```go
   key = "user:" + ctx.UserID
   // Example: "user:user-123"
   ```

2. **API Key / Client ID**
   ```go
   key = "apikey:" + ctx.ApiKey
   // Example: "apikey:sk-abc123"
   ```

3. **IP Address**
   ```go
   key = "ip:" + ctx.IpAddress
   // Example: "ip:192.168.1.1"
   ```

4. **Feature/API Endpoint**
   ```go
   key = "feature:" + endpoint
   // Example: "feature:/generate-report"
   ```

5. **Combination Keys**
   ```go
   key = ctx.UserID + ":" + feature
   // Example: "user-123:/expensive-api"
   ```

### Methods
```go
SetUserID(id string)         // Set user identifier
SetAPIKey(key string)        // Set API key
SetIPAddress(ip string)      // Set IP address
```

### Usage Example
```go
ctx := models.RequestContext{}
ctx.SetUserID("user-123")
ctx.SetAPIKey("sk-abc123")
ctx.SetIPAddress("192.168.1.1")
```

---

## 2. LimitPolicy (Configuration Component)

### Purpose
Defines the rate limiting rules and constraints to enforce.

### Fields
| Field | Type | Description | Example |
|-------|------|-------------|----------|
| `Requests` | int | Number of allowed requests | 100 |
| `Timeframe` | Duration | Time window duration | 60 seconds |
| `MaxBurst` | int | Maximum burst capacity | 150 |
| `ConcurrentRequests` | int | Max concurrent requests | 10 |
| `Entity` | entityType | Rate limit entity (User/IP) | User |

### Policy Types

#### 1. Fixed Rate Policy
```go
policy.SetRequests(100)
policy.SetTimeframe(60 * time.Second)
// 100 requests per minute
```

#### 2. Burst-Tolerant Policy
```go
policy.SetRequests(100)      // Sustained rate
policy.SetMaxBurst(200)      // Allow burst
policy.SetTimeframe(60 * time.Second)
// 100 req/min sustained, burst up to 200
```

#### 3. Feature-Specific Policy
```go
// Expensive API
expensivePolicy.SetRequests(10)
expensivePolicy.SetTimeframe(3600 * time.Second)
// 10 requests per hour

// Regular API
regularPolicy.SetRequests(1000)
regularPolicy.SetTimeframe(60 * time.Second)
// 1000 requests per minute
```

### Configuration Levels

| Level | Scope | Example |
|-------|-------|----------|
| Global | Entire system | 100K RPS system-wide |
| User Tier | Per user tier | Free: 100/min, Pro: 1000/min |
| Feature | Per API endpoint | /login: 5/min, /search: 100/min |
| Combined | User + Feature | user:123:/report: 10/hour |

---

## 3. LimiterRule Interface (Strategy Component)

### Purpose
Defines the contract for rate limiting algorithms using the **Strategy Pattern**.

### Interface Definition
```go
type LimiterRule interface {
    GetKey(ctx RequestContext) string
    Evaluate(ctx, state, policy) (bool, *LimiterState)
}
```

### Methods

#### `GetKey(ctx RequestContext) string`
- **Purpose**: Generate unique key for state lookup
- **Returns**: State key (e.g., "user-123", "ip:192.168.1.1")
- **Usage**: Determines state isolation granularity

#### `Evaluate(ctx, state, policy) (bool, *LimiterState)`
- **Purpose**: Evaluate if request should be allowed
- **Parameters**:
  - `ctx`: Request context
  - `state`: Current limiter state
  - `policy`: Rate limit policy
- **Returns**:
  - `bool`: true if allowed, false if denied
  - `*LimiterState`: Updated state

### Implementations

#### 1. TokenBucketRule ✅ (Implemented)
```go
type TokenBucketRule struct {
    mu          sync.RWMutex
    LimitPolicy LimitPolicy
}
```
**Algorithm**: Tokens added at fixed rate, consumed per request
**Best For**: Burst tolerance, API rate limiting
**Characteristics**:
- Allows bursts up to bucket capacity
- Smooth token refill
- Simple and efficient

#### 2. LeakyBucketRule 🚧 (Planned)
```go
type LeakyBucketRule struct {
    mu          sync.RWMutex
    LimitPolicy LimitPolicy
}
```
**Algorithm**: Requests queued, processed at fixed rate
**Best For**: Smooth traffic shaping
**Characteristics**:
- No bursts allowed
- Constant output rate
- Queue management required

#### 3. FixedWindowRule 🚧 (Planned)
```go
type FixedWindowRule struct {
    LimitPolicy LimitPolicy
}
```
**Algorithm**: Count requests in fixed time windows
**Best For**: Simple quota management
**Characteristics**:
- Easy to implement
- Boundary burst problem
- Memory efficient

#### 4. SlidingWindowRule 🚧 (Planned)
```go
type SlidingWindowRule struct {
    LimitPolicy LimitPolicy
}
```
**Algorithm**: Hybrid of fixed window + moving average
**Best For**: Accurate rate limiting
**Characteristics**:
- Smooth boundary transitions
- More complex
- Better accuracy than fixed window

#### 5. RollingWindowRule 🚧 (Planned)
```go
type RollingWindowRule struct {
    LimitPolicy LimitPolicy
}
```
**Algorithm**: Track individual request timestamps
**Best For**: Precise rate limiting
**Characteristics**:
- Most accurate
- Higher memory usage
- O(n) timestamp management

---

## 4. LimiterState (State Component)

### Purpose
Maintains algorithm-specific state for each entity (user/IP/key).

### Structure
```go
type LimiterState struct {
    TokenBucket TokenBucketState
    LeakyBucket LeakyBucketState
    // Future: FixedWindow, SlidingWindow, RollingWindow
}
```

### State Types

#### TokenBucketState
```go
type TokenBucketState struct {
    Tokens         int       // Current available tokens
    LastRefillTime time.Time // Last refill timestamp
}
```
**Memory**: ~16 bytes per entity

#### LeakyBucketState
```go
type LeakyBucketState struct {
    Water        int       // Current water level
    LastLeakTime time.Time // Last leak timestamp
}
```
**Memory**: ~16 bytes per entity

#### FixedWindowState (Planned)
```go
type FixedWindowState struct {
    WindowStart time.Time // Window start time
    Count       int       // Request count in window
}
```

#### SlidingWindowState (Planned)
```go
type SlidingWindowState struct {
    Windows []Window // Array of sub-windows
}
```

#### RollingWindowState (Planned)
```go
type RollingWindowState struct {
    Timestamps []time.Time // Individual request timestamps
}
```

### State Lifecycle

```
┌──────────────┐
│ Request      │
│ Arrives      │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ GetState()   │ ───No State──▶ Initialize
└──────┬───────┘                    │
       │ State Exists               │
       ▼                            ▼
┌──────────────┐            ┌──────────────┐
│ Evaluate()   │◀───────────│ New State    │
│ - Check      │            │ - Zeros      │
│ - Refill?    │            │ - Defaults   │
│ - Update     │            └──────────────┘
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ SetState()   │
│ Persist      │
└──────────────┘
```

---

## 5. StateStore (Storage Component)

### Purpose
Thread-safe centralized storage for rate limiter state.

### Implementation
```go
type StateStore struct {
    mu   sync.RWMutex
    data map[string]*interfaces.LimiterState
}
```

### Methods

#### `GetState(key string) *LimiterState`
- **Thread Safety**: Read lock (`RLock`)
- **Returns**: Pointer to state or nil if not found
- **Complexity**: O(1)

#### `SetState(key string, state *LimiterState)`
- **Thread Safety**: Write lock (`Lock`)
- **Operation**: Upsert state for key
- **Complexity**: O(1)

### Storage Backends

#### Current: In-Memory (map)
✅ **Pros**:
- Ultra-fast (< 1μs latency)
- Simple implementation
- No external dependencies

❌ **Cons**:
- Not distributed
- Lost on restart
- Limited by RAM

#### Future: Redis
✅ **Pros**:
- Distributed state
- Persistence
- Atomic operations (Lua scripts)
- Pub/Sub for coordination

❌ **Cons**:
- Network latency (~1-5ms)
- External dependency
- Requires setup

**Redis Implementation**:
```go
func (s *RedisStateStore) GetState(key string) (*LimiterState, error) {
    val, err := s.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil
    }
    var state LimiterState
    json.Unmarshal([]byte(val), &state)
    return &state, nil
}

func (s *RedisStateStore) SetState(key string, state *LimiterState) error {
    data, _ := json.Marshal(state)
    return s.client.Set(ctx, key, data, 0).Err()
}
```

#### Future: DynamoDB / Cassandra
✅ **Pros**:
- Highly scalable
- Multi-region
- Managed service

❌ **Cons**:
- Higher latency (~10-20ms)
- Cost
- Eventual consistency

---

## 6. RateLimiterOrchestrator (Orchestration Component)

### Purpose
Central coordinator that ties all components together.

### Structure
```go
type RateLimiterOrchestrator struct {
    stateStore *StateStore      // State persistence
    rule       LimiterRule       // Algorithm strategy
    policy     LimitPolicy       // Rate limit policy
}
```

### Core Methods

#### `Allow(ctx RequestContext) bool`
**Workflow**:
1. Extract key from context using `rule.GetKey()`
2. Retrieve state from `stateStore.GetState()`
3. Initialize state if nil
4. Evaluate request using `rule.Evaluate()`
5. Persist updated state using `stateStore.SetState()`
6. Return allow/deny decision

**Pseudo-code**:
```go
func (o *Orchestrator) Allow(ctx RequestContext) bool {
    key := o.rule.GetKey(ctx)              // "user-123"
    state := o.stateStore.GetState(key)    // Get or nil
    if state == nil {
        state = &LimiterState{}             // Initialize
    }
    allowed, newState := o.rule.Evaluate(ctx, state, o.policy)
    o.stateStore.SetState(key, newState)   // Persist
    return allowed
}
```

#### `SetPolicy(policy LimitPolicy)`
**Purpose**: Update rate limit policy at runtime
**Usage**: Dynamic configuration changes
**Effect**: Immediate - next request uses new policy

#### `SetRule(rule LimiterRule)`
**Purpose**: Swap rate limiting algorithm
**Usage**: A/B testing, algorithm migration
**Effect**: Strategy pattern in action

### Orchestrator Responsibilities

| Responsibility | Description |
|----------------|-------------|
| **State Lifecycle** | Get → Evaluate → Set workflow |
| **Initialization** | Create new state for first requests |
| **Coordination** | Bind rule + policy + state |
| **Abstraction** | Hide algorithm complexity from callers |
| **Extensibility** | Support multiple algorithms via interface |

---

## Component Interactions

### Data Flow
```
Request
   ↓
RequestContext ──────┐
                     │
                     ▼
          RateLimiterOrchestrator
                     │
        ┌────────────┼────────────┐
        ▼            ▼            ▼
   LimiterRule  StateStore  LimitPolicy
        │            │
        ▼            ▼
    Algorithm    LimiterState
        │            │
        └────────────┘
             │
             ▼
        Allow/Deny
```

### Thread Safety Matrix

| Component | Concurrency Model | Mechanism |
|-----------|------------------|------------|
| RequestContext | Read-only | Immutable |
| LimitPolicy | Read-only | Immutable |
| LimiterRule | Thread-safe | sync.RWMutex |
| LimiterState | Protected | Accessed via locks |
| StateStore | Thread-safe | sync.RWMutex |
| Orchestrator | Thread-safe | Delegates to StateStore |

---

## Extensibility Points

### 1. Custom Rate Limiting Keys
```go
func (r *CustomRule) GetKey(ctx RequestContext) string {
    return fmt.Sprintf("%s:%s", ctx.UserID, ctx.Feature)
}
```

### 2. Custom Algorithms
```go
type CustomRule struct {
    // Custom fields
}

func (r *CustomRule) Evaluate(...) (bool, *LimiterState) {
    // Custom logic
}
```

### 3. Storage Backend Swap
```go
type StateStoreInterface interface {
    GetState(key string) *LimiterState
    SetState(key string, state *LimiterState)
}

// Implementations:
// - InMemoryStateStore
// - RedisStateStore
// - DynamoDBStateStore
```

### 4. Multi-Layer Rate Limiting
```go
type MultiLayerOrchestrator struct {
    userLimiter    *RateLimiterOrchestrator
    ipLimiter      *RateLimiterOrchestrator
    globalLimiter  *RateLimiterOrchestrator
}

func (m *MultiLayerOrchestrator) Allow(ctx) bool {
    return m.userLimiter.Allow(ctx) &&
           m.ipLimiter.Allow(ctx) &&
           m.globalLimiter.Allow(ctx)
}
```

---

## Performance Characteristics

| Operation | Current (In-Memory) | Redis | DynamoDB |
|-----------|---------------------|-------|----------|
| GetState | < 1 μs | 1-5 ms | 10-20 ms |
| SetState | < 1 μs | 1-5 ms | 10-20 ms |
| Evaluate | < 10 μs | N/A | N/A |
| **Total Latency** | **< 20 μs** | **2-10 ms** | **20-40 ms** |
| **Throughput** | 1M+ RPS | 100K RPS | 10K RPS |
| **Scalability** | Vertical | Horizontal | Horizontal |

---

## Summary

The rate limiter consists of **6 loosely-coupled components**:

1. ✅ **RequestContext** - Request metadata
2. ✅ **LimitPolicy** - Rate limit rules
3. ✅ **LimiterRule** - Algorithm interface (Strategy Pattern)
4. ✅ **LimiterState** - Algorithm state
5. ✅ **StateStore** - Thread-safe storage
6. ✅ **RateLimiterOrchestrator** - Central coordinator

**Design Principles Applied**:
- Strategy Pattern (swappable algorithms)
- Dependency Injection (testable)
- Interface Segregation (minimal interfaces)
- Single Responsibility (one job per component)
- Open/Closed (extend without modify)
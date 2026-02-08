# Rate Limiter Sequence Diagrams

## 1. Basic Request Flow - Allowed Request

```mermaid
sequenceDiagram
    actor User
    participant API as API Gateway
    participant Orch as RateLimiterOrchestrator
    participant Rule as LimiterRule
    participant Store as StateStore
    participant Backend as Backend Service

    User->>API: HTTP Request
    API->>Orch: Allow(RequestContext)
    Note over API,Orch: ctx = {UserID, ApiKey, IP}
    
    Orch->>Rule: GetKey(ctx)
    Rule-->>Orch: "user-123"
    
    Orch->>Store: GetState("user-123")
    Store-->>Orch: *LimiterState{Tokens: 5}
    
    Orch->>Rule: Evaluate(ctx, state, policy)
    Note over Rule: Check: Tokens >= 1?
    Note over Rule: Decrement: Tokens = 4
    Rule-->>Orch: (allowed=true, newState)
    
    Orch->>Store: SetState("user-123", newState)
    Orch-->>API: true
    
    API->>Backend: Forward Request
    Backend-->>API: Response (200 OK)
    API-->>User: Response + Headers<br/>X-RateLimit-Remaining: 4
```

---

## 2. Rate Limit Exceeded - Blocked Request

```mermaid
sequenceDiagram
    actor User
    participant API as API Gateway
    participant Orch as RateLimiterOrchestrator
    participant Rule as LimiterRule
    participant Store as StateStore

    User->>API: HTTP Request (6th in window)
    API->>Orch: Allow(RequestContext)
    
    Orch->>Rule: GetKey(ctx)
    Rule-->>Orch: "user-123"
    
    Orch->>Store: GetState("user-123")
    Store-->>Orch: *LimiterState{Tokens: 0}
    
    Orch->>Rule: Evaluate(ctx, state, policy)
    Note over Rule: Check: Tokens < 1
    Note over Rule: DENY - No change to state
    Rule-->>Orch: (allowed=false, state)
    
    Orch->>Store: SetState("user-123", state)
    Orch-->>API: false
    
    API-->>User: 429 Too Many Requests<br/>X-RateLimit-Limit: 5<br/>X-RateLimit-Remaining: 0<br/>Retry-After: 8s
```

---

## 3. Token Bucket Refill Flow

```mermaid
sequenceDiagram
    actor User
    participant Orch as RateLimiterOrchestrator
    participant Rule as TokenBucketRule
    participant Store as StateStore
    participant Time as Time.Now()

    Note over User,Time: User exhausted tokens at T=0s
    Note over User,Time: Now attempting request at T=11s
    
    User->>Orch: Allow(RequestContext)
    Orch->>Rule: GetKey(ctx)
    Rule-->>Orch: "user-123"
    
    Orch->>Store: GetState("user-123")
    Store-->>Orch: State{Tokens:0, LastRefill: T=0s}
    
    Orch->>Rule: Evaluate(ctx, state, policy)
    Rule->>Time: time.Since(LastRefillTime)
    Time-->>Rule: 11 seconds
    
    Note over Rule: 11s > 10s (Timeframe)
    Note over Rule: REFILL TRIGGERED
    Rule->>Rule: Tokens = 5 (policy.Requests)
    Rule->>Rule: LastRefillTime = Now()
    
    Note over Rule: Check: Tokens >= 1? ✓
    Rule->>Rule: Decrement: Tokens = 4
    Rule-->>Orch: (allowed=true, newState)
    
    Orch->>Store: SetState("user-123", newState)
    Orch-->>User: true (ALLOWED)
```

---

## 4. Multi-User Concurrent Requests

```mermaid
sequenceDiagram
    participant U1 as User 1
    participant U2 as User 2
    participant Orch as RateLimiterOrchestrator
    participant Store as StateStore
    
    par User 1 Request
        U1->>Orch: Allow(ctx{UserID: "user-1"})
        Orch->>Store: GetState("user-1")
        Store-->>Orch: State{Tokens: 5}
        Note over Orch: Evaluate & Decrement
        Orch->>Store: SetState("user-1", {Tokens: 4})
        Orch-->>U1: ALLOWED
    and User 2 Request
        U2->>Orch: Allow(ctx{UserID: "user-2"})
        Orch->>Store: GetState("user-2")
        Store-->>Orch: State{Tokens: 5}
        Note over Orch: Evaluate & Decrement
        Orch->>Store: SetState("user-2", {Tokens: 4})
        Orch-->>U2: ALLOWED
    end
    
    Note over U1,Store: ✓ Independent state per user<br/>✓ No interference between users
```

---

## 5. First Request (State Initialization)

```mermaid
sequenceDiagram
    actor User as New User
    participant Orch as RateLimiterOrchestrator
    participant Rule as LimiterRule
    participant Store as StateStore

    User->>Orch: Allow(ctx{UserID: "user-999"})
    Orch->>Rule: GetKey(ctx)
    Rule-->>Orch: "user-999"
    
    Orch->>Store: GetState("user-999")
    Store-->>Orch: nil (not found)
    
    Note over Orch: State is nil → Initialize
    Orch->>Orch: state = &LimiterState{}
    Note over Orch: TokenBucket: {Tokens: 0, LastRefill: zero}
    
    Orch->>Rule: Evaluate(ctx, state, policy)
    Note over Rule: LastRefillTime is zero<br/>Time since = ∞ > Timeframe
    Note over Rule: REFILL on first request
    Rule->>Rule: Tokens = 5<br/>LastRefillTime = Now()
    Rule->>Rule: Decrement: Tokens = 4
    Rule-->>Orch: (allowed=true, newState)
    
    Orch->>Store: SetState("user-999", newState)
    Orch-->>User: ALLOWED
```

---

## 6. Dynamic Policy Update Flow

```mermaid
sequenceDiagram
    actor Admin
    participant API as Admin API
    participant Orch as RateLimiterOrchestrator
    actor User

    Admin->>API: POST /admin/rate-limits<br/>{requests: 10, timeframe: "1m"}
    API->>Orch: SetPolicy(newPolicy)
    Note over Orch: Update policy: 5→10 req/min
    Orch-->>API: Success
    API-->>Admin: 200 OK
    
    Note over User,Orch: New policy effective immediately
    
    User->>Orch: Allow(ctx)
    Note over Orch: Uses NEW policy (10 req/min)<br/>for evaluation
    Orch-->>User: Evaluated with updated limits
```

---

## 7. Algorithm Switching at Runtime

```mermaid
sequenceDiagram
    actor Admin
    participant Orch as RateLimiterOrchestrator
    participant TB as TokenBucketRule
    participant LB as LeakyBucketRule
    actor User

    Note over Orch,TB: Currently using Token Bucket
    
    Admin->>Orch: SetRule(LeakyBucketRule)
    Orch->>Orch: rule = LeakyBucketRule
    Orch-->>Admin: Switched to Leaky Bucket
    
    User->>Orch: Allow(ctx)
    Orch->>LB: GetKey(ctx)
    LB-->>Orch: "user-123"
    Note over Orch,LB: Now using Leaky Bucket algorithm
    Orch->>LB: Evaluate(ctx, state, policy)
    LB-->>Orch: (allowed, newState)
    Orch-->>User: Result
    
    Note over Orch,LB: ✓ Zero downtime switch<br/>✓ Strategy pattern in action
```

---

## 8. Error Handling - StateStore Failure

```mermaid
sequenceDiagram
    actor User
    participant API as API Gateway
    participant Orch as RateLimiterOrchestrator
    participant Store as StateStore (Redis)
    participant Backend as Backend Service

    User->>API: HTTP Request
    API->>Orch: Allow(ctx)
    
    Orch->>Store: GetState("user-123")
    Note over Store: ⚠️ Redis Timeout/Failure
    Store--xOrch: Error (connection timeout)
    
    alt Critical API (Fail Open)
        Note over Orch: Policy: Allow on failure
        Orch-->>API: true (ALLOWED)
        API->>Backend: Forward Request
        Backend-->>User: Response (200 OK)
    else Non-Critical API (Fail Closed)
        Note over Orch: Policy: Deny on failure
        Orch-->>API: false (BLOCKED)
        API-->>User: 503 Service Unavailable<br/>Rate limiter unavailable
    end
```

---

## 9. Distributed System - Multi-Instance Coordination

```mermaid
sequenceDiagram
    participant U1 as User (Instance 1)
    participant API1 as API Instance 1
    participant Orch1 as Orchestrator 1
    participant Redis as Redis (Shared State)
    participant Orch2 as Orchestrator 2
    participant API2 as API Instance 2
    participant U2 as User (Instance 2)

    Note over API1,API2: Both instances share Redis state
    
    U1->>API1: Request 1
    API1->>Orch1: Allow(ctx{UserID: "user-123"})
    Orch1->>Redis: INCR user:123:count
    Redis-->>Orch1: 1 (within limit)
    Orch1-->>U1: ALLOWED
    
    U2->>API2: Request 2 (same user)
    API2->>Orch2: Allow(ctx{UserID: "user-123"})
    Orch2->>Redis: INCR user:123:count
    Redis-->>Orch2: 2 (within limit)
    Orch2-->>U2: ALLOWED
    
    Note over Redis: ✓ Atomic operations ensure consistency<br/>✓ Both instances see same state
```

---

## 10. IP-Based + User-Based Combined Limiting

```mermaid
sequenceDiagram
    actor User
    participant Orch as RateLimiterOrchestrator
    participant UserRule as UserRateLimiter
    participant IPRule as IPRateLimiter
    participant Store as StateStore

    User->>Orch: Request from IP: 192.168.1.1
    
    Note over Orch: Check 1: User-level limit
    Orch->>UserRule: Allow(ctx{UserID: "user-123"})
    UserRule->>Store: GetState("user:user-123")
    Store-->>UserRule: {Tokens: 5}
    UserRule-->>Orch: ALLOWED
    
    Note over Orch: Check 2: IP-level limit
    Orch->>IPRule: Allow(ctx{IP: "192.168.1.1"})
    IPRule->>Store: GetState("ip:192.168.1.1")
    Store-->>IPRule: {Tokens: 100}
    IPRule-->>Orch: ALLOWED
    
    Note over Orch: Both checks passed
    Orch-->>User: 200 OK
    
    alt IP Limit Exceeded
        IPRule-->>Orch: BLOCKED
        Orch-->>User: 429 Too Many Requests<br/>(IP-based rate limit)
    end
```

---

## Key Observations

### Thread Safety
- `sync.RWMutex` in StateStore prevents race conditions
- `sync.RWMutex` in Rule implementations protects state modifications
- Concurrent requests handled safely

### State Isolation
- Each user/IP/entity has independent state
- State keyed by `GetKey()` method output
- No interference between different entities

### Atomic Operations
- Get → Evaluate → Set is NOT atomic in current implementation
- For distributed systems, use Redis Lua scripts or CAS operations

### Initialization
- First request auto-initializes state
- Zero-value LastRefillTime triggers immediate refill
- No pre-provisioning needed

### Failure Modes
- Configurable fail-open or fail-closed
- Critical APIs can bypass on errors
- Non-critical APIs deny on errors
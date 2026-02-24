# Notification System Pattern Discussion

## 1. Publish-Subscribe Pattern (Pub/Sub)

**Implementation**: `Broker`, `Consumers`, and `Topics`
- **Why**: Decouples the producers (API layer that accepts notifications) from the consumers (background workers that hit external APIs).
- **Benefit**:
  - Without Pub/Sub, `CreateAndQueue` would have to block until an email, SMS, or Push attempt succeeds, introducing terrible user-facing latency.
  - Allows different channels to scale independently. For instance, if emails take 300ms but SMS takes 50ms, the system can allocate more workers to `email-notifications` without affecting SMS throughput.

## 2. Strategy Pattern

**Implementation**: The `channels.ExternalService` interface and its various implementations (`EmailService`, `SMSService`, `PushService`).
- **Why**: The consumer logic (`makeHandler` in `main.go`) does not need to know the specific mechanics of *how* a notification is sent. It simply invokes `.Send(n)` on the generic interface.
- **Benefit**:
  - Extremely easy to add a new channel in the future (e.g., Slack, WhatsApp).
  - Isolates the external failure points.

## 3. Idempotency Key Pattern

**Implementation**: `IdempotencyService`
- **Why**: Distributed systems occasionally resend identical events (due to client retries or network blips).
- **How**: Before a notification is saved or queued, the orchestrator generates a key from the core payload (`fmt.Sprintf("%s:%s:%s", n.UserID, n.Category, n.Title)`). If this key is already in the cache, the system silently ignores it and returns success.
- **Benefit**: Protects end-users from terrifying spam (e.g., receiving 5 identical banking transaction alerts).

## 4. Intercepting Filter / Pipeline Pattern

**Implementation**: The `NotificationService` encapsulates `RuleEngine.ResolveChannels` and `RateLimiter.Allow`.
- **Why**: An incoming notification doesn't know where it needs to go. The business requirements (Routing based on Category, filtering based on User Preferences, dropping based on Rate Limit) act as a series of filters.
- **Benefit**: 
  - Modularizes the business rules.
  - Any step of the pipeline can dynamically alter the downstream behavior (e.g., dropping SMS if the user hit their 5/min limit, but retaining Email).

## 5. Worker Pool Pattern

**Implementation**: `pub_sub/consumer.go` spawning `N` goroutines.
- **Why**: Listening to a Go channel sequentially operates at $O(1)$ concurrency. Spawning unbounded goroutines inside a loop can starve memory allocations.
- **How**: The consumer spins up a fixed pool of workers (`workerCount`). The workers constantly pull from the shared channel buffered queue.
- **Benefit**: Throttles the concurrency to an optimal level, preventing the OS from context-switching excessively and avoiding "connection refused" errors from overloading external third-party APIs.

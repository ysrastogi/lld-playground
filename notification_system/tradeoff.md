# Notification System Trade-offs

## 1. Data Store: SQLite In-Memory vs PostgreSQL / Cassandra

### ✅ Decision: SQLite (In-Memory)
In the current implementation (`main.go`), the repository uses SQLite via an in-memory connection string (`file::memory:?cache=shared&mode=memory&_journal_mode=OFF&_synchronous=OFF`).

**Pros**:
- **Setup Speed**: Zero external dependencies to run the simulation.
- **Throughput**: Extremely high write performance since I/O is completely bypassed (allows testing purely CPU-bound concurrency limits).
- **Simplicity**: Single-binary execution without Docker compose dependencies.

**Cons (Production Warning)**:
- **Durability**: Completely volatile. A crash wipes out the entire history of notifications.
- **Horizontal Scaling**: An in-memory database cannot be shared natively across multiple distributed API servers.

**Alternative**: Cassandra for massively high write throughput and horizontally scalable persistent storage.

---

## 2. In-Memory Pub/Sub vs External Broker (Kafka / RabbitMQ)

### ✅ Decision: Custom Go Channel-based Pub/Sub
The `infrastructure/pub_sub` package creates a custom message broker using buffered Go channels (`chan Event`).

**Pros**:
- **Zero Network Latency**: Enables millions of events per second on a decent multi-core machine.
- **Ease of Use**: No need to configure Zookeeper, Kafka brokers, or handle network partitions in the playground.
- **Resource Efficient**: Pure goroutine scheduling makes maximum use of OS threads.

**Cons (Production Warning)**:
- **Unreliable Delivery**: If the process crashes, the `50,000` buffered events in the channel are irrecoverably lost.
- **Lack of Replayability**: Unlike Kafka's log segments, consumed channel events disappear. If a deployment introduces a bug that incorrectly consumes messages, they cannot be readily replayed.

---

## 3. Order of Operations: Save-then-Publish vs Outbox Pattern

### ✅ Decision: Synchronous Save then Publish
In `services/notification_service.go`, the application performs:
```go
1. s.notifRepo.Save(n)
2. s.broker.Publish(topic, event)
```

**Pros**:
- **Low Complexity**: Straightforward implementation.
- **Low Latency**: Directly writing to DB and pushing to memory broker is extremely fast.

**Cons (Production Warning)**:
- **Dual Write Problem**: If the application crashes immediately *after* Step 1 but *before* Step 2, the notification shows as `Queued` in the DB but is *never* actually published to the consumers.
- **Inconsistent State**: External systems won't see the event despite the DB recording success.

**Alternative**: The **Transactional Outbox Pattern**. The notification and the "event" are written to a single atomic database table. A separate background process (e.g., Debezium) reads the table log and publishes to the broker, ensuring At-Least-Once delivery semantics.

---

## 4. Polling vs Event-Driven Dead Letter Queue (DLQ)

### ✅ Decision: In-Memory DLQ Channels
The `pub_sub/consumer.go` routes exhausted attempts to an internal `dlq` Go channel.

**Pros**:
- Avoids blocking the main consumer workers.
- Instant feedback.

**Cons (Production Warning)**:
- Hard to inspect interactively. In production, DLQs are routed to separate persistent topics where SREs can inspect, debug, and manually redrive the payloads after a fix is deployed.

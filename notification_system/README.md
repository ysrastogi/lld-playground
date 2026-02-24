# Notification System

A highly concurrent, multi-channel notification dispatcher simulating massive throughput (1 Million+ events) with retry mechanisms, idempotency, and asynchronous processing.

## 🎯 Architecture Summary
This project acts as an advanced distributed system mock-up built in Go. It accepts abstract notifications, processes user preferences, applies rate limits, and dispatches tasks asynchronously via an internal Pub/Sub broker to a pool of worker consumers handling specific channels (Email, SMS, Push, In-App).

## 📚 Documentation Links
- [Requirements](requirements.md)
- [Component Architecture](components.md)
- [Sequence Diagrams](diagrams/sequence.md)
- [Design Trade-offs](tradeoff.md)
- [Pattern Discussion](patterns.md)

## ✨ Key Features
*   **Massive Concurrency Simulator**: The `main.go` file boots workers that stress-test the pipeline with millions of records.
*   **Rule Engine Routing**: Dynamically decides the optimal channel (e.g., Marketing -> Email + InApp) and strips out channels based on User Preferences.
*   **Rate Limiting**: Protects users from being spammed on a specific channel.
*   **In-Memory Pub/Sub**: High-performance buffered channels paired with fixed-size worker pools to efficiently drain tasks.
*   **Reliability Mocking**: Simulates external API unresponsiveness triggers Dead Letter Queues (DLQs) and automated retries.
*   **Idempotency Handling**: Blocks duplicate upstream payloads gracefully.

## 🚀 Running the System
The project includes a 1 Million+ high-scale simulation in `main.go`. Over 10,000 simulated users will have notifications queued aggressively across 100 concurrent producers, which are drained by thousands of consumer go-routines tracking throughput.

```bash
cd notification_system

# Download dependencies (Gorm, SQLite driver, etc)
go mod tidy

# Run the simulation
go run main.go
```

Expect output simulating progress of queued messages, detailed deliveries, simulated network failures, and eventually a DLQ drain.

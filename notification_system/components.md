# Notification System Components

## Overview

The Notification System is designed to handle multi-channel notification delivery with high concurrency, retry mechanisms, and customizable routing rules. It consists of **5 core layers**.

1. **Input / Orchestration Layer**: Manages the incoming requests.
2. **Core Logic Layer**: Resolves rules, user preferences, limits, and deduplication.
3. **Messaging Layer**: Buffers and routes events asynchronously.
4. **Channel Execution Layer**: Interfaces with external providers to deliver notifications.
5. **Storage Layer**: Persists notification states and delivery attempts.

---

## Component Architecture

```
┌────────────────────────────────────────────────────────┐
│             INPUT / ORCHESTRATION LAYER                │
├────────────────────────────────────────────────────────┤
│                  NotificationService                   │
│          - Orchestrates notification creation          │
│          - Coordinates rule engine & rate limiting     │
└──────────────────────────┬─────────────────────────────┘
                           │
                           ▼
┌────────────────────────────────────────────────────────┐
│                   CORE LOGIC LAYER                     │
├──────────────┬───────────────┬─────────────────────────┤
│  RuleEngine  │  RateLimiter  │   IdempotencyService    │
│  - Routing   │  - Throttling │   - Deduplication       │
│  - Prefs     │               │                         │
└─────────┬────┴───────┬───────┴─────────┬───────────────┘
          │            │                 │
          ▼            ▼                 ▼
┌────────────────────────────────────────────────────────┐
│                   MESSAGING LAYER                      │
├────────────────────────────────────────────────────────┤
│                Pub/Sub Message Broker                  │
│                - Topics (e.g., email-notifications)    │
│                - Consumer Worker Pools                 │
│                - Dead Letter Queues (DLQ)              │
└──────────────────────────┬─────────────────────────────┘
                           │
                           ▼
┌────────────────────────────────────────────────────────┐
│               CHANNEL EXECUTION LAYER                  │
├───────┬────────────┬─────────────┬─────────────────────┤
│ Email │    SMS     │    Push     │       In-App        │
└───────┴────────────┴─────────────┴─────────────────────┘
                           │
                           ▼
┌────────────────────────────────────────────────────────┐
│                   STORAGE LAYER                        │
├────────────────────────────┬───────────────────────────┤
│   NotificationRepository   │ DeliveryAttemptRepository │
│   - Stores raw payloads    │ - Tracks retries/status   │
└────────────────────────────┴───────────────────────────┘
```

---

## 1. NotificationService (Orchestration Component)

### Purpose
Acts as the central integration point. It receives incoming requests to send notifications and orchestrates the filtering and routing steps before dispatching them to the message broker.

### Responsibilities
- Verifies if the request is duplicate using the **IdempotencyService**.
- Resolves target channels via the **RuleEngine**.
- Applies throttling via the **RateLimiter**.
- Persists the initial state in the Database.
- Publishes events to the relevant Pub/Sub topics.

---

## 2. Core Logic Components

### IdempotencyService
- **Purpose**: Prevents duplicate notifications from being sent.
- **Mechanism**: Maintains a cache (in-memory or Redis) using a composite key (e.g., `UserID:Category:Title`).
- **Flow**: Returns successfully without generating duplicate events.

### RuleEngine
- **Purpose**: Determines the required notification channels based on the category (e.g., `Transaction`, `Marketing`) and User Preferences.
- **Logic**:
  - Maps `models.NotificationCategory` to default `ChannelType` array.
  - Filters out channels that the user has opted out of (`UserPreference`).

### RateLimiter
- **Purpose**: Prevents a user from being spammed by limiting the number of notifications sent within a specific timeframe per channel.
- **Mechanism**: Token bucket or Sliding window applied per user and per channel type.

---

## 3. Messaging Layer (Pub/Sub)

### Broker
- **Purpose**: Central hub that routes published events to interested consumers.
- **Structure**: Manages separate channels/topics tailored for specific services (e.g., `email-notifications`, `sms-notifications`).

### Consumer
- **Purpose**: Processes events asynchronously.
- **Architecture**:
  - **Worker Pools**: Concurrently drain events (e.g., 500 workers per channel) to support 1M+ throughput.
  - **Retry Limits**: Configurable max retries.
  - **Dead Letter Queue (DLQ)**: Captures events that repeatedly fail.

---

## 4. Channel Services

### Interface Definition
```go
type ExternalService interface {
    Send(notification *models.Notification) error
}
```

### Implementations
- **EmailService**: Sends rich-text notifications.
- **SMSService**: Sends short, urgent text notifications.
- **PushService**: Sends push notifications to mobile devices.
- **InAppService**: Persists notifications to be shown within the client application.

---

## 5. Storage Layer

### NotificationRepository
- **Purpose**: Source of truth for notification details, payload, and initial status (`Created`, `Queued`).

### DeliveryAttemptRepository
- **Purpose**: Highly granular tracking of where a notification is in its lifecycle across different channels. 
- **States Captured**: `DeliveryPending`, `DeliverySuccess`, `DeliveryFailed`.
- **Usage**: Critical for analytical systems to compute delivery rates and latency.

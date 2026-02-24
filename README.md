# LLD Playground

A collection of Low-Level Design (LLD) implementations and architecture case studies. This repository serves as a playground for practicing object-oriented design, system modeling, and implementation patterns in various programming languages.

## Projects

### 1. [Car Rental Lite](car-rental-lite/)
A simplified car rental system implementation (like Turo or Airbnb for cars).

*   **Language**: Go
*   **Key Concepts**: 
    *   Concurrency Control (Double booking prevention)
    *   Inventory Management (Availability Slots)
    *   Clean Architecture (Service-Repository Pattern)
    *   In-Memory Data Structures
*   **Documentation**:
    *   [Requirements](car-rental-lite/requirement.md)
    *   [Architecture Diagrams](car-rental-lite/diagrams/)
    *   [Design Trade-offs](car-rental-lite/tradeoff.md)

### 2. [Rate Limiter](rate-limiter/)
A rate limiter implementation supporting multiple algorithms (Token Bucket, Leaky Bucket, Fixed Window, Sliding Window, Rolling Window).
*   **Language**: Go
*   **Key Concepts**:
    *   Rate Limiting Algorithms
    *   State Management
    *   Orchestrator Pattern
    *   In-Memory State Store
*   **Documentation**:
    *   [Requirements](rate-limiter/requirement.md)
    *   [Architecture Diagrams](rate-limiter/diagrams/architecture.md)
    *   [Components](rate-limiter/components.md)
    *   [Sequence Diagrams](rate-limiter/diagrams/sequence.md)
    *   [Design Trade-offs](rate-limiter/tradeoff.md)
    *   [Pattern Discussion](rate-limiter/patterns.md)

### 3. [LRU Cache](lru-cache/)
An implementation of a Least Recently Used (LRU) Cache.
*   **Language**: Go
*   **Key Concepts**:
    *   LRU Eviction Policy
    *   Doubly Linked List for Usage Tracking
    *   Hash Map for O(1) Access
    *   Thread Safety for Concurrent Access
*   **Documentation**:
    *   [Requirements](lru-cache/requirements.md)
    *   [Architecture Diagrams](lru-cache/diagrams/architecture.md)
    *   [Components](lru-cache/components.md)

### 4. [LFU Cache](lfu-cache/)
An implementation of a Least Frequently Used (LFU) Cache.
*   **Language**: Go
*   **Key Concepts**:
    *   LFU Eviction Policy
    *   Frequency Tracking
    *   Doubly Linked List for Usage Tracking
    *   Hash Map for O(1) Access
    *   Thread Safety for Concurrent Access
*   **Documentation**:
    *   [Requirements](lfu-cache/requirements.md)
    *   [Architecture Diagrams](lfu-cache/diagrams/architecture.md)
    *   [Components](lfu-cache/components.md)

### 5. [URL Shortener](url-shortener/)
An implementation of a URL shortener service.
*   **Language**: Go
*   **Key Concepts**:
    *   URL Shortening
    *   Rate Limiting
    *   Pub Sub Messaging
    *   Idempotency
    *   Retry Infrastructure
*   **Documentation**:
    *   [Requirements](url-shortener/requirements.md)
    *   [Architecture Diagrams](url-shortener/diagrams/architecture.md)
    *   [Components](url-shortener/components.md)
    *   [Sequence Diagrams](url-shortener/diagrams/sequence.md)
    *   [Design Trade-offs](url-shortener/tradeoff.md)
    *   [Pattern Discussion](url-shortener/patterns.md)

### 6. [Notification System](notification_system/)
An implementation of a notification system.
*   **Language**: Go
*   **Key Concepts**:
    *   Notification System
    *   Rate Limiting
    *   Caching
    *   Analytics
*   **Documentation**:
    *   [Requirements](notification_system/requirements.md)
    *   [Architecture Diagrams](notification_system/diagrams/architecture.png)
    *   [Components](notification_system/components.md)
    *   [Sequence Diagrams](notification_system/diagrams/sequence.md)
    *   [Design Trade-offs](notification_system/tradeoff.md)
    *   [Pattern Discussion](notification_system/patterns.md)

## Goals
*   Implement common LLD interview problems.
*   Experiment with different design patterns.
*   Document tradeoffs between different implementation approaches.
*   Future Projects:
    - [x] Car Rental System
    - [x] Rate Limiter
    - [x] LRU Cache
    - [x] LFU Cache
    - [x] URL Shortener
    - [x] Notification System
    - [ ] InMemory Cache
    - [ ] Payment Flow
    - [ ] Job Scheduler
    - [ ] Distributed Locking System
    - [ ] Cache Invalidation System
    - [ ] Search Autocomplete
    - [ ] Ride Matching
    - [ ] Inventory Management

## Usage
Each project folder contains its own `src` and instructions. Navigate to the specific folder to run the code.

```bash
cd car-rental-lite
go run src/main.go
```

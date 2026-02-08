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
    *   [Components] (rate-limiter/components.md)
    *   [Sequence Diagrams](rate-limiter/diagrams/sequence.md)
    *   [Design Trade-offs](rate-limiter/tradeoff.md)
    *   [Pattern Discussion](rate-limiter/patterns.md)


## Goals
*   Implement common LLD interview problems.
*   Experiment with different design patterns.
*   Document tradeoffs between different implementation approaches.
*   Future Projects:
    - [x] Car Rental System
    - [x] Rate Limiter
    - [ ] LRU Cache
    - [ ] LFU Cache
    - [ ] URL Shortener
    - [ ] Notification System
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

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

## Goals
*   Implement common LLD interview problems.
*   Experiment with different design patterns.
*   Document tradeoffs between different implementation approaches.
*   Future Projects:
    - [x] Car Rental System
    - [ ] Rate Limiter
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

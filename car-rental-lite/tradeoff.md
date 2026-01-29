# Design Trade-offs

This document outlines the architectural and design decisions made for the **Car Rental Lite** system, comparing the chosen approach against alternatives.

## 1. Concurrency Control: In-Memory Mutex vs. Database Locking

### **Decision**: In-Memory RWMutex (`sync.RWMutex`)
We implemented the repository layer using in-memory maps protected by Go's `sync.RWMutex`.

| Approach | Pros | Cons |
|:---|:---|:---|
| **In-Memory Mutex** (Chosen) | **Simple to implement** for LLD demo.<br/>**Extremely fast** (nanosecond latency).<br/>Strict consistency within a single instance. | **Not persistent**; data is lost on restart.<br/>**Does not scale** horizontally (multiple replicas would have separate memory states). |
| **Database Pessimistic Locking** (`SELECT FOR UPDATE`) | strict consistency across distributed instances.<br/>Standard industry practice for bookings. | Higher latency (network + disk I/O).<br/>More complex setup (requires SQL DB). |

### **Trade-off Impact**
For a "Lite" version single-instance playground, In-Memory is sufficient. For Production, we would swap the `BookingRepository` implementation to use a Postgres transaction with Row Locking to prevent double bookings across multiple server instances.

---

## 2. Search Implementation: Naive Filtering vs. Geospatial Index

### **Decision**: Naive Iteration + Filtering
The `Search` method iterates through all registered cars to check availability.

| Approach | Pros | Cons |
|:---|:---|:---|
| **Naive Iteration** (Chosen) | Zero external dependencies.<br/>Easy to debug and write.<br/>Sufficient for small datasets (< 10k cars). | **O(N) Complexity**.<br/>Performance degrades linearly with inventory size.<br/>Data fetch inefficiency (fetching all cars to app layer). |
| **Geospatial Index** (PostGIS / Elasticsearch) | **O(log N)** complexity.<br/>Efficiently handles "cars near me". | Requires managing a search infrastructure.<br/>Requires synchronization between "Booking Source of Truth" and "Search Index". |

### **Trade-off Impact**
We favored code simplicity over performance for this stage. In a real-world scenario, checking availability for *every* nearby car via the booking table is expensive. We would likely introduce a pre-calculated "Availability Index" or use a Search Engine.

---

## 3. Availability Model: Slots vs. Simple State

### **Decision**: Explicit Availability Slots
Hosts must "add availability" for a car to be searchable.

| Approach | Pros | Cons |
|:---|:---|:---|
| **Explicit Availability Slots** (Chosen) | Flexible for gig-economy hosts (e.g., "I only rent my car on weekends").<br/>Allows "blackout" dates easily. | **Complex validation logic**: Must check if requested range falls *within* a slot AND doesn't overlap bookings.<br/>Higher storage requirements. |
| **Default Available** (Calendar Blocklists) | Simpler default state.<br/>Easier for commercial fleets (cars always available unless rented). | Harder to model sporadic availability. |

### **Trade-off Impact**
We chose the "Slot" model to support the peer-to-peer (Airbnb-style) nature of the requirements, accepting the higher complexity in the `Search` logic.

# System Requirements

## 1. Overview
The **Car Rental Lite** system connects Car Owners (Hosts) with Renters (Users), allowing users to find and book vehicles for short-term rental periods.

## 2. Functional Requirements

### 2.1 Search & Discovery
- **Search by Location**: Users must be able to find cars within a specific radius of a city or coordinate.
- **Search by Availability**: Users must only see cars that are available for the *entire* requested duration.
- **Filtering**: (Planned) Filter by Price, Car Type (Sedan, SUV), and Amenities.

### 2.2 Booking Management
- **Create Booking**: Users can reserve a car for a specific time window.
- **Prevent Double Booking**: The system must strictly ensure a car cannot be booked by two users for overlapping time slots.
- **Cancel Booking**: Users can cancel a booking (subject to policy windows).

### 2.3 Inventory Management (Host)
- **Car Listing**: Hosts can register cars with details (Make, Model, Year, Price).
- **Availability Slots**: Hosts explicitly define when a car is available for rent.

### 2.4 Trip Operations
- **Inspection Reports**: Support for logging car condition (Fuel, Odometer, Photos) at Pickup and Dropoff.
- **Damage Claims**: Mechanism to file claims with evidence if damage is detected strictly between Pickup and Dropoff reports.

## 3. Non-Functional Requirements

### 3.1 Consistency
- **Strict Consistency** is required for Booking functionality. We cannot tolerate overbooking scenarios.
- **Eventual Consistency** is acceptable for Search results (e.g., a car booked 1ms ago might still appear in search for a brief moment, but booking will fail).

### 3.2 Performance
- **Search Latency**: Target < 200ms for car search queries.
- **Booking Throughput**: Support moderate concurrency for the same vehicle (handling race conditions gracefully).

### 3.3 Scalability
- The system design should allow for separating Read (Search) and Write (Booking) paths in the future.

## 4. Out of Scope
- Payment Processing Integration.
- Real-time GPS tracking of vehicles.
- Insurance Policy generation.
- Identity Verification (KYC).

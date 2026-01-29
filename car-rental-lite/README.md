# Car Rental System

## Actors
- **User (Renter)**: Searches for and books cars.
- **Host (Owner)**: Lists cars and manages availability.

## Goals

### User Goals
- User wants a car for defined time window
- User prefers nearby cars
- User prefers well rated/maintained cars
- Delivering to his location
- Before Use Report 

### Owner Goals
- Host wants to list the car on the platform, based on the availability
- Before and After Use report

## System Requirements
- Allow searching cars by location + availability.
- Allow booking cars for a timed window
- Prevent Double Booking
- Store car condition reports (pre & post booking)
- Record damage claims with evidence
- Associate responsibility with a trip

## Implementation Details (Go)

The system is implemented in Go using a **Clean Architecture** pattern.

### Project Structure
- `src/models`: Core domain entities (User, Car, Booking, etc.). pure data structs.
- `src/repositories`: Interfaces for data persistence and an In-Memory implementation for testing.
- `src/services`: Business logic layer.
    - `InventoryService`: Availability checks and Car search.
    - `BookingService`: Transactional booking logic and double-booking prevention.
- `src/main.go`: Entry point demonstrating the wiring and usage.

### Running the System
```bash
cd car-rental-lite
go run src/main.go
```

# Sequence Diagrams

## 1. Car Search Flow
This flow demonstrates how the system filters cars based on location, host availability, and existing bookings.

```mermaid
sequenceDiagram
    participant U as User
    participant IS as InventoryService
    participant CR as CarRepository
    participant IR as InventoryRepository
    participant BR as BookingRepository

    U->>IS: Search(Location, Dates)
    
    IS->>CR: SearchCars(Location, Radius)
    CR-->>IS: List[Car] (Filtered by Loc)

    loop For Each Car
        IS->>IR: HasAvailabilitySlot(CarID, Dates)
        IR-->>IS: true/false
        
        opt If Has Slot
            IS->>BR: HasOverlappingBooking(CarID, Dates)
            BR-->>IS: true/false (Conflict Check)
        end
    end

    IS-->>U: List[Car] (Available & Conflict Free)
```

## 2. Booking a Car
This flow ensures data consistency and prevents double booking.

```mermaid
sequenceDiagram
    participant U as User
    participant BS as BookingService
    participant CR as CarRepository
    participant IR as InventoryRepository
    participant BR as BookingRepository

    U->>BS: CreateBooking(UserID, CarID, Dates)

    BS->>CR: GetCar(CarID)
    CR-->>BS: Car Details

    par Availability Checks
        BS->>IR: HasAvailabilitySlot(CarID, Dates)
        IR-->>BS: true (Host has listed it)
        
        BS->>BR: HasOverlappingBooking(CarID, Dates)
        BR-->>BS: false (No existing booking)
    end

    BS->>BS: Calculate Price

    BS->>BR: CreateBooking(Booking)
    BR-->>BS: Success

    BS-->>U: Booking Confirmation
```

# Components

The system is divided into three main layers: Models, Services, and Repositories.

## 1. Domain Models (`src/models`)

The business entities that hold state.

### User & Host
- **User**: Represents a renter. Attributes: `ID`, `Name`, `DriverLicense`.
- **Host**: Represents a car owner. Attributes: `ID`, `Name`.

### Inventory
- **Car**: The core asset.
  - Attributes: `ID`, `HostID`, `Type` (Sedan, SUV, etc.), `Location` (Lat/Long, City), `PricePerDay`.
- **AvailabilitySlot**: A time range explicitly defined by a Host when a car is available.
  - Attributes: `CarID`, `StartTime`, `EndTime`.

### Transactional
- **Booking**: A reservation made by a User.
  - Attributes: `ID`, `CarID`, `UserID`, `Status` (Confirmed, Cancelled), `TotalPrice`.
- **InspectionReport**: Records condition before/after trip.
  - Attributes: `BookingID`, `Type` (Pickup/Dropoff), `FuelLevel`, `Odometer`, `Images`.
- **DamageClaim**: A claim filed if damage is found.
  - Attributes: `BookingID`, `Description`, `EvidenceImages`, `Status`.

## 2. Services (`src/services`)

Contain the business logic.

### InventoryService
- **Responsibilities**:
  - Registering new cars.
  - Managing availability slots (Hosts adding time).
  - **Search**: Finding cars that match a location *and* have an availability slot covering the requested duration *and* have no conflicting bookings.

### BookingService
- **Responsibilities**:
  - Validating booking requests.
  - **Concurrency Control**: Ensuring a car isn't double-booked.
  - Calculation of Total Price.
  - Creating the Booking record.

## 3. Repositories (`src/repositories`)

Abstract the data storage.

- **UserRepository**: CRUD for Users and Hosts.
- **CarRepository**: Store and retrieve Car details. Includes Geolocation search (stubbed).
- **InventoryRepository**: Manage `AvailabilitySlot`s.
- **BookingRepository**: Manage `Booking`s. Critical method: `HasOverlappingBooking` for conflict detection.

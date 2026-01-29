package repositories

import (
	"car-rental-lite/src/models"
	"context"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id string) (*models.User, error)
	CreateHost(ctx context.Context, host *models.Host) error
	GetHost(ctx context.Context, id string) (*models.Host, error)
}

type CarRepository interface {
	CreateCar(ctx context.Context, car *models.Car) error
	GetCar(ctx context.Context, id string) (*models.Car, error)
	GetCarsByHost(ctx context.Context, hostID string) ([]*models.Car, error)
	SearchCars(ctx context.Context, location models.Location, radiusKm float64) ([]*models.Car, error)
}

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *models.Booking) error
	GetBooking(ctx context.Context, id string) (*models.Booking, error)
	UpdateBookingStatus(ctx context.Context, id string, status models.BookingStatus) error
	GetBookingsForCar(ctx context.Context, carID string, from, to time.Time) ([]*models.Booking, error)
	HasOverlappingBooking(ctx context.Context, carID string, start, end time.Time) (bool, error)
}

type InventoryRepository interface {
	AddAvailability(ctx context.Context, slot *models.AvailabilitySlot) error
	GetAvailability(ctx context.Context, carID string, from, to time.Time) ([]*models.AvailabilitySlot, error)
	RemoveAvailability(ctx context.Context, slotID string) error
	HasAvailabilitySlot(ctx context.Context, carID string, start, end time.Time) (bool, error)
}

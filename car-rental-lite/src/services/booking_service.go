package services

import (
	"car-rental-lite/src/models"
	"car-rental-lite/src/repositories"
	"context"
	"errors"
	"time"
)

type BookingService struct {
	bookingRepo   repositories.BookingRepository
	inventoryRepo repositories.InventoryRepository
	carRepo       repositories.CarRepository
}

func NewBookingService(bRepo repositories.BookingRepository, iRepo repositories.InventoryRepository, cRepo repositories.CarRepository) *BookingService {
	return &BookingService{
		bookingRepo:   bRepo,
		inventoryRepo: iRepo,
		carRepo:       cRepo,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, userID, carID string, start, end time.Time) (*models.Booking, error) {
	// 1. Validate times
	if start.After(end) {
		return nil, errors.New("invalid booking duration")
	}

	// 2. Check Car existence
	car, err := s.carRepo.GetCar(ctx, carID)
	if err != nil {
		return nil, err
	}
	if !car.IsActive {
		return nil, errors.New("car is not currently active")
	}

	// 3. Check Availability
	// a. Host created a slot?
	hasSlot, err := s.inventoryRepo.HasAvailabilitySlot(ctx, carID, start, end)
	if err != nil {
		return nil, err
	}
	if !hasSlot {
		return nil, errors.New("host has not made the car available for these dates")
	}

	// b. No overlapping booking?
	hasOverlap, err := s.bookingRepo.HasOverlappingBooking(ctx, carID, start, end)
	if err != nil {
		return nil, err
	}
	if hasOverlap {
		return nil, errors.New("car is already booked for these dates")
	}

	// 4. Calculate Price
	days := end.Sub(start).Hours() / 24.0
	// Ceiling or partial day logic usually applies, simple multiply for now
	if days < 1 {
		days = 1
	}
	totalPrice := days * car.PricePerDay

	// 5. Create Booking
	booking := &models.Booking{
		// ID generated in repo or here
		UserID:     userID,
		CarID:      carID,
		StartTime:  start,
		EndTime:    end,
		Status:     models.BookingStatusConfirmed, // Or Pending if payment is separate
		TotalPrice: totalPrice,
		CreatedAt:  time.Now(),
	}

	if err := s.bookingRepo.CreateBooking(ctx, booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, bookingID string) error {
	booking, err := s.bookingRepo.GetBooking(ctx, bookingID)
	if err != nil {
		return err
	}

	if booking.Status == models.BookingStatusCompleted || booking.Status == models.BookingStatusCancelled {
		return errors.New("booking cannot be cancelled")
	}

	if time.Until(booking.StartTime) < 24*time.Hour {
		// return errors.New("cannot cancel within 24 hours")
	}

	return s.bookingRepo.UpdateBookingStatus(ctx, bookingID, models.BookingStatusCancelled)
}

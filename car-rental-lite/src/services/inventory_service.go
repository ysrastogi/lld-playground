package services

import (
	"car-rental-lite/src/models"
	"car-rental-lite/src/repositories"
	"context"
	"errors"
	"time"
)

type InventoryService struct {
	carRepo       repositories.CarRepository
	inventoryRepo repositories.InventoryRepository
	bookingRepo   repositories.BookingRepository
}

func NewInventoryService(carRepo repositories.CarRepository, invRepo repositories.InventoryRepository, bookRepo repositories.BookingRepository) *InventoryService {
	return &InventoryService{
		carRepo:       carRepo,
		inventoryRepo: invRepo,
		bookingRepo:   bookRepo,
	}
}

func (s *InventoryService) RegisterCar(ctx context.Context, car *models.Car) error {
	if car.HostID == "" {
		return errors.New("host ID is required")
	}
	// Basic validation could go here
	return s.carRepo.CreateCar(ctx, car)
}

func (s *InventoryService) AddAvailability(ctx context.Context, carID string, start, end time.Time) error {
	if start.After(end) {
		return errors.New("start time must be before end time")
	}

	// Check if car exists
	_, err := s.carRepo.GetCar(ctx, carID)
	if err != nil {
		return err
	}

	slot := &models.AvailabilitySlot{
		// ID would normally be generated here or in repo
		CarID:     carID,
		StartTime: start,
		EndTime:   end,
	}
	return s.inventoryRepo.AddAvailability(ctx, slot)
}

func (s *InventoryService) Search(ctx context.Context, location models.Location, radiusKm float64, start, end time.Time) ([]*models.Car, error) {
	// 1. Find cars in location
	cars, err := s.carRepo.SearchCars(ctx, location, radiusKm)
	if err != nil {
		return nil, err
	}

	// 2. Filter by availability
	// This is a naive N+1 implementation. In a real system, the DB query would handle this join.
	availableCars := make([]*models.Car, 0)
	for _, car := range cars {
		hasSlot, err := s.inventoryRepo.HasAvailabilitySlot(ctx, car.ID, start, end)
		if err != nil || !hasSlot {
			continue
		}

		hasOverlap, err := s.bookingRepo.HasOverlappingBooking(ctx, car.ID, start, end)
		if err != nil || hasOverlap {
			continue
		}

		availableCars = append(availableCars, car)
	}

	return availableCars, nil
}

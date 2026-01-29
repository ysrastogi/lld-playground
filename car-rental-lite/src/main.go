package main

import (
	"car-rental-lite/src/models"
	"car-rental-lite/src/repositories"
	"car-rental-lite/src/services"
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()

	repo := repositories.NewInMemoryRepo()
	invService := services.NewInventoryService(repo, repo, repo)
	bookService := services.NewBookingService(repo, repo, repo)

	host := &models.Host{ID: "host1", Name: "Host Alice", Email: "alice@example.com"}
	repo.CreateHost(ctx, host)

	user := &models.User{ID: "user1", Name: "Renter Bob", Email: "bob@example.com"}
	repo.CreateUser(ctx, user)

	car := &models.Car{
		ID:           "car1",
		HostID:       "host1",
		Make:         "Toyota",
		Model:        "Camry",
		Year:         2022,
		PricePerDay:  100.0,
		IsActive:     true,
		Location:     models.Location{City: "San Francisco"},
	}
	err := invService.RegisterCar(ctx, car)
	if err != nil {
		panic(err)
	}

	startAvail := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endAvail := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	err = invService.AddAvailability(ctx, car.ID, startAvail, endAvail)
	if err != nil {
		panic(err)
	}
	fmt.Println("Availability Set: Jan 1 - Jan 10")

	searchStart := time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)
	searchEnd := time.Date(2026, 1, 4, 10, 0, 0, 0, time.UTC)
	cars, err := invService.Search(ctx, models.Location{City: "San Francisco"}, 10.0, searchStart, searchEnd)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Search Found %d cars available\n", len(cars))

	fmt.Println("Booking Car 1...")
	booking, err := bookService.CreateBooking(ctx, user.ID, car.ID, searchStart, searchEnd)
	if err != nil {
		fmt.Printf("Booking failed: %v\n", err)
	} else {
		fmt.Printf("Booking Success! ID: %s, Price: %.2f\n", booking.CarID, booking.TotalPrice)
	}

	fmt.Println("Attempting Double Booking...")
	_, err = bookService.CreateBooking(ctx, "user2", car.ID, searchStart, searchEnd)
	if err != nil {
		fmt.Printf("Double Booking Correctly Failed: %v\n", err)
	} else {
		fmt.Println("Error: Double Booking Allowed!")
	}
}

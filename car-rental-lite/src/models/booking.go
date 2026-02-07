package models

import "time"

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
	BookingStatusCompleted BookingStatus = "COMPLETED"
)

type Booking struct {
	ID         string
	CarID      string
	UserID     string
	StartTime  time.Time
	EndTime    time.Time
	Status     BookingStatus
	TotalPrice float64
	CreatedAt  time.Time
}

type AvailabilitySlot struct {
	ID        string
	CarID     string
	StartTime time.Time
	EndTime   time.Time
}

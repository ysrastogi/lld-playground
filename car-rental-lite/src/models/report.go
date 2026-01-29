package models

import "time"

type InspectionType string

const (
	InspectionTypePickup  InspectionType = "PICKUP"
	InspectionTypeDropoff InspectionType = "DROPOFF"
)

type InspectionReport struct {
	ID          string
	BookingID   string
	InspectorID string // UserID or HostID depending on who does it
	Type        InspectionType
	Notes       string
	ImageURLs   []string
	FuelLevel   float64 // Percentage 0-100
	Odometer    int
	CreatedAt   time.Time
}

type ClaimStatus string

const (
	ClaimStatusSubmitted   ClaimStatus = "SUBMITTED"
	ClaimStatusUnderReview ClaimStatus = "UNDER_REVIEW"
	ClaimStatusApproved    ClaimStatus = "APPROVED"
	ClaimStatusRejected    ClaimStatus = "REJECTED"
)

type DamageClaim struct {
	ID                string
	BookingID         string
	Description       string
	EvidenceImageURLs []string
	EstimatedCost     float64
	Status            ClaimStatus
	CreatedAt         time.Time
}

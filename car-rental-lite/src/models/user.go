package models

import "time"

type User struct {
	ID            string
	Name          string
	Email         string
	PhoneNumber   string
	DriverLicense string
	CreatedAt     time.Time
}

type Host struct {
	ID          string
	Name        string
	Email       string
	PhoneNumber string
	CreatedAt   time.Time
}

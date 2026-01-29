package models

type CarType string

const (
	CarTypeSedan   CarType = "SEDAN"
	CarTypeSUV     CarType = "SUV"
	CarTypeCompact CarType = "COMPACT"
	CarTypeLuxury  CarType = "LUXURY"
	CarTypeVan     CarType = "VAN"
)

type Location struct {
	Latitude  float64
	Longitude float64
	Address   string
	City      string
	ZipCode   string
}

type Car struct {
	ID           string
	HostID       string
	Make         string
	Model        string
	Year         int
	Type         CarType
	Location     Location
	PricePerDay  float64
	LicensePlate string
	IsActive     bool
}

package repositories

import (
	"car-rental-lite/src/models"
	"context"
	"errors"
	"sync"
	"time"
)

type InMemoryRepo struct {
	users    map[string]*models.User
	hosts    map[string]*models.Host
	cars     map[string]*models.Car
	bookings map[string]*models.Booking
	slots    map[string][]*models.AvailabilitySlot
	mu       sync.RWMutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		users:    make(map[string]*models.User),
		hosts:    make(map[string]*models.Host),
		cars:     make(map[string]*models.Car),
		bookings: make(map[string]*models.Booking),
		slots:    make(map[string][]*models.AvailabilitySlot),
	}
}

func (r *InMemoryRepo) CreateUser(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}
func (r *InMemoryRepo) GetUser(ctx context.Context, id string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}
func (r *InMemoryRepo) CreateHost(ctx context.Context, host *models.Host) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.hosts[host.ID] = host
	return nil
}
func (r *InMemoryRepo) GetHost(ctx context.Context, id string) (*models.Host, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.hosts[id]
	if !ok {
		return nil, errors.New("host not found")
	}
	return h, nil
}

func (r *InMemoryRepo) CreateCar(ctx context.Context, car *models.Car) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cars[car.ID] = car
	return nil
}
func (r *InMemoryRepo) GetCar(ctx context.Context, id string) (*models.Car, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.cars[id]
	if !ok {
		return nil, errors.New("car not found")
	}
	return c, nil
}
func (r *InMemoryRepo) GetCarsByHost(ctx context.Context, hostID string) ([]*models.Car, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*models.Car
	for _, c := range r.cars {
		if c.HostID == hostID {
			res = append(res, c)
		}
	}
	return res, nil
}
func (r *InMemoryRepo) SearchCars(ctx context.Context, location models.Location, radiusKm float64) ([]*models.Car, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*models.Car
	for _, c := range r.cars {
		if c.IsActive {
			res = append(res, c)
		}
	}
	return res, nil
}

func (r *InMemoryRepo) CreateBooking(ctx context.Context, booking *models.Booking) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bookings[booking.ID] = booking
	return nil
}
func (r *InMemoryRepo) GetBooking(ctx context.Context, id string) (*models.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.bookings[id]
	if !ok {
		return nil, errors.New("booking not found")
	}
	return b, nil
}
func (r *InMemoryRepo) UpdateBookingStatus(ctx context.Context, id string, status models.BookingStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.bookings[id]
	if !ok {
		return errors.New("booking not found")
	}
	b.Status = status
	return nil
}
func (r *InMemoryRepo) GetBookingsForCar(ctx context.Context, carID string, from, to time.Time) ([]*models.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*models.Booking
	for _, b := range r.bookings {
		if b.CarID == carID {
			// Check overlap
			if b.StartTime.Before(to) && b.EndTime.After(from) {
				res = append(res, b)
			}
		}
	}
	return res, nil
}
func (r *InMemoryRepo) HasOverlappingBooking(ctx context.Context, carID string, start, end time.Time) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, b := range r.bookings {
		if b.CarID == carID && b.Status != models.BookingStatusCancelled {
			// Overlap logic: (StartA < EndB) and (EndA > StartB)
			if start.Before(b.EndTime) && end.After(b.StartTime) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (r *InMemoryRepo) AddAvailability(ctx context.Context, slot *models.AvailabilitySlot) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.slots[slot.CarID] = append(r.slots[slot.CarID], slot)
	return nil
}
func (r *InMemoryRepo) GetAvailability(ctx context.Context, carID string, from, to time.Time) ([]*models.AvailabilitySlot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.slots[carID], nil
}
func (r *InMemoryRepo) RemoveAvailability(ctx context.Context, slotID string) error {
	return nil
}
func (r *InMemoryRepo) HasAvailabilitySlot(ctx context.Context, carID string, start, end time.Time) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	slots := r.slots[carID]
	for _, s := range slots {
		if !s.StartTime.After(start) && !s.EndTime.Before(end) {
			return true, nil
		}
	}
	return false, nil
}

package entity

import (
	"github.com/Rhymond/go-money"
	"time"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCanceled  BookingStatus = "canceled"
)

type Booking struct {
	GuestId      int
	CheckInDate  time.Time
	CheckOutDate time.Time
	TotalPrice   *money.Money
	Status       BookingStatus
	CreatedAt    time.Time
}

/*

func (s BookingStatus) Valid() bool {
	switch s {
	case BookingStatusPending, BookingStatusConfirmed, BookingStatusCanceled:
		return true
	default:
		return false
	}
}

// Scan для поддержки чтения из SQL
func (s *BookingStatus) Scan(value interface{}) error {
	if val, ok := value.(string); ok {
		*s = BookingStatus(val)
		if !s.Valid() {
			return errors.New("invalid booking status")
		}
		return nil
	}
	return errors.New("invalid type for booking status")
}

type Booking struct {
	GuestId      int
	CheckInDate  time.Time
	CheckOutDate time.Time
	TotalPrice   *money.Money
	Status       BookingStatus
}
*/

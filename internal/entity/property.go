package entity

import (
	"github.com/Rhymond/go-money"
	"time"
)

type Property struct {
	OwnerId       int
	Title         string
	Description   string
	PricePerNight *money.Money
	MaxGuests     int
	CreatedAt     time.Time
}

type Location struct {
	Address string
	City    string
	Country string
}

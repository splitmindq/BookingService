package entity

import "time"

type Review struct {
	PropertyId int
	GuestId    int
	Rating     int
	Comment    string
	CreatedAt  time.Time
}

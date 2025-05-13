package entity

import "time"

type BlockReason string

// CREATE TYPE blocked_reason AS ENUM ('booking', 'maintenance', 'personal', 'other');
const (
	BookingReason     BlockReason = "booking"
	MaintenanceReason BlockReason = "maintenance"
	PersonalReason    BlockReason = "personal"
	OtherReason       BlockReason = "other"
)

type BlockedDates struct {
	PropertyId int
	StartDate  time.Time
	EndDate    time.Time
	Reason     BlockReason
}

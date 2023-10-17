package storage

import "time"

type Event struct {
	ID          int
	Title       string
	StartDate   time.Time
	EndDate     time.Time
	Description string
	OwnerID     int
	RemindAt    time.Time
}

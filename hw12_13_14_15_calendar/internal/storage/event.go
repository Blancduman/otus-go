package storage

import "time"

type Event struct {
	ID          int64
	Title       string
	StartDate   time.Time
	EndDate     time.Time
	Description string
	OwnerID     int64
	RemindIn    time.Time
}

func (e Event) ExactChanges(c Event) (Event, bool) {
	var isChanged bool

	change := Event{
		ID: c.ID,
	}

	if e.Title != "" && e.Title != c.Title {
		change.Title = e.Title
		isChanged = true
	} else {
		change.Title = c.Title
	}

	if e.Description != "" && e.Description != c.Description {
		change.Description = e.Description
		isChanged = true
	} else {
		change.Description = c.Description
	}

	if e.OwnerID != 0 && e.OwnerID != c.OwnerID {
		change.OwnerID = e.OwnerID
		isChanged = true
	} else {
		change.OwnerID = c.OwnerID
	}

	if e.StartDate.String() != "1970-01-01 00:00:00 +0000 UTC" && e.StartDate != c.StartDate {
		change.StartDate = e.StartDate
		isChanged = true
	} else {
		change.StartDate = c.StartDate
	}

	if e.EndDate.String() != "1970-01-01 00:00:00 +0000 UTC" && e.EndDate != c.EndDate {
		change.EndDate = e.EndDate
		isChanged = true
	} else {
		change.EndDate = c.EndDate
	}

	if e.RemindIn.String() != "0s" && e.RemindIn != c.RemindIn {
		change.RemindIn = e.RemindIn
		isChanged = true
	} else {
		change.RemindIn = c.RemindIn
	}

	return change, isChanged
}

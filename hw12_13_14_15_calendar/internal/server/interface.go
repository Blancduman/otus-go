package server

import (
	"context"
	"time"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
}

type Application interface {
	CreateEvent(
		ctx context.Context,
		title string,
		description string,
		ownerID int64,
		startDate time.Time,
		endDate time.Time,
		remindIn time.Duration,
	) (int64, error)
	EditEvent(
		ctx context.Context,
		id int64,
		title string,
		description string,
		ownerID int64,
		startDate time.Time,
		endDate time.Time,
		remindIn time.Duration,
	) error
	RemoveEvent(ctx context.Context, id int64) error
	GetEvent(ctx context.Context, id int64) (storage.Event, error)
	GetDateTimeRangeEvents(ctx context.Context, startDate time.Time, endDate time.Time) ([]storage.Event, error)
}

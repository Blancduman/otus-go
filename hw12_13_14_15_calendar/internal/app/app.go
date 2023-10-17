package app

import (
	"context"
	"time"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
}

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Add(event *storage.Event) error
	Edit(event *storage.Event) error
	Delete(id int) error
	Get(id int) (storage.Event, error)
	GetDateTimeRange(from, to time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(_ context.Context, _ string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

package app

import (
	"context"
	"time"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/pkg/errors"
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
	Add(ctx context.Context, event *storage.Event) error
	Edit(ctx context.Context, event *storage.Event) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (storage.Event, error)
	GetDateTimeRange(ctx context.Context, from, to time.Time) ([]storage.Event, error)
	GetReminders(ctx context.Context, reminder time.Time) ([]*storage.Event, error)
	RemoveOlds(ctx context.Context, mark time.Time) error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(
	ctx context.Context,
	title string,
	description string,
	ownerID int64,
	startDate time.Time,
	endDate time.Time,
	remindIn time.Time,
) (int64, error) {
	event := storage.Event{
		ID:          0,
		Title:       title,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		OwnerID:     ownerID,
		RemindIn:    remindIn,
	}

	select {
	case <-ctx.Done():
		return 0, errors.New("create event ctx done")
	default:
		err := a.Storage.Add(ctx, &event)
		if err != nil {
			return 0, errors.Wrap(err, "create event")
		}

		return event.ID, nil
	}
}

func (a *App) EditEvent(
	ctx context.Context,
	id int64,
	title string,
	description string,
	ownerID int64,
	startDate time.Time,
	endDate time.Time,
	remindIn time.Time,
) error {
	ch := storage.Event{
		ID:          id,
		Title:       title,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		OwnerID:     ownerID,
		RemindIn:    remindIn,
	}

	select {
	case <-ctx.Done():
		return errors.New("edit event ctx done")
	default:
		event, err := a.Storage.Get(ctx, id)
		if err != nil {
			return errors.Wrap(err, "edit event get event")
		}

		if change, isChanged := ch.ExactChanges(event); isChanged {
			err := a.Storage.Edit(ctx, &change)
			if err != nil {
				return errors.Wrap(err, "edit event")
			}
		}

		return nil
	}
}

func (a *App) RemoveEvent(ctx context.Context, id int64) error {
	select {
	case <-ctx.Done():
		return errors.New("remove event ctx done")

	default:
		err := a.Storage.Delete(ctx, id)
		if err != nil {
			return errors.Wrap(err, "remove event")
		}

		return nil
	}
}

func (a *App) GetEvent(ctx context.Context, id int64) (storage.Event, error) {
	select {
	case <-ctx.Done():
		return storage.Event{}, errors.New("get event ctx done")
	default:
		event, err := a.Storage.Get(ctx, id)
		if err != nil {
			return storage.Event{}, errors.Wrap(err, "get event")
		}

		return event, nil
	}
}

func (a *App) GetDateTimeRangeEvents(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
) ([]storage.Event, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("get datetime range events ctx done")

	default:
		events, err := a.Storage.GetDateTimeRange(ctx, startDate, endDate)
		if err != nil {
			return nil, errors.Wrap(err, "get datetime range events")
		}

		return events, nil
	}
}

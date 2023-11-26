package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/rb"
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/pkg/errors"
)

type Scheduler struct {
	Storage  Storage
	Client   rb.Client
	Logger   *logger.Logger
	Period   int64
	Mark     time.Duration
	Exchange string
	Key      string
}

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	GetReminders(ctx context.Context, reminder time.Time) ([]*storage.Event, error)
	RemoveOlds(ctx context.Context, mark time.Time) error
}

func New(
	storage Storage,
	client rb.Client,
	logger *logger.Logger,
	period int64,
	mark int64,
	exchange string,
	key string,
) *Scheduler {
	return &Scheduler{
		Storage:  storage,
		Client:   client,
		Logger:   logger,
		Period:   period,
		Mark:     time.Duration(mark) * time.Second,
		Exchange: exchange,
		Key:      key,
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(s.Period) * time.Second)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return errors.New("scheduler context is done")
		case <-ticker.C:
			err := s.Storage.RemoveOlds(ctx, time.Now().Add(-s.Mark))
			if err != nil {
				return errors.Wrap(err, "scheduler fail to remove old events")
			}

			events, err := s.Storage.GetReminders(ctx, time.Now())
			if err != nil {
				return errors.Wrap(err, "scheduler fail to get reminders")
			}

			publishedAmount := 0

			for _, event := range events {
				message := rb.Message{
					ID:        event.ID,
					Title:     event.Title,
					OwnerID:   event.OwnerID,
					StartDate: event.StartDate,
					EndDate:   event.EndDate,
				}

				byteMessage, err := json.Marshal(message)
				if err != nil {
					s.Logger.Error(fmt.Sprintf("scheduler fail to marshal event %d", event.ID))
					continue
				}

				if err = s.Client.Publish(s.Exchange, s.Key, byteMessage); err != nil {
					return errors.Wrap(err, "scheduler fail to publish message")
				}

				publishedAmount++
			}

			if len(events) == 0 {
				s.Logger.Info("no event to notify")
			}

			s.Logger.Info(fmt.Sprintf("Amount of events to notify: %d", len(events)))
			s.Logger.Info(fmt.Sprintf("Amount of events was notified: %d", publishedAmount))
		}
	}
}

package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	app.Storage
	mu        sync.RWMutex
	eventsMap map[int64]*storage.Event
	lastID    int64
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(_ context.Context) error {
	s.eventsMap = make(map[int64]*storage.Event)

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.eventsMap = nil

	return nil
}

func (s *Storage) Add(_ context.Context, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = s.lastID
	s.lastID++
	s.eventsMap[event.ID] = event

	return nil
}

func (s *Storage) Edit(_ context.Context, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.eventsMap[event.ID]; !ok {
		return storage.ErrNotFound
	}

	s.eventsMap[event.ID] = event

	return nil
}

func (s *Storage) Delete(_ context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.eventsMap, id)

	return nil
}

func (s *Storage) Get(_ context.Context, id int64) (storage.Event, error) {
	event, ok := s.eventsMap[id]

	if !ok {
		return storage.Event{}, storage.ErrNotFound
	}

	return *event, nil
}

func (s *Storage) GetDateTimeRange(_ context.Context, from, to time.Time) ([]storage.Event, error) {
	result := make([]storage.Event, 0)

	for _, event := range s.eventsMap {
		if event.StartDate.After(from) && event.StartDate.Before(to) {
			result = append(result, *event)
		}
	}

	return result, nil
}

func (s *Storage) GetReminders(_ context.Context, reminder time.Time) ([]*storage.Event, error) {
	result := make([]*storage.Event, 0)

	for _, event := range s.eventsMap {
		if event.RemindIn.Equal(reminder) || event.RemindIn.After(reminder) && event.RemindIn.Before(time.Now()) {
			result = append(result, event)
		}
	}

	return result, nil
}

func (s *Storage) RemoveOlds(_ context.Context, mark time.Time) error {
	for _, event := range s.eventsMap {
		if event.EndDate.Before(mark) {
			delete(s.eventsMap, event.ID)
		}
	}

	return nil
}

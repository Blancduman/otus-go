package memorystorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	store := New()
	ctx := context.TODO()
	_ = store.Connect(ctx)

	defer func() {
		err := store.Close(ctx)
		if err != nil {
			fmt.Println("can not close memory storage")
		}
	}()

	for _, event := range []storage.Event{fixture1(), fixture2(), fixture3(), fixture4(), fixture5(), fixture6()} {
		event := event
		err := store.Add(&event)
		require.NoError(t, err)
	}

	t.Run("get date time range all", func(t *testing.T) {
		events, err := store.GetDateTimeRange(time.Now(), time.Now().Add(time.Hour*10))
		require.NoError(t, err)
		require.Len(t, events, 6)
	})

	t.Run("get date time range out", func(t *testing.T) {
		events, err := store.GetDateTimeRange(time.Now().Add(-time.Hour*100), time.Now().Add(-time.Hour*10))
		require.NoError(t, err)
		require.Len(t, events, 0)
	})

	t.Run("edit event", func(t *testing.T) {
		now := time.Now()
		newEvent := storage.Event{
			Title:       "edit",
			StartDate:   now.Add(time.Hour * 24),
			EndDate:     now.Add(time.Hour * 25),
			Description: "edit test",
		}

		err := store.Add(&newEvent)
		require.NoError(t, err)
		require.NotEqual(t, 0, newEvent.ID)

		updateEvent := storage.Event{
			ID:          newEvent.ID,
			Title:       "edit2",
			StartDate:   now.Add(time.Hour * 25),
			EndDate:     now.Add(time.Hour * 26),
			Description: "edit test2",
		}
		err = store.Edit(&updateEvent)
		require.NoError(t, err)

		updatedEvent, err := store.Get(newEvent.ID)
		require.NoError(t, err)
		require.Equal(t, "edit2", updatedEvent.Title)
		require.Equal(t, "edit test2", updatedEvent.Description)
		require.Equal(t, now.Add(time.Hour*25), updatedEvent.StartDate)
		require.Equal(t, now.Add(time.Hour*26), updatedEvent.EndDate)
	})

	t.Run("delete event", func(t *testing.T) {
		removableEvents, err := store.GetDateTimeRange(time.Now(), time.Now().Add(time.Hour*10))
		require.NoError(t, err)
		require.NotEqual(t, 0, len(removableEvents))

		for _, event := range removableEvents {
			err := store.Delete(event.ID)
			require.NoError(t, err)
		}

		events, err := store.GetDateTimeRange(time.Now(), time.Now().Add(time.Hour*10))
		require.NoError(t, err)
		require.Len(t, events, 0)
	})
}

func fixture1() storage.Event {
	return storage.Event{
		Title:       "Test Event 1",
		StartDate:   time.Now().Add(time.Hour),
		EndDate:     time.Now().Add(time.Hour * 2),
		Description: "This is test event 1",
	}
}

func fixture2() storage.Event {
	return storage.Event{
		Title:       "Test Event 2",
		StartDate:   time.Now().Add(time.Hour * 2),
		EndDate:     time.Now().Add(time.Hour * 3),
		Description: "This is test event 2",
	}
}

func fixture3() storage.Event {
	return storage.Event{
		Title:       "Test Event 3",
		StartDate:   time.Now().Add(time.Hour * 3),
		EndDate:     time.Now().Add(time.Hour * 4),
		Description: "This is test event 2",
	}
}

func fixture4() storage.Event {
	return storage.Event{
		Title:       "Test Event 4",
		StartDate:   time.Now().Add(time.Hour * 4),
		EndDate:     time.Now().Add(time.Hour * 5),
		Description: "This is test event 4",
	}
}

func fixture5() storage.Event {
	return storage.Event{
		Title:       "Test Event 5",
		StartDate:   time.Now().Add(time.Hour * 5),
		EndDate:     time.Now().Add(time.Hour * 6),
		Description: "This is test event 5",
	}
}

func fixture6() storage.Event {
	return storage.Event{
		Title:       "Test Event 1",
		StartDate:   time.Now().Add(time.Hour * 6),
		EndDate:     time.Now().Add(time.Hour * 7),
		Description: "This is test event 1",
	}
}

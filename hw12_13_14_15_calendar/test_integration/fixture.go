//go:build integration

package test_integration

import (
	"time"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage"
)

func Fixture1() storage.Event {
	return storage.Event{
		ID:          1,
		Title:       "Test 1",
		StartDate:   time.Now().Add(time.Hour),
		EndDate:     time.Now().Add(time.Hour * 2),
		Description: "Test 1 description",
		OwnerID:     1,
		RemindIn:    time.Now().Add(time.Minute * 30),
	}
}

func Fixture2() storage.Event {
	return storage.Event{
		ID:          2,
		Title:       "Test 2",
		StartDate:   time.Now().Add(time.Hour * 3),
		EndDate:     time.Now().Add(time.Hour * 4),
		Description: "Test 2 description",
		OwnerID:     1,
		RemindIn:    time.Now().Add(time.Minute*30 + time.Hour*2),
	}
}

func Fixture3() storage.Event {
	return storage.Event{
		ID:          3,
		Title:       "Test 3",
		StartDate:   time.Now().Add(time.Hour * 25),
		EndDate:     time.Now().Add(time.Hour * 26),
		Description: "Test 3 description",
		OwnerID:     2,
		RemindIn:    time.Now().Add(time.Minute*30 + time.Hour*24),
	}
}

func Fixture4() storage.Event {
	return storage.Event{
		ID:          4,
		Title:       "Test 4",
		StartDate:   time.Now().Add(time.Hour * 30),
		EndDate:     time.Now().Add(time.Hour * 31),
		Description: "Test 4 description",
		OwnerID:     2,
		RemindIn:    time.Now().Add(time.Minute*30 + time.Hour*29),
	}
}

func Fixture5() storage.Event {
	return storage.Event{
		ID:          5,
		Title:       "Test 5",
		StartDate:   time.Now().Add(time.Hour * 50),
		EndDate:     time.Now().Add(time.Hour * 51),
		Description: "Test 5 description",
		OwnerID:     3,
		RemindIn:    time.Now().Add(time.Minute*30 + time.Hour*49),
	}
}

func Fixture6() storage.Event {
	return storage.Event{
		ID:          6,
		Title:       "Test 6",
		StartDate:   time.Now().Add(time.Hour * 53),
		EndDate:     time.Now().Add(time.Hour * 54),
		Description: "Test 6 description",
		OwnerID:     3,
		RemindIn:    time.Now().Add(time.Minute*30 + time.Hour*52),
	}
}

func Fixture7() storage.Event {
	return storage.Event{
		ID:          7,
		Title:       "Test 7",
		StartDate:   time.Now().Add(time.Hour * 701),
		EndDate:     time.Now().Add(time.Hour * 702),
		Description: "Test 7 description",
		OwnerID:     4,
		RemindIn:    time.Now().Add(time.Minute*30 + time.Hour*700),
	}
}

func Fixture8() storage.Event {
	return storage.Event{
		ID:          8,
		Title:       "Test 8",
		StartDate:   time.Now().Add(time.Hour * 710),
		EndDate:     time.Now().Add(time.Hour * 711),
		Description: "Test 8 description",
		OwnerID:     4,
		RemindIn:    time.Now().Add(time.Minute*30 + time.Hour*709),
	}
}

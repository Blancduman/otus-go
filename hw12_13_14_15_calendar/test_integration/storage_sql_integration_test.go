//go:build integration

package test_integration

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
)

const dsn = "host=db port=5432 user=otus_user password=otus_password dbname=calendar_test sslmode=disable TimeZone=UTC"

type PostgresDBRepoTestSuite struct {
	suite.Suite
	store *sqlstorage.Storage
}

func TestRepoPostgresDBTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresDBRepoTestSuite))
}

func (s *PostgresDBRepoTestSuite) createDB(ctx context.Context) {
	conn, err := pgx.Connect(ctx, dsn)
	s.Require().NoError(err)

	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS event (
		   id serial primary key,
		   title text,
		   start_date timestamp with time zone not null,
		   end_date timestamp with time zone not null,
		   description text,
		   owner_id bigint,
		   remind_in timestamp with time zone
		)`,
	)
	s.Require().NoError(err)

	err = conn.Close(ctx)
	s.Require().NoError(err)
}

func (s *PostgresDBRepoTestSuite) SetupTest() {
	ctx := context.TODO()

	s.createDB(ctx)

	s.store = sqlstorage.New(dsn)
	err := s.store.Connect(ctx)
	s.Require().NoError(err)

	fixtures := []storage.Event{
		Fixture1(),
		Fixture2(),
		Fixture3(),
		Fixture4(),
		Fixture5(),
		Fixture6(),
		Fixture7(),
		Fixture8(),
	}

	for _, event := range fixtures {
		err := s.store.Add(ctx, &event)
		s.Require().NoError(err)
	}
}

func (s *PostgresDBRepoTestSuite) TearDownTest() {
	ctx := context.TODO()

	conn, err := pgx.Connect(ctx, dsn)
	s.Require().NoError(err)

	_, err = conn.Exec(ctx, `DROP TABLE event`)
	s.Require().NoError(err)

	err = conn.Close(ctx)
	s.Require().NoError(err)
}

func (s *PostgresDBRepoTestSuite) TearDownSuite() {
	ctx := context.TODO()

	err := s.store.Close(ctx)
	s.Require().NoError(err)
}

func (s *PostgresDBRepoTestSuite) TestAddSuccess() {
	ctx := context.TODO()

	event := storage.Event{
		ID:          0,
		Title:       "Test 9",
		StartDate:   time.Now().Add(time.Hour * 6),
		EndDate:     time.Now().Add(time.Hour * 7),
		Description: "Test 9 description",
		OwnerID:     5,
		RemindIn:    time.Now().Add(time.Hour*5 + time.Minute*30),
	}

	err := s.store.Add(ctx, &event)
	s.Require().NoError(err)

	ev, err := s.store.Get(ctx, event.ID)
	s.Require().NoError(err)
	s.Require().Equal(int64(9), ev.ID)
}

func (s *PostgresDBRepoTestSuite) TestAddFail() {
	ctx := context.TODO()

	var event *storage.Event

	err := s.store.Add(ctx, event)
	s.Require().Error(err)
}

func (s *PostgresDBRepoTestSuite) TestGetDay() {
	ctx := context.TODO()

	events, err := s.store.GetDateTimeRange(ctx, time.Now(), time.Now().Add(time.Hour*24))
	s.Require().NoError(err)
	s.Require().Len(events, 2)
}

func (s *PostgresDBRepoTestSuite) TestGetWeek() {
	ctx := context.TODO()

	events, err := s.store.GetDateTimeRange(ctx, time.Now(), time.Now().Add(time.Hour*24*7))
	s.Require().NoError(err)
	s.Require().Len(events, 6)
}

func (s *PostgresDBRepoTestSuite) TestGetMonth() {
	ctx := context.TODO()

	events, err := s.store.GetDateTimeRange(ctx, time.Now(), time.Now().Add(time.Hour*730))
	s.Require().NoError(err)
	s.Require().Len(events, 8)
}

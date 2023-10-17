package sqlstorage

import (
	"context"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type Storage struct { // TODO
	app.Storage
	dsn  string
	conn *pgx.Conn
}

func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, s.dsn)
	if err != nil {
		return errors.Wrap(err, "fail to connect db")
	}

	if err = conn.Ping(ctx); err != nil {
		return errors.Wrap(err, "fail to ping connection db")
	}

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if s.conn == nil {
		return nil
	}

	err := s.conn.Close(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to close db connection")
	}

	return nil
}

func (s *Storage) Add(event *storage.Event) error {
	var id int

	query := `INSERT INTO event(title, start_date, end_date, description) VALUES ($1, $2, $3, $4) RETURNING id`

	err := s.conn.QueryRow(
		context.Background(),
		query,
		event.Title,
		event.StartDate,
		event.EndDate,
		event.Description,
	).Scan(&id)
	if err != nil {
		return errors.Wrap(err, "storage fail to insert event")
	}

	event.ID = id

	return nil
}

func (s *Storage) Edit(event *storage.Event) error {
	query := `UPDATE event SET title=$1, start_date=$2, end_date=$3, description=$4 WHERE id=$5`

	_, err := s.conn.Exec(
		context.Background(),
		query,
		event.Title,
		event.StartDate,
		event.EndDate,
		event.Description,
		event.ID,
	)
	if err != nil {
		return errors.Wrap(err, "storage fail to set event")
	}

	return nil
}

func (s *Storage) Delete(id int) error {
	query := `DELETE FROM event WHERE id=$1`

	_, err := s.conn.Exec(context.Background(), query, id)
	if err != nil {
		return errors.Wrap(err, "storage fail to delete")
	}

	return nil
}

func (s *Storage) Get(id int) (storage.Event, error) {
	event := storage.Event{}
	query := `SELECT id, title, start_date, end_date, description FROM event WHERE id=$1`

	row := s.conn.QueryRow(context.Background(), query, id)

	err := row.Scan(&event.ID, &event.Title, &event.StartDate, &event.EndDate, &event.Description)
	if err != nil {
		return event, errors.Wrap(err, "storage fail to get event")
	}

	return event, nil
}

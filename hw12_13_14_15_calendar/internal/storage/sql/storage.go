package sqlstorage

import (
	"context"
	"strings"
	"time"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type Storage struct {
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

	println("pinged")

	s.conn = conn

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

func (s *Storage) Add(ctx context.Context, event *storage.Event) error {
	var id int64

	if event == nil {
		return errors.New("empty event")
	}

	query := `
		INSERT INTO event(title, start_date, end_date, description, owner_id, remind_in)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := s.conn.QueryRow(
		ctx,
		query,
		event.Title,
		event.StartDate.UTC(),
		event.EndDate.UTC(),
		event.Description,
		event.OwnerID,
		event.RemindIn.UTC(),
	).Scan(&id)
	if err != nil {
		return errors.Wrap(err, "storage fail to insert event")
	}

	event.ID = id

	return nil
}

func (s *Storage) Edit(ctx context.Context, event *storage.Event) error {
	query := `
		UPDATE event
		SET title=$1, start_date=$2, end_date=$3, description=$4, owner_id=$5, remind_in=$6
		WHERE id=$7
	`

	_, err := s.conn.Exec(
		ctx,
		query,
		event.Title,
		event.StartDate.UTC(),
		event.EndDate.UTC(),
		event.Description,
		event.OwnerID,
		event.RemindIn.UTC(),
		event.ID,
	)
	if err != nil {
		return errors.Wrap(err, "storage fail to set event")
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM event WHERE id=$1`

	_, err := s.conn.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "storage fail to delete")
	}

	return nil
}

func (s *Storage) Get(ctx context.Context, id int64) (storage.Event, error) {
	event := storage.Event{}
	query := `SELECT id, title, start_date, end_date, description, remind_in, owner_id FROM event WHERE id=$1`

	row := s.conn.QueryRow(ctx, query, id)

	err := row.Scan(
		&event.ID,
		&event.Title,
		&event.StartDate,
		&event.EndDate,
		&event.Description,
		&event.RemindIn,
		&event.OwnerID,
	)
	if err != nil {
		return event, errors.Wrap(err, "storage fail to get event")
	}

	return event, nil
}

func (s *Storage) GetDateTimeRange(ctx context.Context, from, to time.Time) ([]storage.Event, error) {
	query := `
		SELECT id, title, description, start_date, end_date, remind_in, owner_id
		FROM event
		WHERE start_date BETWEEN $1 AND $2
	`

	rows, err := s.conn.Query(ctx, query, from.UTC(), to.UTC())
	if err != nil {
		return nil, errors.Wrap(err, "storage fail to get datetime range")
	}

	defer rows.Close()

	var result []storage.Event

	for rows.Next() {
		event := storage.Event{}

		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDate,
			&event.EndDate,
			&event.RemindIn,
			&event.OwnerID,
		); err != nil {
			return nil, errors.Wrap(err, "storage fail to read row")
		}

		result = append(result, event)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "storage row read error")
	}

	return result, nil
}

func (s *Storage) GetReminders(ctx context.Context, reminder time.Time) ([]*storage.Event, error) {
	query := `SELECT id, title, description, start_date, end_date, remind_in, owner_id
		FROM event
		WHERE remind_in BETWEEN $1 AND $2
	`

	rows, err := s.conn.Query(
		ctx,
		query,
		strings.Trim(reminder.UTC().String(), " UTC"),
		strings.Trim(time.Now().UTC().String(), " UTC"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "storage fail to get datetime range")
	}

	defer rows.Close()

	var result []*storage.Event

	for rows.Next() {
		event := storage.Event{}

		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDate,
			&event.EndDate,
			&event.RemindIn,
			&event.OwnerID,
		); err != nil {
			return nil, errors.Wrap(err, "storage fail to read row")
		}

		result = append(result, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "storage row read error")
	}

	return result, nil
}

func (s *Storage) RemoveOlds(ctx context.Context, mark time.Time) error {
	query := `DELETE FROM event WHERE end_date < $1`

	_, err := s.conn.Exec(ctx, query, mark.UTC())
	if err != nil {
		return errors.Wrap(err, "storage fail to delete old events")
	}

	return nil
}

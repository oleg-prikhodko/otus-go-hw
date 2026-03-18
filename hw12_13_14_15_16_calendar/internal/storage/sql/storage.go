package sqlstorage

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/common"
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

var ErrNotConnected = errors.New("database not connected")

type Storage struct {
	db     *sqlx.DB
	logger common.Logger
	addr   string
}

func New(logger common.Logger, addr string) *Storage {
	return &Storage{logger: logger, addr: addr}
}

func (s *Storage) Connect() error {
	db, err := sqlx.Connect("postgres", s.addr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	s.db = db
	return nil
}

func (s *Storage) Close() error {
	if s.db == nil {
		return nil
	}

	return s.db.Close()
}

func (s *Storage) Create(ev storage.Event) error {
	if s.db == nil {
		return ErrNotConnected
	}

	query := `
		INSERT INTO events (title, event_time, duration, description, owner_id, notify_before)
		VALUES (:title, :event_time, :duration, :description, :owner_id, :notify_before)
		RETURNING id;`

	_, err := s.db.NamedExec(query, ev)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

func (s *Storage) Update(ev storage.Event) error {
	if s.db == nil {
		return ErrNotConnected
	}

	query := `
		UPDATE events
		SET title = $1, event_time = $2, duration = $3, description = $4, owner_id = $5, notify_before = $6
		WHERE id = $7
		RETURNING *;`

	res, err := s.db.NamedExec(query, ev)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	if affected == 0 {
		return common.NotFoundErr
	}

	return nil
}

func (s *Storage) Delete(id string) error {
	if s.db == nil {
		return ErrNotConnected
	}

	res, err := s.db.Exec("DELETE FROM events WHERE id = $1;", id)
	if err != nil {
		return fmt.Errorf("failed to delete event %s: %w", id, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to delete event %s: %w", id, err)
	}
	if affected == 0 {
		return common.NotFoundErr
	}

	return nil
}

func (s *Storage) List(from, to time.Time) ([]storage.Event, error) {
	if s.db == nil {
		return nil, ErrNotConnected
	}

	query := `
		SELECT id, title, event_time, duration, description, owner_id, notify_before
		FROM events
		WHERE event_time >= $1 AND event_time <= $2;`

	events := make([]storage.Event, 0)
	err := s.db.Select(&events, query, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	return events, nil
}

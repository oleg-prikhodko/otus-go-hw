package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	internalhttp "github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

var ErrNotConnected = errors.New("database not connected")

type Storage struct {
	db     *sqlx.DB
	logger internalhttp.Logger
}

func New(logger internalhttp.Logger) *Storage {
	return &Storage{logger: logger}
}

func (s *Storage) Connect(ctx context.Context, connStr string) error {
	db, err := sqlx.ConnectContext(ctx, "postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	s.db = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if s.db == nil {
		return nil
	}

	return s.db.Close()
}

// repository impl

func (s *Storage) Create(ctx context.Context, ev *storage.Event) error {
	if s.db == nil {
		return ErrNotConnected
	}

	query := `
		INSERT INTO events (title, event_time, duration, description, owner_id, notify_before)
		VALUES (:title, :event_time, :duration, :description, :owner_id, :notify_before)
		RETURNING id;`

	stmt, err := s.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &ev.ID, ev)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

func (s *Storage) Update(ctx context.Context, id string, ev *storage.Event) error {
	if s.db == nil {
		return ErrNotConnected
	}

	query := `
		UPDATE events
		SET title = $1, event_time = $2, duration = $3, description = $4, owner_id = $5, notify_before = $6
		WHERE id = $7
		RETURNING *;`

	var updated storage.Event
	err := s.db.GetContext(ctx, &updated, query,
		ev.Title, ev.Time, ev.Duration, ev.Description, ev.OwnerID, ev.NotifyBefore, id)

	if err != nil {
		return fmt.Errorf("failed to update event %s: %w", id, err)
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	if s.db == nil {
		return ErrNotConnected
	}

	q := `DELETE FROM events WHERE id = $1;`
	_, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("failed to delete event %s: %w", id, err)
	}

	return nil
}

func (s *Storage) List(ctx context.Context, from, to time.Time) ([]*storage.Event, error) {
	if s.db == nil {
		return nil, ErrNotConnected
	}

	query := `
		SELECT id, title, event_time, duration, description, owner_id, notify_before
		FROM events
		WHERE event_time >= $1 AND event_time <= $2;`

	events := make([]*storage.Event, 0)
	err := s.db.SelectContext(ctx, &events, query, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	return events, nil
}

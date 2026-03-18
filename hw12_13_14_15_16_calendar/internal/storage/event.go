package storage

import (
	"context"
	"time"
)

type Event struct {
	ID           string         `db:"id"`
	Title        string         `db:"title"`
	Time         time.Time      `db:"event_time"`
	Duration     time.Duration  `db:"duration"`
	Description  *string        `db:"description"`
	OwnerID      string         `db:"owner_id"`
	NotifyBefore *time.Duration `db:"notify_before"`
}

type EventStorage interface {
	Create(ctx context.Context, ev *Event) error
	Update(ctx context.Context, id string, ev *Event) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, from time.Time, to time.Time) ([]*Event, error)
}

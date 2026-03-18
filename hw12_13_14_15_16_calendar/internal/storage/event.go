package storage

import (
	"io"
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
	io.Closer
	Create(ev Event) error
	Update(ev Event) error
	Delete(id string) error
	List(from time.Time, to time.Time) ([]Event, error)
}

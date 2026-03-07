package storage

import "time"

type Event struct {
	ID           string
	Title        string
	Time         time.Time
	Duration     time.Duration
	Description  *string
	OwnerID      string
	NotifyBefore *time.Duration
}

type EventRepository interface {
	Create(ev *Event) (*Event, error)
	Update(id string, ev *Event) (*Event, error)
	Delete(id string) (*Event, error)
	List(from time.Time, to time.Time) ([]*Event, error)
}

package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events sync.Map
}

func (s *Storage) Create(ctx context.Context, ev *storage.Event) error {
	ev.ID = uuid.New().String()
	s.events.Store(ev.ID, ev)

	return nil
}

func (s *Storage) Update(ctx context.Context, id string, ev *storage.Event) error {
	_, ok := s.events.Load(id)
	if !ok {
		return errors.New("not found")
	}

	ev.ID = id
	s.events.Store(id, ev)

	return nil
}

func (s *Storage) Delete(ctx context.Context, id string) error {
	_, ok := s.events.LoadAndDelete(id)
	if !ok {
		return errors.New("not found")
	}

	return nil
}

func (s *Storage) List(ctx context.Context, from time.Time, to time.Time) ([]*storage.Event, error) {
	var events []*storage.Event

	s.events.Range(func(k, v interface{}) bool {
		ev := v.(*storage.Event)
		if ev.Time.After(from) && ev.Time.Before(to) {
			// Return a copy to prevent external mutation
			events = append(events, copyEvent(ev))
		}

		return true
	})

	return events, nil
}

func copyEvent(ev *storage.Event) *storage.Event {
	copy := &storage.Event{
		ID:       ev.ID,
		Title:    ev.Title,
		Time:     ev.Time,
		Duration: ev.Duration,
		OwnerID:  ev.OwnerID,
	}
	if ev.Description != nil {
		desc := *ev.Description
		copy.Description = &desc
	}
	if ev.NotifyBefore != nil {
		notify := *ev.NotifyBefore
		copy.NotifyBefore = &notify
	}
	return copy
}

func New() *Storage {
	return &Storage{}
}

package memorystorage

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events sync.Map
	mu     sync.RWMutex //nolint:unused
}

func (s *Storage) Create(ev *storage.Event) (*storage.Event, error) {
	//s.mu.Lock()
	//defer s.mu.Unlock()

	e := &storage.Event{
		ID:           uuid.New().String(),
		Title:        ev.Title,
		Time:         ev.Time,
		Duration:     ev.Duration,
		Description:  ev.Description,
		OwnerID:      ev.OwnerID,
		NotifyBefore: ev.NotifyBefore,
	}

	s.events.Store(ev.ID, e)

	return e, nil
}

func (s *Storage) Update(id string, ev *storage.Event) (*storage.Event, error) {
	_, ok := s.events.Load(id)
	if !ok {
		return nil, errors.New("not found")
	}

	s.events.Store(id, ev)

	return ev, nil
}

func (s *Storage) Delete(id string) (*storage.Event, error) {
	e, ok := s.events.Load(id)
	if !ok {
		return nil, errors.New("not found")
	}

	s.events.Delete(id)

	return e.(*storage.Event), nil
}

func (s *Storage) List(from time.Time, to time.Time) ([]*storage.Event, error) {
	var events []*storage.Event

	s.events.Range(func(k, v interface{}) bool {
		ev := v.(*storage.Event)
		start := ev.Time
		end := ev.Time.Add(ev.Duration)
		if start.Before(to) && end.After(from) {
			events = append(events, ev)
		}

		return true
	})

	return events, nil
}

func New() *Storage {
	return &Storage{}
}

package memorystorage

import (
	"sync"
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/common"
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events sync.Map
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) Create(ev storage.Event) error {
	s.events.Store(ev.ID, ev)

	return nil
}

func (s *Storage) Update(ev storage.Event) error {
	_, ok := s.events.Load(ev.ID)
	if !ok {
		return common.NotFoundErr
	}

	s.events.Store(ev.ID, ev)

	return nil
}

func (s *Storage) Delete(id string) error {
	_, ok := s.events.LoadAndDelete(id)
	if !ok {
		return common.NotFoundErr
	}

	return nil
}

func (s *Storage) List(from time.Time, to time.Time) ([]storage.Event, error) {
	var events []storage.Event

	s.events.Range(func(k, v any) bool {
		ev := v.(storage.Event)
		if ev.Time.After(from) && ev.Time.Before(to) {
			events = append(events, ev)
		}

		return true
	})

	return events, nil
}

func New() *Storage {
	return &Storage{events: sync.Map{}}
}

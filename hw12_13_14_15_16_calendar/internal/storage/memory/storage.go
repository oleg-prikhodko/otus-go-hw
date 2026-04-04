package memorystorage

import (
	"sync"
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/common"  //nolint:depguard
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
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
		return common.ErrNotFound
	}

	s.events.Store(ev.ID, ev)

	return nil
}

func (s *Storage) Delete(id string) error {
	_, ok := s.events.LoadAndDelete(id)
	if !ok {
		return common.ErrNotFound
	}

	return nil
}

func (s *Storage) List(from time.Time, to time.Time) ([]storage.Event, error) {
	var events []storage.Event

	s.events.Range(func(_, v any) bool {
		ev := v.(storage.Event)
		if ev.Time.After(from) && ev.Time.Before(to) {
			events = append(events, ev)
		}

		return true
	})

	return events, nil
}

func (s *Storage) ListForNotification(from time.Time, to time.Time) ([]storage.Event, error) {
	var events []storage.Event

	s.events.Range(func(_, v any) bool {
		ev := v.(storage.Event)
		var before time.Duration
		if ev.NotifyBefore != nil {
			before = *ev.NotifyBefore
		}
		if ev.Time.Add(-before).After(from) && ev.Time.Add(-before).Before(to) {
			events = append(events, ev)
		}

		return true
	})

	return events, nil
}

func (s *Storage) DeleteOutdated(olderThan time.Time) error {
	s.events.Range(func(k, v any) bool {
		ev := v.(storage.Event)
		if ev.Time.Before(olderThan) {
			s.events.Delete(k)
		}
		return true
	})
	return nil
}

func New() *Storage {
	return &Storage{events: sync.Map{}}
}

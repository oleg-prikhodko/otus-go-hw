package common

import (
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

type MockLogger struct {
	InfoMsg string
	ErrMsg  string
}

func (m *MockLogger) Info(msg string) {
	m.InfoMsg = msg
}

func (m *MockLogger) Error(msg string) {
	m.ErrMsg = msg
}

type MockApp struct {
	EventCreateCalled bool
	EventCreate       storage.Event
	EventCreateErr    error

	EventUpdateCalled bool
	EventUpdate       storage.Event
	EventUpdateErr    error

	EventDeleteCalled bool
	EventDeleteID     string
	EventDeleteErr    error

	ListDayEvents []storage.Event
	ListDayErr    error

	ListWeekEvents []storage.Event
	ListWeekErr    error

	ListMonthEvents []storage.Event
	ListMonthErr    error
}

func (m *MockApp) CreateEvent(ev storage.Event) error {
	m.EventCreateCalled = true
	m.EventCreate = ev
	return m.EventCreateErr
}

func (m *MockApp) UpdateEvent(ev storage.Event) error {
	m.EventUpdateCalled = true
	m.EventUpdate = ev
	return m.EventUpdateErr
}

func (m *MockApp) DeleteEvent(id string) error {
	m.EventDeleteCalled = true
	m.EventDeleteID = id
	return m.EventDeleteErr
}

func (m *MockApp) ListEventsForDay(date time.Time) ([]storage.Event, error) {
	return m.ListDayEvents, m.ListDayErr
}

func (m *MockApp) ListEventsForWeek(date time.Time) ([]storage.Event, error) {
	return m.ListWeekEvents, m.ListWeekErr
}

func (m *MockApp) ListEventsForMonth(date time.Time) ([]storage.Event, error) {
	return m.ListMonthEvents, m.ListMonthErr
}

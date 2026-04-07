package internalhttp

import (
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type CreateEventRequest struct {
	Title        string    `json:"title"`
	Time         time.Time `json:"time"`
	Duration     int64     `json:"duration"`
	Description  *string   `json:"description,omitempty"`
	OwnerID      string    `json:"ownerId"`
	NotifyBefore *int64    `json:"notifyBefore,omitempty"`
}

type UpdateEventRequest struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Time         time.Time `json:"time"`
	Duration     int64     `json:"duration"`
	Description  *string   `json:"description,omitempty"`
	OwnerID      string    `json:"ownerId"`
	NotifyBefore *int64    `json:"notifyBefore,omitempty"`
}

type DeleteEventRequest struct {
	ID string `json:"id"`
}

type ListEventsRequest struct {
	Date time.Time `json:"date"`
}

type EventResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Time         time.Time `json:"time"`
	Duration     int64     `json:"duration"`
	Description  *string   `json:"description,omitempty"`
	OwnerID      string    `json:"ownerId"`
	NotifyBefore *int64    `json:"notifyBefore,omitempty"`
}

func EventToResponse(ev storage.Event) EventResponse {
	var notifyBefore *int64
	if ev.NotifyBefore != nil {
		nano := ev.NotifyBefore.Nanoseconds()
		notifyBefore = &nano
	}
	return EventResponse{
		ID:           ev.ID,
		Title:        ev.Title,
		Time:         ev.Time,
		Duration:     ev.Duration.Nanoseconds(),
		Description:  ev.Description,
		OwnerID:      ev.OwnerID,
		NotifyBefore: notifyBefore,
	}
}

type ListEventsResponse struct {
	Events []EventResponse `json:"events"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateEventResponse struct {
	ID string `json:"id"`
}

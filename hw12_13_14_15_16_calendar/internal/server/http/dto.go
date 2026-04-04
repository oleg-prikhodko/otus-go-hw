package internalhttp

import (
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type CreateEventRequest struct {
	Title        string         `json:"title"`
	Time         time.Time      `json:"time"`
	Duration     time.Duration  `json:"duration"`
	Description  *string        `json:"description,omitempty"`
	OwnerID      string         `json:"ownerId"`
	NotifyBefore *time.Duration `json:"notifyBefore,omitempty"`
}

type UpdateEventRequest struct {
	ID           string         `json:"id"`
	Title        string         `json:"title"`
	Time         time.Time      `json:"time"`
	Duration     time.Duration  `json:"duration"`
	Description  *string        `json:"description,omitempty"`
	OwnerID      string         `json:"ownerId"`
	NotifyBefore *time.Duration `json:"notifyBefore,omitempty"`
}

type DeleteEventRequest struct {
	ID string `json:"id"`
}

type ListEventsRequest struct {
	Date time.Time `json:"date"`
}

type EventResponse struct {
	ID           string         `json:"id"`
	Title        string         `json:"title"`
	Time         time.Time      `json:"time"`
	Duration     time.Duration  `json:"duration"`
	Description  *string        `json:"description,omitempty"`
	OwnerID      string         `json:"ownerId"`
	NotifyBefore *time.Duration `json:"notifyBefore,omitempty"`
}

func EventToResponse(ev storage.Event) EventResponse {
	return EventResponse{
		ID:           ev.ID,
		Title:        ev.Title,
		Time:         ev.Time,
		Duration:     ev.Duration,
		Description:  ev.Description,
		OwnerID:      ev.OwnerID,
		NotifyBefore: ev.NotifyBefore,
	}
}

type ListEventsResponse struct {
	Events []EventResponse `json:"events"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

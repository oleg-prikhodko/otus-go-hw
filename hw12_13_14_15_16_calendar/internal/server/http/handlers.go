package internalhttp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/common"  //nolint:depguard
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

func handleCreateEvent(app common.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateEventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		ev := storage.Event{
			Title:        req.Title,
			Time:         req.Time,
			Duration:     req.Duration,
			Description:  req.Description,
			OwnerID:      req.OwnerID,
			NotifyBefore: req.NotifyBefore,
		}

		if err := app.CreateEvent(ev); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func handleUpdateEvent(app common.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		var req UpdateEventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		ev := storage.Event{
			ID:           id,
			Title:        req.Title,
			Time:         req.Time,
			Duration:     req.Duration,
			Description:  req.Description,
			OwnerID:      req.OwnerID,
			NotifyBefore: req.NotifyBefore,
		}

		if err := app.UpdateEvent(ev); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handleDeleteEvent(app common.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		if err := app.DeleteEvent(id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func handleListDayEvents(app common.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dateStr := r.URL.Query().Get("date")
		if dateStr == "" {
			writeError(w, http.StatusBadRequest, "date parameter is required")
			return
		}

		date, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid date format, use RFC3339")
			return
		}

		events, err := app.ListEventsForDay(date)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, ListEventsResponse{Events: eventsToResponse(events)})
	}
}

func handleListWeekEvents(app common.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dateStr := r.URL.Query().Get("date")
		if dateStr == "" {
			writeError(w, http.StatusBadRequest, "date parameter is required")
			return
		}

		date, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid date format, use RFC3339")
			return
		}

		events, err := app.ListEventsForWeek(date)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, ListEventsResponse{Events: eventsToResponse(events)})
	}
}

func handleListMonthEvents(app common.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dateStr := r.URL.Query().Get("date")
		if dateStr == "" {
			writeError(w, http.StatusBadRequest, "date parameter is required")
			return
		}

		date, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid date format, use RFC3339")
			return
		}

		events, err := app.ListEventsForMonth(date)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, ListEventsResponse{Events: eventsToResponse(events)})
	}
}

func eventsToResponse(events []storage.Event) []EventResponse {
	result := make([]EventResponse, len(events))
	for i, ev := range events {
		result[i] = EventToResponse(ev)
	}
	return result
}

func writeError(w http.ResponseWriter, code int, message string) {
	writeJSON(w, code, ErrorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

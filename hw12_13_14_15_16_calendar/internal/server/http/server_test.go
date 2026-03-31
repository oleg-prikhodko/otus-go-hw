package internalhttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/common"
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

func TestHandleCreateEvent(t *testing.T) {
	app := &common.MockApp{}

	handler := handleCreateEvent(app)

	body := `{"title":"Test Event","time":"2024-01-01T10:00:00Z","duration":3600000000000,"ownerId":"user1"}`
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
	if !app.EventCreateCalled {
		t.Error("CreateEvent was not called")
	}
	if app.EventCreate.Title != "Test Event" {
		t.Errorf("expected title 'Test Event', got '%s'", app.EventCreate.Title)
	}
}

func TestHandleUpdateEvent(t *testing.T) {
	app := &common.MockApp{}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /events/{id}", handleUpdateEvent(app))

	body := `{"title":"Updated Event","time":"2024-01-01T10:00:00Z","duration":7200000000000,"ownerId":"user1"}`
	req := httptest.NewRequest(http.MethodPut, "/events/123", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if !app.EventUpdateCalled {
		t.Error("UpdateEvent was not called")
	}
	if app.EventUpdate.ID != "123" {
		t.Errorf("expected id '123', got '%s'", app.EventUpdate.ID)
	}
}

func TestHandleDeleteEvent(t *testing.T) {
	app := &common.MockApp{}

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /events/{id}", handleDeleteEvent(app))

	req := httptest.NewRequest(http.MethodDelete, "/events/123", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
	if !app.EventDeleteCalled {
		t.Error("DeleteEvent was not called")
	}
	if app.EventDeleteID != "123" {
		t.Errorf("expected id '123', got '%s'", app.EventDeleteID)
	}
}

func TestHandleListDayEvents(t *testing.T) {
	app := &common.MockApp{
		ListDayEvents: []storage.Event{
			{ID: "1", Title: "Event 1"},
			{ID: "2", Title: "Event 2"},
		},
	}

	handler := handleListDayEvents(app)

	req := httptest.NewRequest(http.MethodGet, "/events/day?date=2024-01-01T10:00:00Z", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp struct {
		Events []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"events"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Events) != 2 {
		t.Errorf("expected 2 events, got %d", len(resp.Events))
	}
}

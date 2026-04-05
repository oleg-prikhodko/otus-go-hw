package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid" //nolint:depguard
	_ "github.com/lib/pq"    //nolint:depguard
)

var (
	db          *sql.DB
	calendarURL string
)

func TestMain(m *testing.M) {
	dbConn := os.Getenv("DB_CONN")
	if dbConn == "" {
		panic("DB_CONN env var is not set")
	}

	var err error
	db, err = sql.Open("postgres", dbConn)
	if err != nil {
		panic(fmt.Sprintf("failed to open db: %v", err))
	}

	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("failed to ping db: %v", err))
	}

	calendarHost := os.Getenv("CALENDAR_HOST")
	if calendarHost == "" {
		panic("CALENDAR_HOST env var is not set")
	}
	calendarPort := os.Getenv("CALENDAR_HTTP_PORT")
	if calendarPort == "" {
		panic("CALENDAR_HTTP_PORT env var is not set")
	}
	calendarURL = fmt.Sprintf("http://%s:%s", calendarHost, calendarPort)

	time.Sleep(2 * time.Second)

	code := m.Run()
	_ = db.Close()
	os.Exit(code)
}

func truncateEvents(t *testing.T) {
	t.Helper()
	_, err := db.Exec("TRUNCATE TABLE events")
	if err != nil {
		t.Fatalf("failed to truncate events: %v", err)
	}
}

func TestCreateEvent_Success(t *testing.T) {
	t.Cleanup(func() { truncateEvents(t) })

	body := map[string]any{
		"title":    "Test Event",
		"time":     "2025-01-01T10:00:00Z",
		"duration": time.Hour * 3,
		"ownerId":  uuid.New().String(),
	}
	bodyBytes, _ := json.Marshal(body)

	resp, err := http.Post(calendarURL+"/events", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("failed to POST /events: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d: %s", http.StatusCreated, resp.StatusCode, string(bodyBytes))
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM events WHERE title = 'Test Event'").Scan(&count)
	if err != nil {
		t.Fatalf("failed to count events: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 event in db, got %d", count)
	}
}

func createTestEvent(t *testing.T, title string, eventTime time.Time) string {
	t.Helper()
	body := map[string]any{
		"title":    title,
		"time":     eventTime.Format(time.RFC3339),
		"duration": int64(time.Hour),
		"ownerId":  uuid.New().String(),
	}
	bodyBytes, _ := json.Marshal(body)

	resp, err := http.Post(calendarURL+"/events", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("failed to create test event: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create test event: status %d", resp.StatusCode)
	}

	var result struct {
		ID string `json:"id"`
	}
	respBody, _ := io.ReadAll(resp.Body)
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
	}

	if result.ID != "" {
		return result.ID
	}

	var id string
	err = db.QueryRow("SELECT id FROM events WHERE title = $1", title).Scan(&id)
	if err != nil {
		t.Fatalf("failed to get event id from db: %v", err)
	}
	return id
}

func TestListDayEvents_Success(t *testing.T) {
	t.Cleanup(func() { truncateEvents(t) })

	eventTime := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	createTestEvent(t, "Day Event", eventTime)

	resp, err := http.Get(calendarURL + "/events/day?date=2025-01-01T10:00:00Z")
	if err != nil {
		t.Fatalf("failed to GET /events/day: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var result struct {
		Events []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"events"`
	}
	respBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(result.Events) != 1 {
		t.Errorf("expected 1 event, got %d", len(result.Events))
	}
	if result.Events[0].Title != "Day Event" {
		t.Errorf("expected title 'Day Event', got '%s'", result.Events[0].Title)
	}
}

func TestListWeekEvents_Success(t *testing.T) {
	t.Cleanup(func() { truncateEvents(t) })

	eventTime := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	createTestEvent(t, "Week Event", eventTime)

	resp, err := http.Get(calendarURL + "/events/week?date=2025-01-01T10:00:00Z")
	if err != nil {
		t.Fatalf("failed to GET /events/week: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var result struct {
		Events []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"events"`
	}
	respBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(result.Events) != 1 {
		t.Errorf("expected 1 event, got %d", len(result.Events))
	}
	if result.Events[0].Title != "Week Event" {
		t.Errorf("expected title 'Week Event', got '%s'", result.Events[0].Title)
	}
}

func TestListMonthEvents_Success(t *testing.T) {
	t.Cleanup(func() { truncateEvents(t) })

	eventTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	createTestEvent(t, "Month Event", eventTime)

	resp, err := http.Get(calendarURL + "/events/month?date=2025-01-15T10:00:00Z")
	if err != nil {
		t.Fatalf("failed to GET /events/month: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var result struct {
		Events []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"events"`
	}
	respBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(result.Events) != 1 {
		t.Errorf("expected 1 event, got %d", len(result.Events))
	}
	if result.Events[0].Title != "Month Event" {
		t.Errorf("expected title 'Month Event', got '%s'", result.Events[0].Title)
	}
}

func TestDeleteEvent_Success(t *testing.T) {
	t.Cleanup(func() { truncateEvents(t) })

	eventTime := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	eventID := createTestEvent(t, "To Delete", eventTime)

	req, err := http.NewRequest(http.MethodDelete, calendarURL+"/events/"+eventID, nil)
	if err != nil {
		t.Fatalf("failed to create DELETE request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to DELETE /events: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Errorf("expected status %d, got %d: %s", http.StatusNoContent, resp.StatusCode, string(bodyBytes))
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM events WHERE id = $1", eventID).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count events: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 events after delete, got %d", count)
	}
}

func TestDeleteEvent_NotFound(t *testing.T) {
	t.Cleanup(func() { truncateEvents(t) })

	nonExistentID := uuid.New().String()

	req, err := http.NewRequest(http.MethodDelete, calendarURL+"/events/"+nonExistentID, nil)
	if err != nil {
		t.Fatalf("failed to create DELETE request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to DELETE /events: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Errorf("expected status %d, got %d: %s", http.StatusNotFound, resp.StatusCode, string(bodyBytes))
	}
}

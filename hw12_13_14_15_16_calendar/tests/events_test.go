//go:build integration

package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq" //nolint:depguard
)

var (
	db          *sql.DB
	calendarURL string
)

func TestMain(m *testing.M) {
	dbConn := os.Getenv("DB_CONN")
	if dbConn == "" {
		dbConn = "postgres://postgres:postgres@localhost:5432/calendar?sslmode=disable"
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
		calendarHost = "localhost"
	}
	calendarPort := os.Getenv("CALENDAR_HTTP_PORT")
	if calendarPort == "" {
		calendarPort = "8080"
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

	body := map[string]interface{}{
		"title":    "Test Event",
		"time":     "2025-01-01T10:00:00Z",
		"duration": "2h0m0s",
		"ownerId":  "user-123",
	}
	bodyBytes, _ := json.Marshal(body)

	resp, err := http.Post(calendarURL+"/events", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("failed to POST /events: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM events WHERE title = 'Test Event' AND owner_id = 'user-123'").Scan(&count)
	if err != nil {
		t.Fatalf("failed to count events: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 event in db, got %d", count)
	}

	var id string
	err = db.QueryRow("SELECT id FROM events WHERE title = 'Test Event'").Scan(&id)
	if err != nil {
		t.Fatalf("event not found in db: %v", err)
	}
	if id == "" {
		t.Error("expected non-empty id from db")
	}
}

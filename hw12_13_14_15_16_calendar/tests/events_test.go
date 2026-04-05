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

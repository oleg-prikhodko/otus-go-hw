//go:build integration

package sqlstorage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/common"
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

var testStorage *Storage

func TestMain(m *testing.M) {
	logger := &mockLogger{}
	//testStorage = New(logger, os.Getenv("TEST_DB_CONN"))
	testStorage = New(logger, "postgres://postgres:postgres@localhost:5432/calendar?sslmode=disable")

	if err := testStorage.Connect(); err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	code := m.Run()

	if err := testStorage.Close(); err != nil {
		panic(fmt.Sprintf("failed to close storage: %v", err))
	}

	os.Exit(code)
}

func truncateTable(t *testing.T) {
	t.Helper()
	_, err := testStorage.db.Exec("TRUNCATE TABLE events")
	if err != nil {
		t.Errorf("failed to truncate table: %v", err)
	}
}

func TestSqlStorage_Create(t *testing.T) {
	ev := storage.Event{
		ID:           uuid.New().String(),
		Title:        "test title",
		Time:         time.Now().Add(time.Hour * 72),
		Duration:     time.Hour * 2,
		Description:  nil,
		OwnerID:      uuid.New().String(),
		NotifyBefore: nil,
	}

	t.Cleanup(func() { truncateTable(t) })

	t.Run("success", func(t *testing.T) {
		if err := testStorage.Create(ev); err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		evFromDb := storage.Event{}
		if err := testStorage.db.Get(&evFromDb, "SELECT * FROM events WHERE id = $1", ev.ID); err != nil {
			t.Fatalf("failed to get event from DB: %v", err)
		}

		if ev.Title != evFromDb.Title {
			t.Fatalf("title mismatch: expected %v, got %v", ev.Title, evFromDb.Title)
		}
		if ev.Time.UTC().Format(time.RFC3339) != evFromDb.Time.UTC().Format(time.RFC3339) {
			t.Fatalf("time mismatch: expected %v, got %v", ev.Time.UTC().Format(time.RFC3339), evFromDb.Time.UTC().Format(time.RFC3339))
		}
		if ev.Duration != evFromDb.Duration {
			t.Fatalf("duration mismatch: expected %v, got %v", ev.Duration, evFromDb.Duration)
		}
		if ev.OwnerID != evFromDb.OwnerID {
			t.Fatalf("owner_id mismatch: expected %v, got %v", ev.OwnerID, evFromDb.OwnerID)
		}
	})
}

func TestSqlStorage_Update(t *testing.T) {
	ev := storage.Event{
		ID:           uuid.New().String(),
		Title:        "test title",
		Time:         time.Now().Add(time.Hour * 72),
		Duration:     time.Hour * 2,
		Description:  nil,
		OwnerID:      uuid.New().String(),
		NotifyBefore: nil,
	}

	t.Cleanup(func() { truncateTable(t) })

	// Insert initial event
	query := `INSERT INTO events (id, title, event_time, duration, description, owner_id, notify_before)
		VALUES (:id, :title, :event_time, :duration, :description, :owner_id, :notify_before)`
	_, err := testStorage.db.NamedExec(query, ev)
	if err != nil {
		t.Fatalf("failed to insert test event: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		evFromDb := storage.Event{}
		if err := testStorage.db.Get(&evFromDb, "SELECT * FROM events WHERE id = $1", ev.ID); err != nil {
			t.Fatalf("failed to get event: %v", err)
		}

		evFromDb.Title = "new title"
		if err := testStorage.Update(evFromDb); err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		if err := testStorage.db.Get(&evFromDb, "SELECT * FROM events WHERE id = $1", evFromDb.ID); err != nil {
			t.Fatalf("failed to get updated event: %v", err)
		}

		if evFromDb.Title != "new title" {
			t.Fatalf("expected title 'new title', got %v", evFromDb.Title)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		evFromDb := storage.Event{
			ID:       uuid.New().String(),
			Title:    "nonexistent",
			Time:     time.Now(),
			Duration: time.Hour,
			OwnerID:  uuid.New().String(),
		}
		err := testStorage.Update(evFromDb)
		if err == nil {
			t.Fatal("expected error for non-existent ID")
		}
		if !errors.Is(err, common.NotFoundErr) {
			t.Errorf("expected common.NotFoundErr, got %v", err)
		}
	})
}

func TestSqlStorage_Delete(t *testing.T) {
	ev := storage.Event{
		ID:           uuid.New().String(),
		Title:        "test title",
		Time:         time.Now().Add(time.Hour * 72),
		Duration:     time.Hour * 2,
		Description:  nil,
		OwnerID:      uuid.New().String(),
		NotifyBefore: nil,
	}

	t.Cleanup(func() { truncateTable(t) })

	// Insert test event
	query := `INSERT INTO events (id, title, event_time, duration, description, owner_id, notify_before)
		VALUES (:id, :title, :event_time, :duration, :description, :owner_id, :notify_before)`
	_, err := testStorage.db.NamedExec(query, ev)
	if err != nil {
		t.Fatalf("failed to insert test event: %v", err)
	}

	var id string
	if err := testStorage.db.QueryRowx("SELECT id FROM events LIMIT 1").Scan(&id); err != nil {
		t.Fatalf("failed to get event ID: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		if err := testStorage.Delete(id); err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		var ev storage.Event
		if err := testStorage.db.QueryRowx("SELECT * FROM events WHERE id = $1", id).Scan(&ev); !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("expected sql.ErrNoRows, got %v", err)
		}
	})
}

func TestSqlStorage_List(t *testing.T) {
	ev := storage.Event{
		ID:           uuid.New().String(),
		Title:        "test title",
		Time:         time.Now().Add(time.Hour * 72),
		Duration:     time.Hour * 2,
		Description:  nil,
		OwnerID:      uuid.New().String(),
		NotifyBefore: nil,
	}

	t.Cleanup(func() { truncateTable(t) })

	// Insert test event
	query := `INSERT INTO events (id, title, event_time, duration, description, owner_id, notify_before)
		VALUES (:id, :title, :event_time, :duration, :description, :owner_id, :notify_before)`
	_, err := testStorage.db.NamedExec(query, ev)
	if err != nil {
		t.Fatalf("failed to insert test event: %v", err)
	}

	t.Run("with_intersection", func(t *testing.T) {
		events, err := testStorage.List(ev.Time.Add(-time.Hour*24), ev.Time.Add(time.Hour*24))
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(events) != 1 {
			t.Fatalf("expected 1 event, got %d", len(events))
		}

		if events[0].Title != ev.Title {
			t.Fatalf("expected title %v, got %v", ev.Title, events[0].Title)
		}
	})

	t.Run("no_intersection", func(t *testing.T) {
		events, err := testStorage.List(ev.Time.Add(time.Hour*24), ev.Time.Add(time.Hour*48))
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(events) != 0 {
			t.Fatalf("expected 0 events, got %d", len(events))
		}
	})
}

type mockLogger struct{}

func (m mockLogger) Info(_ string)  {}
func (m mockLogger) Error(_ string) {}

package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

func TestSqlStorage(t *testing.T) {
	logger := &mockLogger{}
	st := New(logger)

	err := st.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/calendar?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close(context.Background())

	truncate := func() {
		fmt.Println("cleaning up")
		_, _ = st.db.Exec("TRUNCATE TABLE events")
	}

	ev := storage.Event{
		Title:        "test title",
		Time:         time.Now().Add(time.Hour * 72),
		Duration:     time.Hour * 2,
		Description:  nil,
		OwnerID:      uuid.New().String(),
		NotifyBefore: nil,
	}

	t.Run("create", func(t *testing.T) {
		t.Cleanup(truncate)

		if err := st.Create(context.Background(), &ev); err != nil {
			t.Fatal(err)
		}

		evFromDb := storage.Event{}
		if err := st.db.Get(&evFromDb, "SELECT * FROM events WHERE id = $1", ev.ID); err != nil {
			t.Fatal(err)
		}

		if ev.Title != evFromDb.Title {
			t.Fatalf("expected %v, got %v\n", ev.Title, evFromDb.Title)
		}
		if ev.Time.UTC().Format(time.RFC3339) != evFromDb.Time.UTC().Format(time.RFC3339) {
			t.Fatalf("expected %v, got %v\n", ev.Time.UTC().Format(time.RFC3339), evFromDb.Time.UTC().Format(time.RFC3339))
		}
		if ev.Duration != evFromDb.Duration {
			t.Fatalf("expected %v, got %v\n", ev.Duration, evFromDb.Duration)
		}
		if ev.Description != evFromDb.Description {
			t.Fatalf("expected %v, got %v\n", ev.Description, evFromDb.Description)
		}
		if ev.OwnerID != evFromDb.OwnerID {
			t.Fatalf("expected %v, got %v\n", ev.OwnerID, evFromDb.OwnerID)
		}
		if ev.NotifyBefore != evFromDb.NotifyBefore {
			t.Fatalf("expected %v, got %v\n", ev.NotifyBefore, evFromDb.NotifyBefore)
		}
	})

	t.Run("update", func(t *testing.T) {
		t.Cleanup(truncate)

		query := `
		INSERT INTO events (title, event_time, duration, description, owner_id, notify_before)
		VALUES (:title, :event_time, :duration, :description, :owner_id, :notify_before)`

		_, err := st.db.NamedExec(query, ev)
		if err != nil {
			t.Fatal(err)
		}

		evFromDb := storage.Event{}
		if err := st.db.Get(&evFromDb, "SELECT * FROM events LIMIT 1"); err != nil {
			t.Fatal(err)
		}
		evFromDb.Title = "new title"

		if err := st.Update(context.Background(), evFromDb.ID, &evFromDb); err != nil {
			t.Fatal(err)
		}

		if err := st.db.Get(&evFromDb, "SELECT * FROM events WHERE id = $1", evFromDb.ID); err != nil {
			t.Fatal(err)
		}

		if evFromDb.Title != "new title" {
			t.Fatalf("expected %v, got %v\n", evFromDb.Title, "new title")
		}
	})

	t.Run("delete", func(t *testing.T) {
		t.Cleanup(truncate)

		query := `
		INSERT INTO events (title, event_time, duration, description, owner_id, notify_before)
		VALUES (:title, :event_time, :duration, :description, :owner_id, :notify_before)`
		_, err := st.db.NamedExec(query, ev)
		if err != nil {
			t.Fatal(err)
		}

		var id string
		if err := st.db.QueryRowx("select id from events limit 1").Scan(&id); err != nil {
			t.Fatal(err)
		}

		if err := st.Delete(context.Background(), id); err != nil {
			t.Fatal(err)
		}

		var ev storage.Event
		if err := st.db.QueryRowx("select * from events where id = $1", id).Scan(&ev); !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("expected sql.ErrNoRows, got %v\n", err)
		}
	})

	t.Run("list", func(t *testing.T) {
		t.Cleanup(truncate)

		query := `
		INSERT INTO events (title, event_time, duration, description, owner_id, notify_before)
		VALUES (:title, :event_time, :duration, :description, :owner_id, :notify_before)`
		_, err := st.db.NamedExec(query, ev)
		if err != nil {
			t.Fatal(err)
		}

		// has intersection
		events, err := st.List(context.Background(), ev.Time.Add(-time.Hour*24), ev.Time.Add(time.Hour*24))
		if err != nil {
			t.Fatal(err)
		}

		if len(events) != 1 {
			t.Fatalf("expected 1 event, got %v", len(events))
		}

		if events[0].Title != ev.Title {
			t.Fatalf("expected %v, got %v", ev.Title, events[0].Title)
		}

		// no intersection
		events, err = st.List(context.Background(), ev.Time.Add(time.Hour*24), ev.Time.Add(time.Hour*48))
		if err != nil {
			t.Fatal(err)
		}

		if len(events) != 0 {
			t.Fatalf("expected 0 event, got %v", len(events))
		}
	})
}

type mockLogger struct {
}

func (m mockLogger) Info(_ string) {
}

func (m mockLogger) Error(_ string) {
}

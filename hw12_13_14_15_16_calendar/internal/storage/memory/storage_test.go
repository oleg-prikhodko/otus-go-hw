package memorystorage

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/common"  //nolint:depguard
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

func TestStorage_Create(t *testing.T) {
	s := New()

	ev := storage.Event{
		ID:    "test-id",
		Title: "Test Event",
		Time:  time.Now(),
	}

	err := s.Create(ev)
	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}
}

func TestStorage_Update(t *testing.T) {
	t.Run("update existing event", func(t *testing.T) {
		s := New()

		ev := storage.Event{
			ID:    "test-id",
			Title: "Test Event",
			Time:  time.Now(),
		}

		if err := s.Create(ev); err != nil {
			t.Fatalf("Create() returned error: %v", err)
		}

		ev.Title = "Updated Event"
		err := s.Update(ev)
		if err != nil {
			t.Fatalf("Update() returned error: %v", err)
		}
	})

	t.Run("update non-existing event", func(t *testing.T) {
		s := New()

		ev := storage.Event{
			ID:    "non-existing-id",
			Title: "Test Event",
			Time:  time.Now(),
		}

		err := s.Update(ev)
		if err == nil {
			t.Fatal("Update() expected error for non-existing event, got nil")
		}
		if !errors.Is(err, common.ErrNotFound) {
			t.Errorf("Update() expected ErrNotFound, got: %v", err)
		}
	})
}

func TestStorage_Delete(t *testing.T) {
	t.Run("delete existing event", func(t *testing.T) {
		s := New()

		ev := storage.Event{
			ID:    "test-id",
			Title: "Test Event",
			Time:  time.Now(),
		}

		if err := s.Create(ev); err != nil {
			t.Fatalf("Create() returned error: %v", err)
		}

		err := s.Delete(ev.ID)
		if err != nil {
			t.Fatalf("Delete() returned error: %v", err)
		}
	})

	t.Run("delete non-existing event", func(t *testing.T) {
		s := New()

		err := s.Delete("non-existing-id")
		if err == nil {
			t.Fatal("Delete() expected error for non-existing event, got nil")
		}
		if !errors.Is(err, common.ErrNotFound) {
			t.Errorf("Delete() expected ErrNotFound, got: %v", err)
		}
	})
}

func TestStorage_List(t *testing.T) {
	t.Run("list events in range", func(t *testing.T) {
		s := New()

		now := time.Now()
		events := []storage.Event{
			{ID: "1", Title: "Event 1", Time: now.Add(-2 * time.Hour)},
			{ID: "2", Title: "Event 2", Time: now.Add(-1 * time.Hour)},
			{ID: "3", Title: "Event 3", Time: now},
			{ID: "4", Title: "Event 4", Time: now.Add(1 * time.Hour)},
			{ID: "5", Title: "Event 5", Time: now.Add(2 * time.Hour)},
		}

		for _, ev := range events {
			if err := s.Create(ev); err != nil {
				t.Fatalf("Create() returned error: %v", err)
			}
		}

		result, err := s.List(now.Add(-30*time.Minute), now.Add(30*time.Minute))
		if err != nil {
			t.Fatalf("List() returned error: %v", err)
		}

		if len(result) != 1 {
			t.Errorf("List() expected 1 event, got %d", len(result))
		}
		if len(result) > 0 && result[0].ID != "3" {
			t.Errorf("List() expected event ID 3, got %s", result[0].ID)
		}
	})

	t.Run("list events empty range", func(t *testing.T) {
		s := New()

		result, err := s.List(time.Now(), time.Now().Add(time.Hour))
		if err != nil {
			t.Fatalf("List() returned error: %v", err)
		}

		if len(result) != 0 {
			t.Errorf("List() expected 0 events, got %d", len(result))
		}
	})
}

func TestStorage_ThreadSafety(t *testing.T) {
	t.Run("concurrent creates", func(_ *testing.T) {
		s := New()
		var wg sync.WaitGroup
		numGoroutines := 100

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				ev := storage.Event{
					ID:    string(rune(id)),
					Title: "Event",
					Time:  time.Now(),
				}
				_ = s.Create(ev)
			}(i)
		}
		wg.Wait()
	})

	t.Run("concurrent updates", func(t *testing.T) {
		s := New()

		ev := storage.Event{
			ID:    "test-id",
			Title: "Test Event",
			Time:  time.Now(),
		}
		if err := s.Create(ev); err != nil {
			t.Fatalf("Create() returned error: %v", err)
		}

		var wg sync.WaitGroup
		numGoroutines := 100

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				updated := storage.Event{
					ID:    "test-id",
					Title: "Updated",
					Time:  time.Now(),
				}
				_ = s.Update(updated)
			}()
		}
		wg.Wait()
	})

	t.Run("concurrent deletes", func(t *testing.T) {
		s := New()

		var wg sync.WaitGroup
		numGoroutines := 100

		for i := 0; i < numGoroutines; i++ {
			ev := storage.Event{
				ID:    string(rune(i)),
				Title: "Event",
				Time:  time.Now(),
			}
			if err := s.Create(ev); err != nil {
				t.Fatalf("Create() returned error: %v", err)
			}
		}

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()
				_ = s.Delete(string(rune(id)))
			}(i)
		}
		wg.Wait()
	})
}

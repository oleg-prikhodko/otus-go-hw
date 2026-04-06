//go:build integration

package tests

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"              //nolint:depguard
	amqp "github.com/rabbitmq/amqp091-go" //nolint:depguard
)

var (
	rabbitmqURL     string
	sentEventsQueue = "sent_events"
)

func createTestEventForNotification(t *testing.T, title string) string {
	t.Helper()

	notifyBefore := 5 * time.Minute
	eventTime := time.Now().Add(notifyBefore + time.Second*10)

	id := uuid.New().String()
	ownerID := uuid.New().String()
	duration := time.Hour

	_, err := db.Exec(`
		INSERT INTO events (id, title, event_time, duration, owner_id, notify_before)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, id, title, eventTime, int64(duration), ownerID, int64(notifyBefore))
	if err != nil {
		t.Fatalf("failed to insert event: %v", err)
	}

	return id
}

func TestSchedulerSender_SendsToSentEventsQueue(t *testing.T) {
	t.Cleanup(func() { truncateEvents(t) })

	rabbitmqURL = os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		t.Fatal("RABBITMQ_URL env var is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		t.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("failed to open channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		sentEventsQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to register consumer: %v", err)
	}

	eventID := createTestEventForNotification(t, "Notification Test Event")

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("timeout waiting for message on %s queue. eventID=%s", sentEventsQueue, eventID)
		case msg, ok := <-msgs:
			if !ok {
				t.Fatalf("channel closed unexpectedly")
			}
			var ev struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(msg.Body, &ev); err != nil {
				t.Logf("failed to unmarshal message: %v", err)
				continue
			}
			if ev.ID != eventID {
				t.Fatalf("unexpected eventID: expected %s, got %s", eventID, ev.ID)
			}
			return
		case <-time.After(time.Second):
		}
	}
}

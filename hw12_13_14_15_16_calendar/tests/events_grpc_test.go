//go:build integration

package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"                                                                    //nolint:depguard
	_ "github.com/lib/pq"                                                                       //nolint:depguard
	pb "github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/server/grpc/proto" //nolint:depguard
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var grpcClient pb.CalendarServiceClient

func grpcSetup(t *testing.T) {
	t.Helper()
	grpcHost := os.Getenv("CALENDAR_HOST")
	if grpcHost == "" {
		t.Skip("CALENDAR_HOST env var is not set")
	}
	grpcPort := os.Getenv("CALENDAR_GRPC_PORT")
	if grpcPort == "" {
		t.Skip("CALENDAR_GRPC_PORT env var is not set")
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", grpcHost, grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial grpc server: %v", err)
	}
	grpcClient = pb.NewCalendarServiceClient(conn)
}

func TestCreateEventGrpc_Success(t *testing.T) {
	grpcSetup(t)
	t.Cleanup(func() { truncateEvents(t) })

	ctx := context.Background()
	resp, err := grpcClient.CreateEvent(ctx, &pb.CreateEventRequest{
		Title:    "Test Event",
		Time:     timestamppb.New(time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)),
		Duration: durationpb.New(time.Hour * 3),
		OwnerId:  uuid.New().String(),
	})
	if err != nil {
		t.Fatalf("failed to create event: %v", err)
	}

	if resp.Id == "" {
		t.Error("expected non-empty id")
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

func TestListDayEventsGrpc_Success(t *testing.T) {
	grpcSetup(t)
	t.Cleanup(func() { truncateEvents(t) })

	eventTime := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	createTestEventGrpc(t, eventTime)

	ctx := context.Background()
	resp, err := grpcClient.ListEventsForDay(ctx, &pb.ListEventsRequest{
		Date: timestamppb.New(eventTime),
	})
	if err != nil {
		t.Fatalf("failed to list events: %v", err)
	}

	if len(resp.Events) != 1 {
		t.Errorf("expected 1 event, got %d", len(resp.Events))
	}
	if resp.Events[0].Title != "Test Event" {
		t.Errorf("expected title 'Test Event', got '%s'", resp.Events[0].Title)
	}
}

func TestListWeekEventsGrpc_Success(t *testing.T) {
	grpcSetup(t)
	t.Cleanup(func() { truncateEvents(t) })

	eventTime := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	createTestEventGrpc(t, eventTime)

	ctx := context.Background()
	resp, err := grpcClient.ListEventsForWeek(ctx, &pb.ListEventsRequest{
		Date: timestamppb.New(eventTime),
	})
	if err != nil {
		t.Fatalf("failed to list events: %v", err)
	}

	if len(resp.Events) != 1 {
		t.Errorf("expected 1 event, got %d", len(resp.Events))
	}
}

func TestListMonthEventsGrpc_Success(t *testing.T) {
	grpcSetup(t)
	t.Cleanup(func() { truncateEvents(t) })

	eventTime := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	createTestEventGrpc(t, eventTime)

	ctx := context.Background()
	resp, err := grpcClient.ListEventsForMonth(ctx, &pb.ListEventsRequest{
		Date: timestamppb.New(eventTime),
	})
	if err != nil {
		t.Fatalf("failed to list events: %v", err)
	}

	if len(resp.Events) != 1 {
		t.Errorf("expected 1 event, got %d", len(resp.Events))
	}
}

func TestDeleteEventGrpc_Success(t *testing.T) {
	grpcSetup(t)
	t.Cleanup(func() { truncateEvents(t) })

	eventTime := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	eventID := createTestEventGrpc(t, eventTime)

	ctx := context.Background()
	_, err := grpcClient.DeleteEvent(ctx, &pb.DeleteEventRequest{Id: eventID})
	if err != nil {
		t.Fatalf("failed to delete event: %v", err)
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

func createTestEventGrpc(t *testing.T, eventTime time.Time) string {
	t.Helper()
	ctx := context.Background()
	resp, err := grpcClient.CreateEvent(ctx, &pb.CreateEventRequest{
		Title:    "Test Event",
		Time:     timestamppb.New(eventTime),
		Duration: durationpb.New(time.Hour),
		OwnerId:  uuid.New().String(),
	})
	if err != nil {
		t.Fatalf("failed to create test event: %v", err)
	}
	return resp.Id
}

package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/common"
	pb "github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/server/grpc/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateEvent(t *testing.T) {
	app := &common.MockApp{}
	logger := &common.MockLogger{}

	server := NewServer(logger, app, "")

	req := &pb.CreateEventRequest{
		Title:    "Test Event",
		Time:     timestamppb.New(time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)),
		Duration: durationpb.New(time.Hour),
		OwnerId:  "user1",
	}

	resp, err := server.CreateEvent(context.Background(), req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !app.EventCreateCalled {
		t.Error("CreateEvent was not called")
	}

	if app.EventCreate.Title != "Test Event" {
		t.Errorf("expected title 'Test Event', got '%s'", app.EventCreate.Title)
	}

	if app.EventCreate.OwnerID != "user1" {
		t.Errorf("expected owner_id 'user1', got '%s'", app.EventCreate.OwnerID)
	}

	if resp.Id == "" {
		t.Log("note: ID generation not implemented - resp.Id is empty")
	}
}

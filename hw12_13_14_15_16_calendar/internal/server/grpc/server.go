package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/common"               //nolint:depguard
	pb "github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/server/grpc/proto" //nolint:depguard
	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/storage"              //nolint:depguard
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedCalendarServiceServer
	server *grpc.Server
	app    common.Application
	logger common.Logger
	addr   string
}

func NewServer(logger common.Logger, app common.Application, addr string) *Server {
	return &Server{
		app:    app,
		logger: logger,
		addr:   addr,
	}
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.server = grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor(s.logger)))
	pb.RegisterCalendarServiceServer(s.server, s)

	s.logger.Info(fmt.Sprintf("starting gRPC server at %s", s.addr))

	go func() {
		if err := s.server.Serve(lis); err != nil {
			s.logger.Error(fmt.Sprintf("gRPC server error: %s", err))
		}
	}()

	<-ctx.Done()
	s.logger.Info("shutting down gRPC server")

	s.Stop()
	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}

func (s *Server) CreateEvent(_ context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	ev := storage.Event{
		Title:        req.Title,
		Time:         req.Time.AsTime(),
		Duration:     req.Duration.AsDuration(),
		Description:  &req.Description,
		OwnerID:      req.OwnerId,
		NotifyBefore: parseProtoDuration(req.NotifyBefore),
	}

	id, err := s.app.CreateEvent(ev)
	if err != nil {
		return nil, err
	}

	return &pb.CreateEventResponse{Id: id}, nil
}

func (s *Server) UpdateEvent(_ context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	ev := storage.Event{
		ID:           req.Id,
		Title:        req.Title,
		Time:         req.Time.AsTime(),
		Duration:     req.Duration.AsDuration(),
		Description:  &req.Description,
		OwnerID:      req.OwnerId,
		NotifyBefore: parseProtoDuration(req.NotifyBefore),
	}

	if err := s.app.UpdateEvent(ev); err != nil {
		return nil, err
	}

	return &pb.UpdateEventResponse{}, nil
}

func (s *Server) DeleteEvent(_ context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	if err := s.app.DeleteEvent(req.Id); err != nil {
		return nil, err
	}

	return &pb.DeleteEventResponse{}, nil
}

func (s *Server) ListEventsForDay(_ context.Context, req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	events, err := s.app.ListEventsForDay(req.Date.AsTime())
	if err != nil {
		return nil, err
	}

	return &pb.ListEventsResponse{Events: eventsToProto(events)}, nil
}

func (s *Server) ListEventsForWeek(_ context.Context, req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	events, err := s.app.ListEventsForWeek(req.Date.AsTime())
	if err != nil {
		return nil, err
	}

	return &pb.ListEventsResponse{Events: eventsToProto(events)}, nil
}

func (s *Server) ListEventsForMonth(_ context.Context, req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	events, err := s.app.ListEventsForMonth(req.Date.AsTime())
	if err != nil {
		return nil, err
	}

	return &pb.ListEventsResponse{Events: eventsToProto(events)}, nil
}

func parseProtoDuration(d *durationpb.Duration) *time.Duration {
	if d == nil {
		return nil
	}
	v := d.AsDuration()
	return &v
}

func eventsToProto(events []storage.Event) []*pb.Event {
	result := make([]*pb.Event, len(events))
	for i, ev := range events {
		var description string
		if ev.Description != nil {
			description = *ev.Description
		}

		var notifyBefore *durationpb.Duration
		if ev.NotifyBefore != nil {
			notifyBefore = durationpb.New(*ev.NotifyBefore)
		}

		result[i] = &pb.Event{
			Id:           ev.ID,
			Title:        ev.Title,
			Time:         timestamppb.New(ev.Time),
			Duration:     durationpb.New(ev.Duration),
			Description:  description,
			OwnerId:      ev.OwnerID,
			NotifyBefore: notifyBefore,
		}
	}
	return result
}

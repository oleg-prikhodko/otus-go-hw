package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/oleg-prikhodko/otus-go-hw/hw12_13_14_15_calendar/internal/common" //nolint:depguard
)

type Server struct {
	server *http.Server
	app    common.Application
	logger common.Logger
}

func NewServer(logger common.Logger, app common.Application, addr string) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("hello world"))
	}))

	server := &http.Server{Addr: addr, Handler: loggingMiddleware(logger, mux)} //nolint:gosec

	return &Server{server, app, logger}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info(fmt.Sprintf("starting server at %s", s.server.Addr))

	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			s.logger.Error(fmt.Sprintf("http server error: %s", err))
		}
	}()

	<-ctx.Done()
	s.logger.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return s.Stop(shutdownCtx)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

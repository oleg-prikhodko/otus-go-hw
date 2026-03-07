package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	app    Application
	logger Logger
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application, addr string) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello world"))
	}))

	server := &http.Server{Addr: addr, Handler: loggingMiddleware(logger, mux)}

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

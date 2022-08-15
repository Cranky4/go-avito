package internalhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	logger     Logger
	app        Application
	httpServer *http.Server
}

type Logger interface {
	Info(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, id storage.EventID, title string, startsAt, endsAt time.Time) error
	UpdateEvent(ctx context.Context, id storage.EventID, title string, startsAt, endsAt time.Time) error
	GetEvent(ctx context.Context, id storage.EventID) (storage.Event, error)
	DeleteEvent(ctx context.Context, id storage.EventID) error
	GetDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetWeekEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetMonthEvents(ctx context.Context, date time.Time) ([]storage.Event, error)
}

func NewServer(logger Logger, app Application, addr string) *Server {
	handler := &httpHandler{}

	mux := http.NewServeMux()
	mux.Handle("/", loggingMiddleware(handler))

	httpServer := &http.Server{
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           mux,
		Addr:              addr,
	}

	return &Server{
		logger:     logger,
		app:        app,
		httpServer: httpServer,
	}
}

func (s *Server) Start(ctx context.Context) error {
	go s.httpServer.ListenAndServe()
	s.logger.Info("http server started...")

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.httpServer.Close()
	s.logger.Info("http server stopped...")

	return nil
}

type httpHandler struct{}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("hello world")
}

package internalhttp

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	logger                Logger
	app                   Application
	httpServer            *http.Server
	requestLogFileHandler *os.File
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

func NewServer(logger Logger, app Application, addr, requestLogFile string) *Server {
	handler := &httpHandler{}

	file, err := os.OpenFile(requestLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	requestLogger := log.New(file, "", log.LstdFlags|log.Lshortfile)

	mux := http.NewServeMux()
	mux.Handle("/", loggingMiddleware(handler, requestLogger))

	httpServer := &http.Server{
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           mux,
		Addr:              addr,
	}

	return &Server{
		logger:                logger,
		app:                   app,
		httpServer:            httpServer,
		requestLogFileHandler: file,
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

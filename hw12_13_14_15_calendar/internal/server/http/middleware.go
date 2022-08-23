package internalhttp

import (
	"log"
	"net/http"
	"time"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func loggingMiddleware(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}

		started := time.Now()
		next.ServeHTTP(recorder, r)
		finished := time.Now()

		logger.Printf(
			"%s [%s] %s %s?%s %s %d %d \"%s\"",
			r.RemoteAddr,
			time.Now().Format(time.RFC822Z),
			r.Method,
			r.URL.Path,
			r.URL.Query().Encode(),
			r.Proto,
			recorder.Status,
			finished.Sub(started).Microseconds(),
			r.Header.Get("User-Agent"),
		)
	})
}

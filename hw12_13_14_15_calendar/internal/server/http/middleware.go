package internalhttp

import (
	"log"
	"net/http"
	"os"
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.OpenFile("./logs/http-server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		log.SetOutput(file)

		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}

		started := time.Now()
		next.ServeHTTP(recorder, r)
		finished := time.Now()

		log.Printf(
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

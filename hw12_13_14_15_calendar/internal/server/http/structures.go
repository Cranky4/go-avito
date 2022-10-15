package internalhttp

import "time"

type EventRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	StartsAt    string `json:"startsAt"`
	EndsAt      string `json:"endsAt"`
	NotifyAfter string `json:"notify,omitempty"`
}

type ErrorResponse struct {
	Code    int
	Message string
	Data    []interface{}
}

type EventResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	StartsAt    time.Time  `json:"startsAt"`
	EndsAt      time.Time  `json:"endsAt"`
	NotifyAfter *time.Time `json:"notify,omitempty"`
}

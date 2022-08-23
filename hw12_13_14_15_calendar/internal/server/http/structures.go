package internalhttp

type EventRequest struct {
	ID       string
	Title    string
	StartsAt string
	EndsAt   string
}

type ErrorResponse struct {
	Code    int
	Message string
	Data    []interface{}
}

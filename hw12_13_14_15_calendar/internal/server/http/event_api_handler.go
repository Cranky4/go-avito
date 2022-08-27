package internalhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/storage"
)

type EventAPIHandler struct {
	app Application
}

func NewEventAPIHandler(app Application) *EventAPIHandler {
	return &EventAPIHandler{
		app: app,
	}
}

func (h *EventAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	switch r.Method {
	case http.MethodPost:
		/*
			POST /events
			{
			    "id": "48cd8858-9103-4c6a-9a83-1d58307f071b",
			    "title": "zxc",
			    "startsAt": "2022-08-23T15:04:05+07:00",
			    "endsAt": "2022-08-23T15:04:05+07:00"
			}
		*/
		createEvent(ctx, h.app, w, r)
	case http.MethodDelete:
		/*
			DELETE /events?id=48cd8858-9103-4c6a-9a83-1d58307f071b
		*/
		deleteEvent(ctx, h.app, w, r)
	case http.MethodPut:
		/*
			PUT /events
			{
			    "id": "48cd8858-9103-4c6a-9a83-1d58307f071b",
			    "title": "zxc",
			    "startsAt": "2022-08-23T15:04:05+07:00",
			    "endsAt": "2022-08-23T15:04:05+07:00"
			}
		*/
		updateEvent(ctx, h.app, w, r)
	default:
		/*
			GET /events?day=2022-08-10&period=month
		*/
		getEvents(ctx, h.app, w, r)
	}
}

func processEvent(ctx context.Context, isNew bool, app Application, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var request EventRequest
	err := decoder.Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusUnprocessableEntity,
			},
		)

		return
	}

	// TODO: validation

	eventID, err := storage.NewEventIDFromString(request.ID)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(
			ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusUnprocessableEntity,
				Data:    []interface{}{"id"},
			},
		)
		return
	}

	startsAt, err := time.Parse(time.RFC3339, request.StartsAt)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(
			ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusUnprocessableEntity,
				Data:    []interface{}{"startsAt"},
			},
		)
		return
	}
	endsAt, err := time.Parse(time.RFC3339, request.EndsAt)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(
			ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusUnprocessableEntity,
				Data:    []interface{}{"endsAt"},
			},
		)
		return
	}

	var notifyAfter storage.NotifyAfter
	if request.NotifyAfter != "" {
		notifyAfterTime, err := time.Parse(time.RFC3339, request.NotifyAfter)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(
				ErrorResponse{
					Message: err.Error(),
					Code:    http.StatusUnprocessableEntity,
					Data:    []interface{}{"notifyAfter"},
				},
			)
			return
		}

		notifyAfter.Time = notifyAfterTime
		notifyAfter.IsSet = true
	}

	if isNew {
		err = app.CreateEvent(
			ctx,
			eventID,
			request.Title,
			startsAt,
			endsAt,
			notifyAfter,
		)
	} else {
		err = app.UpdateEvent(
			ctx,
			eventID,
			request.Title,
			startsAt,
			endsAt,
			notifyAfter,
		)
	}

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(
			ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusUnprocessableEntity,
			},
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func createEvent(ctx context.Context, app Application, w http.ResponseWriter, r *http.Request) {
	processEvent(ctx, true, app, w, r)
}

func getEvents(ctx context.Context, app Application, w http.ResponseWriter, r *http.Request) {
	var events []storage.Event
	var err error

	period := r.URL.Query().Get("period")
	day := r.URL.Query().Get("day")
	startsAt := time.Now()
	if day != "" {
		startsAt, err = time.Parse("2006-01-02", day)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(
				ErrorResponse{
					Message: err.Error(),
					Code:    http.StatusBadRequest,
					Data:    []interface{}{"day"},
				},
			)
			return
		}
	}

	switch period {
	case "month":
		events, err = app.GetMonthEvents(ctx, startsAt)
	case "week":
		events, err = app.GetWeekEvents(ctx, startsAt)
	default:
		events, err = app.GetDayEvents(ctx, startsAt)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			},
		)
		return
	}

	response := make([]EventResponse, 0, len(events))
	for _, ev := range events {
		response = append(response, EventResponse{
			ID:          ev.ID.String(),
			Title:       ev.Title,
			StartsAt:    ev.StartsAt,
			EndsAt:      ev.EndsAt,
			NotifyAfter: ev.NotifyAfter.Time,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func deleteEvent(ctx context.Context, app Application, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		json.NewEncoder(w).Encode(
			ErrorResponse{
				Message: "empty id parameter",
				Code:    http.StatusBadRequest,
				Data:    []interface{}{"id"},
			},
		)
		return
	}

	eventID, err := storage.NewEventIDFromString(id)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(
			ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusUnprocessableEntity,
				Data:    []interface{}{"id"},
			},
		)
		return
	}

	if app.DeleteEvent(ctx, eventID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			ErrorResponse{
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			},
		)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func updateEvent(ctx context.Context, app Application, w http.ResponseWriter, r *http.Request) {
	processEvent(ctx, false, app, w, r)
}

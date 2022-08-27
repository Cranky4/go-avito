package internalgrpc

import (
	"context"
	"log"

	pb "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/api/EventService"
	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/app"
	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler = pb.EventServiceServer

type handler struct {
	pb.UnimplementedEventServiceServer
	app  *app.App
	logg *log.Logger
}

func NewHandler(app *app.App, logger *log.Logger) (Handler, error) {
	return &handler{app: app, logg: logger}, nil
}

func (h *handler) CreateEvent(ctx context.Context, r *(pb.CreateEventRequest)) (*emptypb.Empty, error) {
	h.logg.Printf("%#v", r)

	eventID, err := storage.NewEventIDFromString(r.Id)
	if err != nil {
		return nil, err
	}

	var notifyAfter storage.NotifyAfter
	if r.NotifyAfter != nil {
		notifyAfter.Time = r.NotifyAfter.AsTime()
		notifyAfter.IsSet = true
	}

	err = h.app.CreateEvent(
		ctx,
		eventID,
		r.Title,
		r.StartsAt.AsTime(),
		r.EndsAt.AsTime(),
		notifyAfter,
	)

	return &emptypb.Empty{}, err
}

func (h *handler) UpdateEvent(ctx context.Context, r *(pb.UpdateEventRequest)) (*emptypb.Empty, error) {
	h.logg.Printf("%v", r)

	eventID, err := storage.NewEventIDFromString(r.Id)
	if err != nil {
		return nil, err
	}

	var notifyAfter storage.NotifyAfter
	if r.NotifyAfter != nil {
		notifyAfter.Time = r.NotifyAfter.AsTime()
		notifyAfter.IsSet = true
	}

	err = h.app.UpdateEvent(
		ctx,
		eventID,
		r.Title,
		r.StartsAt.AsTime(),
		r.EndsAt.AsTime(),
		notifyAfter,
	)

	return &emptypb.Empty{}, err
}

func (h *handler) DeleteEvent(ctx context.Context, r *pb.DeleteEventRequest) (*emptypb.Empty, error) {
	h.logg.Printf("%v", r)

	eventID, err := storage.NewEventIDFromString(r.Id)
	if err != nil {
		return nil, err
	}

	err = h.app.DeleteEvent(ctx, eventID)

	return &emptypb.Empty{}, err
}

func (h *handler) GetDayEvents(ctx context.Context, r *timestamppb.Timestamp) (*pb.EventsResponse, error) {
	h.logg.Printf("%v", r)

	evs, err := h.app.GetDayEvents(ctx, r.AsTime())

	return &pb.EventsResponse{Events: collectEvents(evs)}, err
}

func (h *handler) GetWeekEvents(ctx context.Context, r *timestamppb.Timestamp) (*pb.EventsResponse, error) {
	h.logg.Printf("%v", r)

	evs, err := h.app.GetWeekEvents(ctx, r.AsTime())

	return &pb.EventsResponse{Events: collectEvents(evs)}, err
}

func (h *handler) GetMonthEvents(ctx context.Context, r *timestamppb.Timestamp) (*pb.EventsResponse, error) {
	h.logg.Printf("%v", r)

	evs, err := h.app.GetMonthEvents(ctx, r.AsTime())

	return &pb.EventsResponse{Events: collectEvents(evs)}, err
}

func collectEvents(evs []storage.Event) []*pb.Event {
	events := make([]*pb.Event, 0, len(evs))

	for _, e := range evs {
		var notifyAfter timestamppb.Timestamp
		if e.NotifyAfter.IsSet {
			notifyAfter.Seconds = e.NotifyAfter.Time.Unix()
		}

		events = append(events, &pb.Event{
			Id:    e.ID.String(),
			Title: e.Title,
			StartsAt: &timestamppb.Timestamp{
				Seconds: e.StartsAt.Unix(),
			},
			EndsAt: &timestamppb.Timestamp{
				Seconds: e.EndsAt.Unix(),
			},
			NotifyAfter: &notifyAfter,
		})
	}

	return events
}

package internalhttp

import (
	"context"
	"fmt"
	"net"

	"github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/app"
	"google.golang.org/grpc"

	pb "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/api/EventService"
)

type Server struct {
	pb.UnimplementedEventServiceServer
	grpcServer *grpc.Server
	logger     Logger
	app        *app.App
}

func New(app *app.App, logg Logger) *Server {
	return &Server{app: app, logger: logg}
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func (s *Server) Start(ctx context.Context, addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}
	s.grpcServer = grpc.NewServer(opts...)

	handler, err := NewHandler(s.app)
	if err != nil {
		return err
	}

	pb.RegisterEventServiceServer(s.grpcServer, handler)
	go s.grpcServer.Serve(listener)
	s.logger.Info(fmt.Sprintf("grpc server started and listen %s...", addr))

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) {
	s.grpcServer.Stop()
	s.logger.Info("grpc server stopped")
}

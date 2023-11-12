package internalgrpc

import (
	"context"
	"fmt"
	"net"

	event "github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/gen_buf/grpc/v1"
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/server"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Server struct {
	app     server.Application
	logger  server.Logger
	server  *grpc.Server
	address string
}

func NewServer(logger server.Logger, app server.Application, address string) *Server {
	return &Server{
		app:     app,
		logger:  logger,
		address: address,
	}
}

func (s *Server) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		logMiddleware := loggingMiddleware(s.logger)
		s.server = grpc.NewServer(logMiddleware)

		event.RegisterEventServiceServer(s.server, &EventServer{
			App:    s.app,
			Logger: s.logger,
		})

		s.logger.Info(fmt.Sprintf("grpc server start %s", s.address))

		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			return errors.Wrap(err, "I CAN'T HEAR")
		}

		return s.server.Serve(lis)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		s.server.GracefulStop()

		return nil
	}
}

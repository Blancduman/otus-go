package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	pb "github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/gen_buf/grpc/v1"
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/server"
	internalgrpc "github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
)

type Server struct {
	app     server.Application
	logger  server.Logger
	server  *http.Server
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
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			handler := http.NewServeMux()
			handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("hello-world"))
				if err != nil {
					s.logger.Error(err.Error())
				} else {
					w.WriteHeader(200)
				}
			})

			runtimeMux := runtime.NewServeMux()

			err := s.registerGRPC(ctx, runtimeMux)
			if err != nil {
				return errors.Wrap(err, "register grpc api handlers")
			}

			handler.Handle("/api/v1/", runtimeMux)

			s.server = &http.Server{
				Addr:         s.address,
				Handler:      loggingMiddleware(handler, s.logger),
				ReadTimeout:  time.Second * 10,
				WriteTimeout: time.Second * 10,
			}

			s.logger.Info(fmt.Sprintf("rest server start %s", s.address))

			return s.server.ListenAndServe()
		}
	}
}

func (s *Server) Stop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		return s.server.Close()
	}
}

func (s *Server) registerGRPC(ctx context.Context, runtimeMux *runtime.ServeMux) error {
	err := pb.RegisterEventServiceHandlerServer(
		ctx,
		runtimeMux,
		internalgrpc.EventServer{
			App:    s.app,
			Logger: s.logger,
		},
	)
	if err != nil {
		return errors.Wrap(err, "register grpc")
	}

	return nil
}

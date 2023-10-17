package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	app     Application
	logger  Logger
	server  *http.Server
	address string
}

type Logger interface {
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application, address string) *Server {
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
		handler := http.NewServeMux()
		handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("hello-world"))
			if err != nil {
				s.logger.Error(err.Error())
			} else {
				w.WriteHeader(200)
			}
		})

		s.server = &http.Server{
			Addr:         s.address,
			Handler:      loggingMiddleware(handler, s.logger),
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}

		s.logger.Info(fmt.Sprintf("listening %s", s.address))

		return s.server.ListenAndServe()
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

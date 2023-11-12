package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func loggingMiddleware(log server.Logger) grpc.ServerOption {
	return grpc.UnaryInterceptor(
		func(
			ctx context.Context,
			req interface{},
			info *grpc.UnaryServerInfo,
			handler grpc.UnaryHandler,
		) (interface{}, error) {
			now := time.Now()
			response, err := handler(ctx, req)

			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return response, err
			}
			latency := time.Since(now).Milliseconds()

			log.Info(fmt.Sprintf(
				"%s [%s] %s %d %s",
				md.Get(":authority"),
				now.String(),
				info.FullMethod,
				latency,
				md.Get("user-agent"),
			))

			return response, err
		},
	)
}

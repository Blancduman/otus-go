package main

import (
	"context"
	"net/http"
	"time"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/server/grpc"
	memorystorage "github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func grpcServerCmd(ctx context.Context, conf Config) *cobra.Command {
	return &cobra.Command{
		Use:   "grpc-server",
		Short: "start grpc server",
		RunE: func(cmd *cobra.Command, args []string) error {
			var storage app.Storage

			logg := logger.New(conf.Logger.Level)

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			switch conf.Database.Type {
			case "sql":
				storage = sqlstorage.New(conf.Database.URL)
			default:
				storage = memorystorage.New()
			}

			err := storage.Connect(ctx)
			if err != nil {
				return errors.Wrap(err, "could not connect to db")
			}
			defer func() {
				err := storage.Close(ctx)
				if err != nil {
					logg.Error("fail to close storage: " + err.Error())
				}
			}()

			calendar := app.New(logg, storage)
			server := internalgrpc.NewServer(logg, calendar, conf.GetGRPCAddr())

			go func() {
				<-ctx.Done()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				if err := server.Stop(ctx); err != nil {
					logg.Error("failed to stop grpc server: " + err.Error())
				}
			}()

			logg.Info("calendar is running...")

			if err := server.Start(ctx); !errors.Is(err, http.ErrServerClosed) {
				logg.Error(errors.Wrap(err, "failed to start grpc server").Error())
			}

			return nil
		},
	}
}

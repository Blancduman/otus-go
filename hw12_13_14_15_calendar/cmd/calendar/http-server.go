package main

import (
	"context"
	"net/http"
	"time"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func httpServerCmd(ctx context.Context, conf Config) *cobra.Command {
	return &cobra.Command{
		Use:   "http-server --config",
		Short: "start http server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			logg := logger.New(conf.Logger.Level)
			storage := memorystorage.New()
			calendar := app.New(logg, storage)
			server := internalhttp.NewServer(logg, calendar, conf.GetAddr())

			go func() {
				<-ctx.Done()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				if err := server.Stop(ctx); err != nil {
					logg.Error("failed to stop http server: " + err.Error())
				}
			}()

			logg.Info("calendar is running...")

			if err := server.Start(ctx); !errors.Is(err, http.ErrServerClosed) {
				logg.Error(errors.Wrap(err, "failed to start http server").Error())
			}

			return nil
		},
	}
}

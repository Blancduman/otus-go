package main

import (
	"context"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/scheduler"
	memorystorage "github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func schedulerCmd(ctx context.Context, conf Config) *cobra.Command {
	return &cobra.Command{
		Use:   "scheduler",
		Short: "start scheduler",
		RunE: func(cmd *cobra.Command, args []string) error {
			var storage app.Storage

			logg := logger.New(conf.Logger.Level)

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			switch conf.Database.Type {
			case "sql":
				storage = sqlstorage.New(conf.Database.URL)
				defer func() {
					err := storage.Close(ctx)
					if err != nil {
						logg.Error("fail to close storage: " + err.Error())
					}
				}()
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

			rbClient, err := RabbitMqBuild(conf)
			if err != nil {
				return errors.Wrap(err, "could not build rabbitmq client")
			}
			defer func() {
				if err := rbClient.Close(); err != nil {
					logg.Error(errors.Wrap(err, "fail to close rb client").Error())
				}
			}()

			schedler := scheduler.New(
				storage,
				rbClient,
				logg,
				conf.Scheduler.Period,
				conf.Scheduler.Mark,
				conf.RabbitMQ.Exchange,
				conf.RabbitMQ.Key,
			)

			logg.Info("scheduler running")

			if err := schedler.Start(ctx); err != nil {
				logg.Error(errors.Wrap(err, "fail to start scheduler").Error())
			}

			return nil
		},
	}
}

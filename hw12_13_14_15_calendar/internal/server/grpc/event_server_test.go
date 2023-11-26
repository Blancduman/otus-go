package internalgrpc_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/app"
	event "github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/gen_buf/grpc/v1"
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/server/grpc"
	memorystorage "github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_EventServer(t *testing.T) {
	ctx := context.TODO()

	server, err := getGRPCServer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := server.Stop(ctx)
		require.NoError(t, err)
	}()

	t.Run("add", func(t *testing.T) {
		ctx := context.TODO()

		client := getGRPCClient(ctx)
		e, err := client.Add(ctx, &event.AddRequest{
			Title:       "Test title 1",
			Description: "Test description 1",
			OwnerId:     1,
			StartDate:   timestamppb.New(time.Now().Add(time.Hour * 2)),
			EndDate:     timestamppb.New(time.Now().Add(time.Hour * 3)),
			RemindIn:    timestamppb.New(time.Now().Add(time.Hour*2 + time.Minute*30)),
		})
		require.NotNil(t, e)
		require.NoError(t, err)

		e2, err2 := client.Add(ctx, &event.AddRequest{
			Title:       "Test title 3",
			Description: "Test description 3",
			OwnerId:     1,
			StartDate:   timestamppb.New(time.Now().Add(time.Hour * 20)),
			EndDate:     timestamppb.New(time.Now().Add(time.Hour * 21)),
			RemindIn:    timestamppb.New(time.Now().Add(time.Hour*20 + time.Minute*30)),
		})
		require.NotNil(t, e2)
		require.NoError(t, err2)
	})

	t.Run("edit", func(t *testing.T) {
		ctx := context.TODO()

		client := getGRPCClient(ctx)
		_, err := client.Edit(ctx, &event.EditRequest{
			Id: 0,
			Title: func() *string {
				str := "Test title 2"

				return &str
			}(),
			Description: func() *string {
				str := "Test description 2"

				return &str
			}(),
		})
		require.NoError(t, err)

		ev, err := client.Get(ctx, &event.GetRequest{Id: 0})
		require.NoError(t, err)
		require.Equal(t, "Test title 2", ev.Event.GetTitle())
	})

	t.Run("get datetime range", func(t *testing.T) {
		ctx := context.TODO()

		client := getGRPCClient(ctx)
		response, err := client.GetDateTimeRange(ctx, &event.RangeRequest{
			StartDate: timestamppb.New(time.Now()),
			EndDate:   timestamppb.New(time.Now().Add(time.Hour * 5)),
		})
		require.NoError(t, err)
		require.Len(t, response.Events, 1)
		require.Equal(t, int64(0), response.Events[0].Id)
	})

	t.Run("delete", func(t *testing.T) {
		ctx := context.TODO()

		client := getGRPCClient(ctx)
		_, err := client.Remove(ctx, &event.RemoveRequest{Id: 0})
		require.NoError(t, err)

		_, err = client.Get(ctx, &event.GetRequest{Id: 0})
		require.Error(t, err)
	})
}

func getGRPCServer(ctx context.Context) (*internalgrpc.Server, error) {
	logg := logger.New("debug")
	storage := memorystorage.New()
	err := storage.Connect(ctx)
	if err != nil {
		return nil, err
	}

	application := app.New(logg, storage)
	server := internalgrpc.NewServer(logg, application, "localhost:9000")

	go func() {
		if err := server.Start(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	return server, nil
}

func getGRPCClient(ctx context.Context) event.EventServiceClient {
	cc, err := grpc.DialContext(
		ctx,
		"localhost:9000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	return event.NewEventServiceClient(cc)
}

package internalgrpc

import (
	"context"

	pb "github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/gen_buf/grpc/v1"
	"github.com/blancduman/otus-go/hw12_13_14_15_calendar/internal/server"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventServer struct {
	App    server.Application
	Logger server.Logger
	pb.UnimplementedEventServiceServer
}

func (s EventServer) Add(ctx context.Context, r *pb.AddRequest) (*pb.AddResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("ctx done")
	default:
		id, err := s.App.CreateEvent(
			ctx,
			r.GetTitle(),
			r.GetDescription(),
			r.GetOwnerId(),
			r.GetStartDate().AsTime(),
			r.GetEndDate().AsTime(),
			r.GetRemindIn().AsTime(),
		)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, errors.Wrap(err, "could not create event")
		}

		return &pb.AddResponse{Id: id}, nil
	}
}

func (s EventServer) Edit(ctx context.Context, r *pb.EditRequest) (*pb.EditResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("ctx done")
	default:
		err := s.App.EditEvent(
			ctx,
			r.GetId(),
			r.GetTitle(),
			r.GetDescription(),
			r.GetOwnerId(),
			r.GetStartDate().AsTime(),
			r.GetEndDate().AsTime(),
			r.GetRemindIn().AsTime(),
		)
		if err != nil {
			return nil, errors.Wrap(err, "could not edit event")
		}

		return &pb.EditResponse{}, nil
	}
}

func (s EventServer) Remove(ctx context.Context, r *pb.RemoveRequest) (*pb.RemoveResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("ctx done")
	default:
		err := s.App.RemoveEvent(ctx, r.GetId())
		if err != nil {
			return nil, errors.Wrap(err, "could not remove pb")
		}

		return &pb.RemoveResponse{}, nil
	}
}

func (s EventServer) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("ctx done")
	default:
		ev, err := s.App.GetEvent(ctx, r.GetId())
		if err != nil {
			return nil, errors.Wrap(err, "could not get event")
		}

		return &pb.GetResponse{Event: &pb.Event{
			Id:          ev.ID,
			Title:       ev.Title,
			Description: ev.Description,
			OwnerId:     ev.OwnerID,
			StartDate:   timestamppb.New(ev.StartDate),
			EndDate:     timestamppb.New(ev.EndDate),
			RemindIn:    timestamppb.New(ev.RemindIn),
		}}, nil
	}
}

func (s EventServer) GetDateTimeRange(ctx context.Context, r *pb.RangeRequest) (*pb.RangeResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("ctx done")
	default:
		evs, err := s.App.GetDateTimeRangeEvents(ctx, r.GetStartDate().AsTime(), r.GetEndDate().AsTime())
		if err != nil {
			return nil, errors.Wrap(err, "could not get datetime range events")
		}

		var events []*pb.Event

		for _, ev := range evs {
			events = append(events, &pb.Event{
				Id:          ev.ID,
				Title:       ev.Title,
				Description: ev.Description,
				OwnerId:     ev.OwnerID,
				StartDate:   timestamppb.New(ev.StartDate),
				EndDate:     timestamppb.New(ev.EndDate),
				RemindIn:    timestamppb.New(ev.RemindIn),
			})
		}

		return &pb.RangeResponse{Events: events}, nil
	}
}

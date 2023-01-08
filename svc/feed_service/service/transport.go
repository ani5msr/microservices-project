package service

import (
	"context"

	pb "github.com/ani5msr/microservices-project/pbu"
	om "github.com/ani5msr/microservices-project/pkg/object_model"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func newEvent(e *om.PostManagerEvent) (event *pb.Event) {
	event = &pb.Event{
		EventType: (pb.EventType)(e.EventType),
		Username:  e.Username,
		Url:       e.Url,
	}

	seconds := e.Timestamp.Unix()
	nanos := (int32(e.Timestamp.UnixNano() - 1e9*seconds))
	event.Timestamp = &timestamp.Timestamp{Seconds: seconds, Nanos: nanos}
	return
}

func decodeGetFeedRequest(_ context.Context, r interface{}) (interface{}, error) {
	request := r.(*pb.GetFeedRequest)
	return om.GetFeedRequest{
		Username:   request.Username,
		StartToken: request.StartToken,
	}, nil
}

func encodeGetFeedResponse(_ context.Context, r interface{}) (interface{}, error) {
	return r, nil
}

func makeGetFeedEndpoint(svc om.FeedManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(om.GetFeedRequest)
		r, err := svc.GetFeed(req)
		res := &pb.GetFeedResponse{
			Events:    []*pb.Event{},
			NextToken: r.NextToken,
		}
		if err != nil {
			res.Err = err.Error()
		}
		for _, e := range r.Events {
			event := newEvent(e)
			res.Events = append(res.Events, event)
		}
		return res, nil
	}
}

type handler struct {
	getFeed grpctransport.Handler
}

func (s *handler) GetFeed(ctx context.Context, r *pb.GetFeedRequest) (*pb.GetFeedResponse, error) {
	_, resp, err := s.getFeed.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.GetFeedResponse), nil
}

func newFeedServer(svc om.FeedManager) pb.FeedServer {
	return &handler{
		getFeed: grpctransport.NewServer(
			makeGetFeedEndpoint(svc),
			decodeGetFeedRequest,
			encodeGetFeedResponse,
		),
	}
}

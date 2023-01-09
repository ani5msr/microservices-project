package feed_manager_client

import (
	"context"
	"errors"
	"time"

	pb "github.com/ani5msr/microservices-project/pbu"
	om "github.com/ani5msr/microservices-project/pkg/object_model"
	"github.com/go-kit/kit/endpoint"
)

type EndpointSet struct {
	GetFeedEndpoint endpoint.Endpoint
}

func encodeGetFeedRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(om.GetFeedRequest)
	return &pb.GetFeedRequest{
		Username:   req.Username,
		StartToken: req.StartToken,
	}, nil
}

func newEvent(e *pb.Event) (event *om.PostManagerEvent) {
	return &om.PostManagerEvent{
		EventType: (om.PostManagerEventTypeEnum)(e.EventType),
		Username:  e.Username,
		Url:       e.Url,
		Timestamp: time.Unix(e.Timestamp.GetSeconds(), (int64)(e.Timestamp.GetNanos())),
	}
}

func decodeGetFeedResponse(_ context.Context, r interface{}) (interface{}, error) {
	gnr := r.(*pb.GetFeedResponse)
	if gnr.Err != "" {
		return nil, errors.New(gnr.Err)
	}

	res := &om.GetFeedResult{
		NextToken: gnr.NextToken,
	}

	for _, e := range gnr.Events {
		res.Events = append(res.Events, newEvent(e))
	}
	return res, nil
}

func (s EndpointSet) GetFeed(req om.GetFeedRequest) (result om.GetFeedResult, err error) {
	resp, err := s.GetFeedEndpoint(context.Background(), req)
	if err != nil {
		return
	}
	result = *(resp.(*om.GetFeedResult))
	return
}

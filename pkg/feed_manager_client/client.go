package feed_manager_client

import (
	pb "github.com/ani5msr/microservices-project/pbu"
	om "github.com/ani5msr/microservices-project/pkg/object_model"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type DisconnectFunc func() error

func NewClient(grpcAddr string) (cli om.FeedManager, disconnectFunc DisconnectFunc, err error) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	disconnectFunc = func() (err error) {
		if conn == nil {
			return
		}

		err = conn.Close()
		return
	}

	if err != nil {
		return
	}
	var getFeedEndpoint = grpctransport.NewClient(
		conn, "pb.Feed", "GetFeed",
		encodeGetFeedRequest,
		decodeGetFeedResponse,
		pb.GetFeedResponse{},
	).Endpoint()

	cli = EndpointSet{
		GetFeedEndpoint: getFeedEndpoint,
	}
	return
}

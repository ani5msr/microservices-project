package service

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/ani5msr/microservices-project/pbu"
	nm "github.com/ani5msr/microservices-project/pkg/feed_manager"
	"google.golang.org/grpc"
)

func Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "6060"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	redisHostname := os.Getenv("FEED_MANAGER_REDIS_SERVICE_HOST")
	redisPort := os.Getenv("FEED_MANAGER_REDIS_SERVICE_PORT")

	var store nm.Store
	if redisHostname == "" {
		store = nm.NewInMemoryFeedStore()
	} else {
		address := fmt.Sprintf("%s:%s", redisHostname, redisPort)
		store, err = nm.NewRedisFeedStore(address)
		if err != nil {
			log.Fatal(err)
		}
	}

	natsHostname := os.Getenv("NATS_CLUSTER_SERVICE_HOST")
	natsPort := os.Getenv("NATS_CLUSTER_SERVICE_PORT")

	svc, err := nm.NewFeedManager(store, natsHostname, natsPort)
	if err != nil {
		log.Fatal(err)
	}

	gRPCServer := grpc.NewServer()
	pb.RegisterFeedServer(gRPCServer, newFeedServer(svc))

	fmt.Printf("News service is listening on port %s...\n", port)
	err = gRPCServer.Serve(listener)
	fmt.Println("Serve() failed", err)
}

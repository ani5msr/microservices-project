package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	jeagerconfig "github.com/uber/jaeger-client-go/config"

	"github.com/ani5msr/microservices-project/pkg/db_utils"
	"github.com/ani5msr/microservices-project/pkg/log"
	om "github.com/ani5msr/microservices-project/pkg/object_model"
	lm "github.com/ani5msr/microservices-project/pkg/post_manager"
	"github.com/ani5msr/microservices-project/pkg/post_manager_events"
	sgm "github.com/ani5msr/microservices-project/pkg/social_graph_client"
)

type EventSink struct {
}

type postManagerMiddleware func(om.PostManager) om.PostManager

func (s *EventSink) OnPostAdded(username string, post *om.Post) {
	//log.Println("Post added")
}

func (s *EventSink) OnPostUpdated(username string, post *om.Post) {
	//log.Println("Post updated")
}

func (s *EventSink) OnPostDeleted(username string, url string) {
	//log.Println("Post deleted")
}

// createTracer returns an instance of Jaeger Tracer that samples
// 100% of traces and logs all spans to stdout.
func createTracer(service string) (opentracing.Tracer, io.Closer) {
	cfg := &jeagerconfig.Configuration{
		ServiceName: service,
		Sampler: &jeagerconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jeagerconfig.ReporterConfig{
			LogSpans: true,
		},
	}
	logger := jeagerconfig.Logger(jaeger.StdLogger)
	tracer, closer, err := cfg.NewTracer(logger)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot create tracer: %v\n", err))
	}
	return tracer, closer
}

func Run() {
	dbHost, dbPort, err := db_utils.GetDbEndpoint("post")
	if err != nil {
		log.Fatal(err)
	}
	store, err := lm.NewDbPostStore(dbHost, dbPort, "postgres", "postgres")
	if err != nil {
		log.Fatal(err)
	}

	sgHost := os.Getenv("SOCIAL_GRAPH_MANAGER_SERVICE_HOST")
	if sgHost == "" {
		sgHost = "localhost"
	}

	sgPort := os.Getenv("SOCIAL_GRAPH_MANAGER_SERVICE_PORT")
	if sgPort == "" {
		sgPort = "9090"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	maxPostsPerUserStr := os.Getenv("MAX_POSTS_PER_USER")
	if maxPostsPerUserStr == "" {
		maxPostsPerUserStr = "10"
	}

	maxPostsPerUser, err := strconv.ParseInt(maxPostsPerUserStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	socialGraphClient, err := sgm.NewClient(fmt.Sprintf("%s:%s", sgHost, sgPort))
	if err != nil {
		log.Fatal(err)
	}

	natsHostname := os.Getenv("NATS_CLUSTER_SERVICE_HOST")
	natsPort := os.Getenv("NATS_CLUSTER_SERVICE_PORT")

	natsUrl := ""
	var eventSink om.PostManagerEvents
	if natsHostname != "" {
		natsUrl = natsHostname + ":" + natsPort
		eventSink, err = post_manager_events.NewEventSender(natsUrl)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		eventSink = &EventSink{}
	}

	// Create a logger
	logger := log.NewLogger("post manager")

	// Create a tracer
	tracer, closer := createTracer("post-manager")
	defer closer.Close()

	// Create the service implementation
	svc, err := lm.NewPostManager(store, socialGraphClient, natsUrl, eventSink, maxPostsPerUser)
	if err != nil {
		log.Fatal(err)
	}

	// Hook up the logging middleware to the service and the logger
	svc = newLoggingMiddleware(logger)(svc)

	// Hook up the metrics middleware
	svc = newMetricsMiddleware()(svc)

	// Hook up the tracing middleware
	svc = newTracingMiddleware(tracer)(svc)

	getPostsHandler := httptransport.NewServer(
		makeGetPostEndpoint(svc),
		decodeGetPostRequest,
		encodeResponse,
	)

	addPostHandler := httptransport.NewServer(
		makeAddPostEndpoint(svc),
		decodeAddPostRequest,
		encodeResponse,
	)

	updatePostHandler := httptransport.NewServer(
		makeUpdatePostEndpoint(svc),
		decodeUpdatePostRequest,
		encodeResponse,
	)

	deletePostHandler := httptransport.NewServer(
		makeDeletePostEndpoint(svc),
		decodeDeletePostRequest,
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Methods("GET").Path("/posts").Handler(getPostsHandler)
	r.Methods("POST").Path("/posts").Handler(addPostHandler)
	r.Methods("PUT").Path("/posts").Handler(updatePostHandler)
	r.Methods("DELETE").Path("/posts").Handler(deletePostHandler)
	r.Methods("GET").Path("/metrics").Handler(promhttp.Handler())

	logger.Log("msg", "*** listening on ***", "port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

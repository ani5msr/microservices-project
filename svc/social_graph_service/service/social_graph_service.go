package service

import (
	"log"
	"net/http"
	"os"

	"github.com/ani5msr/microservices-project/pkg/db_utils"
	sgm "github.com/ani5msr/microservices-project/pkg/social_graph_manager"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Run() {
	log.Println("Service started...")
	dbHost, dbPort, err := db_utils.GetDbEndpoint("social_graph")
	if err != nil {
		log.Fatal(err)
	}
	store, err := sgm.NewDbSocialGraphStore(dbHost, dbPort, "postgres", "postgres")
	check(err)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	svc, err := sgm.NewSocialGraphManager(store)
	check(err)

	followHandler := httptransport.NewServer(
		makeFollowEndpoint(svc),
		decodeFollowRequest,
		encodeResponse,
	)

	unfollowHandler := httptransport.NewServer(
		makeUnfollowEndpoint(svc),
		decodeUnfollowRequest,
		encodeResponse,
	)

	getFollowingHandler := httptransport.NewServer(
		makeGetFollowingEndpoint(svc),
		decodeGetFollowingRequest,
		encodeResponse,
	)

	getFollowersHandler := httptransport.NewServer(
		makeGetFollowersEndpoint(svc),
		decodeGetFollowersRequest,
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Methods("POST").Path("/follow").Handler(followHandler)
	r.Methods("POST").Path("/unfollow").Handler(unfollowHandler)
	r.Methods("GET").Path("/following/{username}").Handler(getFollowingHandler)
	r.Methods("GET").Path("/followers/{username}").Handler(getFollowersHandler)

	log.Printf("Listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

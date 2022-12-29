package service

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/the-gigi/delinkcious/pkg/db_util"

	sgm "github.com/ani5msr/microservices-project/pkg/user_manager"
	httptransport "github.com/go-kit/kit/transport/http"
)

func Run() {
	dbHost, dbPort, err := db_util.GetDbEndpoint("user")
	if err != nil {
		log.Fatal(err)
	}
	store, err := sgm.NewDbUserStore(dbHost, dbPort, "postgres", "postgres")
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "7070"
	}

	svc, err := sgm.NewUserManager(store)
	if err != nil {
		log.Fatal(err)
	}

	registerHandler := httptransport.NewServer(
		makeRegisterEndpoint(svc),
		decodeRegisterRequest,
		encodeResponse,
	)

	loginHandler := httptransport.NewServer(
		makeLoginEndpoint(svc),
		decodeLoginRequest,
		encodeResponse,
	)

	logoutHandler := httptransport.NewServer(
		makeLogoutEndpoint(svc),
		decodeLogoutRequest,
		encodeResponse,
	)

	r := mux.NewRouter()
	r.Methods("POST").Path("/register").Handler(registerHandler)
	r.Methods("POST").Path("/login").Handler(loginHandler)
	r.Methods("POST").Path("/logout").Handler(logoutHandler)

	log.Printf("Listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

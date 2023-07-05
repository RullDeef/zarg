package main

import (
	"fmt"
	"net/http"
	"server/cmd/compaign/repository"
	"server/cmd/lobby/handler"
	"server/cmd/lobby/service"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const defaultPort = 4667

func main() {
	log := logrus.New()
	log.SetReportCaller(true)
	logger := log.WithField("microservice", "lobby")

	profileRepo := repository.NewProfile(logger)
	participator := service.NewParticipator(logger)

	router := mux.NewRouter()
	profileRouter := router.PathPrefix("/profiles").Subrouter()
	profileRouter.HandleFunc("/new", handler.BuildProfileCreateAnonymous(logger, profileRepo)).Methods("POST")
	profileRouter.HandleFunc("/{profile_id}", handler.BuildProfileGet(logger, profileRepo)).Methods("GET")

	router.HandleFunc("/compaigns/new", handler.BuildOnNewCompaignRequest(logger, participator))

	http.Handle("/", router)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", defaultPort), nil); err != nil {
		logger.Error(err)
		panic(err)
	}
}

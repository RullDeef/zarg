package main

import (
	"fmt"
	"net/http"
	"server/cmd/compaign/handler"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const defaultPort = 4668

var logger *logrus.Entry

func main() {
	log := logrus.New()
	log.SetReportCaller(true)
	logger = log.WithField("microservice", "compaign")

	router := mux.NewRouter()
	router.Path("/compaign/{id}").HandlerFunc(handler.GetCompaign).Methods("GET")

	http.Handle("/", router)

	err := http.ListenAndServe(fmt.Sprintf(":%d", defaultPort), nil)
	if err != nil {
		logger.Error(err)
		panic(err)
	}
}

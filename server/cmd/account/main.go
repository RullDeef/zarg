package main

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

const defaultPort = 4666

func main() {
	log := logrus.New()
	log.SetReportCaller(true)
	logger := log.WithField("microservice", "accounter")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", defaultPort), nil); err != nil {
		logger.Error(err)
		panic(err)
	}
}

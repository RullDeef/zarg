package server

import (
	"net/http"

	"go.uber.org/zap"
)

type profileMux struct {
	http.ServeMux

	logger *zap.SugaredLogger
}

func newProfileMux(logger *zap.SugaredLogger) *profileMux {
	pm := profileMux{
		ServeMux: *http.NewServeMux(),
		logger:   logger,
	}

	return &pm
}

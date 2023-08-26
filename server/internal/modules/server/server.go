package server

import (
	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	logger *zap.SugaredLogger
	mux    *http.ServeMux
}

func NewServer(logger *zap.SugaredLogger, pm *profileMux, lm *lobbyMux) *Server {
	s := Server{
		logger: logger,
		mux:    http.NewServeMux(),
	}

	// /api/lobby/join?mode=[single|guild|random]
	s.mux.Handle("/api/lobby/", http.StripPrefix("/api/lobby", lm))

	// /api/profiles/new
	// /api/profiles/[id]
	s.mux.Handle("/api/profiles/", http.StripPrefix("/api/profiles", pm))

	return &s
}

func (s *Server) Run(address string) error {
	s.logger.Infow("Run", "address", address)

	return http.ListenAndServe(address, s.mux)
}

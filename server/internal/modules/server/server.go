package server

import (
	"net/http"
	"server/internal/modules/lobby"

	"go.uber.org/zap"
)

type Server struct {
	logger *zap.SugaredLogger
	mux    *http.ServeMux
	lobby  *lobby.Lobby
}

func New(
	logger *zap.SugaredLogger,
	lobby *lobby.Lobby,
) *Server {
	s := Server{
		logger: logger,
		mux:    http.NewServeMux(),
		lobby:  lobby,
	}

	// /api/lobby/join?mode=[single|guild|random]
	s.mux.Handle("/api/lobby/join",
		panicWrapperMiddleware(
			loggerMiddleware(
				s.logger,
				http.HandlerFunc(s.apiLobbyJoin),
			),
		),
	)

	return &s
}

func (s *Server) Run(address string) error {
	return http.ListenAndServe(address, s.mux)
}

func (s *Server) apiLobbyJoin(w http.ResponseWriter, r *http.Request) {
	// check auth token
	authToken := r.Header.Get("Authorization")

	profile, ok := s.authManager.ValidateToken(authToken)
	if authToken != "" && !ok {
		// token is invalid but not empty
	}

	// check mode
	mode := r.FormValue("mode")

}

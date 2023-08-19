package server

import (
	"context"
	"net/http"
	"server/internal/modules/lobby"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Server struct {
	logger     *zap.SugaredLogger
	mux        *http.ServeMux
	controller *lobby.Controller
}

func New(logger *zap.SugaredLogger, controller *lobby.Controller) *Server {
	s := Server{
		logger:     logger,
		mux:        http.NewServeMux(),
		controller: controller,
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
	s.logger.Infow("Run", "address", address)

	return http.ListenAndServe(address, s.mux)
}

func (s *Server) apiLobbyJoin(w http.ResponseWriter, r *http.Request) {
	mode := r.FormValue("mode")
	authToken := r.Header.Get("Authorization")
	profileID, err := s.controller.AcceptJoinRequest(mode, authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	mw := wsUpgradeMiddleware(func(conn *websocket.Conn) {
		conn.SetCloseHandler(func(code int, text string) error {
			s.controller.CancelRequest(profileID)
			return nil
		})

		conn.WriteJSON(map[string]string{
			"profile_id": string(profileID),
			"status":     "waiting",
		})

		compaignID, err := s.controller.WaitJoin(context.Background(), profileID)
		if err != nil {
			conn.WriteJSON(map[string]string{"error": err.Error()})
		} else {
			conn.WriteJSON(map[string]string{"compaign_id": string(compaignID)})
		}

		conn.Close()
	})

	go mw.ServeHTTP(w, r)
}

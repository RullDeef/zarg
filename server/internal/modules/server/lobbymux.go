package server

import (
	"context"
	"net/http"
	"server/internal/modules/lobby"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type lobbyMux struct {
	http.ServeMux

	logger     *zap.SugaredLogger
	controller *lobby.Controller
}

func newLobbyMux(logger *zap.SugaredLogger, controller *lobby.Controller) *lobbyMux {
	lm := lobbyMux{
		ServeMux:   *http.NewServeMux(),
		logger:     logger,
		controller: controller,
	}

	// /api/lobby/join?mode=[single|guild|random]
	lm.Handle("/join",
		panicWrapperMiddleware(
			loggerMiddleware(
				lm.logger,
				http.HandlerFunc(lm.apiLobbyJoin),
			),
		),
	)

	return &lm
}

func (lm *lobbyMux) apiLobbyJoin(w http.ResponseWriter, r *http.Request) {
	mode := r.FormValue("mode")
	authToken := r.Header.Get("Authorization")
	profileID, err := lm.controller.AcceptJoinRequest(mode, authToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	mw := wsUpgradeMiddleware(func(conn *websocket.Conn) {
		conn.SetCloseHandler(func(code int, text string) error {
			lm.controller.CancelRequest(profileID)
			return nil
		})

		conn.WriteJSON(map[string]string{
			"profile_id": string(profileID),
			"status":     "waiting",
		})

		compaignID, err := lm.controller.WaitJoin(context.Background(), profileID)
		if err != lobby.ErrRequestCancelled {
			if err != nil {
				conn.WriteJSON(map[string]string{"error": err.Error()})
			} else {
				conn.WriteJSON(map[string]string{"compaign_id": string(compaignID)})
			}
			conn.Close()
		}
	})

	go mw.ServeHTTP(w, r)
}

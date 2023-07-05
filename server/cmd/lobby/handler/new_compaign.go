package handler

import (
	"net/http"
	"server/cmd/lobby/service"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func BuildOnNewCompaignRequest(logger *logrus.Entry, participator *service.Participator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.WithField("r", r)
		logger.Info("onNewCompaignRequest")

		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error(err)
			return
		}

		var request service.ParticipationRequest
		if err := conn.ReadJSON(&request); err != nil {
			logger.Error(err)
			return
		}

		if err := participator.SubmitRequest(conn, request); err != nil {
			conn.WriteJSON(map[string]string{"error": err.Error()})
			logger.Error(err)
			conn.Close()
			return
		}
	}
}

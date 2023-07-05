package service

import (
	"container/list"
	"errors"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ParticipationMode string

const (
	participationSingle = ParticipationMode("single")
	participationRandom = ParticipationMode("random")
	participationGuild  = ParticipationMode("guild")
)

var (
	errInvalidMode    = errors.New("bad mode value")
	errModeNotAllowed = errors.New("mode not allowed")
)

type Participator struct {
	logger *logrus.Entry

	referees map[ParticipationMode]*Referee
}

// структура запроса на участие в походе
type ParticipationRequest struct {
	Mode      ParticipationMode `json:"mode"`
	Anonymous bool              `json:"anonymous"`
	ProfileID string            `json:"profile_id"`
	JWT       string            `json:"jwt"`
}

// внутреннее представление участника похода
type participant struct {
	ws       *websocket.Conn
	req      ParticipationRequest
	playerID string
}

func NewParticipator(logger *logrus.Entry) *Participator {
	p := &Participator{
		referees: make(map[ParticipationMode]*Referee),
	}

	p.referees[participationSingle] = NewReferee(logger, RefereeConfig{
		MinPlayersInTeam: 1,
		MaxPlayersInTeam: 1,
	}, p.onTeamReady)

	p.referees[participationRandom] = NewReferee(logger, RefereeConfig{
		MinPlayersInTeam: 2,
		MaxPlayersInTeam: 4,
	}, p.onTeamReady)

	p.referees[participationGuild] = NewReferee(logger, RefereeConfig{
		MinPlayersInTeam: 2,
		MaxPlayersInTeam: 4,
	}, p.onTeamReady)

	p.logger = logger.WithField("service", p)
	return p
}

func (p *Participator) SubmitRequest(ws *websocket.Conn, req ParticipationRequest) error {
	logger := p.logger.WithField("req", req)

	if !modeIsValid(req.Mode) {
		return errInvalidMode
	}

	// ожидаем, что игрок будет анонимным пока что
	if !req.Anonymous {
		logger.Error(errModeNotAllowed)
		return errModeNotAllowed
	}

	partis := &participant{
		ws:       ws,
		req:      req,
		playerID: uuid.New().String(),
	}
	logger.WithField("partis", partis).Info("new participant")

	p.referees[req.Mode].EnqueuePlayer(partis)
	return nil
}

func (p *Participator) onTeamReady(team *list.List) {
	logger := p.logger.WithField("team", team)
	logger.Info("onTeamReady")

	// создать новый поход
	compaignID := uuid.New().String()

	playerIDs := make([]string, 0, team.Len())
	for node := team.Front(); node != nil; node = node.Next() {
		participant := node.Value.(*participant)
		playerIDs = append(playerIDs, participant.playerID)
	}

	message := map[string]any{
		"compaign_id": compaignID,
		"players":     playerIDs,
	}

	// отправить каждому участнику ответ с успехом и разорвать соединение
	for node := team.Front(); node != nil; node = node.Next() {
		participant := node.Value.(*participant)
		participant.ws.WriteJSON(message)

		go func(ws *websocket.Conn) {
			ws.ReadMessage()
			ws.Close()
		}(participant.ws)
	}
}

func modeIsValid(mode ParticipationMode) bool {
	if mode == participationSingle {
		return true
	} else if mode == participationRandom {
		return true
	} else if mode == participationGuild {
		return true
	}
	return false
}

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

type CompaignFactoryFunc func([]ParticipationRequest) (map[string]any, error)

type Participator struct {
	logger *logrus.Entry

	referees        map[ParticipationMode]*Referee
	compaignFactory CompaignFactoryFunc
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

func NewParticipator(logger *logrus.Entry, factory CompaignFactoryFunc) *Participator {
	p := &Participator{
		referees:        make(map[ParticipationMode]*Referee),
		compaignFactory: factory,
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

	p.logger = logger //.WithField("service", p)
	return p
}

func (p *Participator) SubmitRequest(ws *websocket.Conn, req ParticipationRequest) error {
	logger := p.logger.WithField("req", req)

	if !modeIsValid(req.Mode) {
		return errInvalidMode
	}

	// TODO: ожидаем, что игрок будет анонимным пока что
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

	if p.playerInAnyQueue(req.ProfileID) {
		logger.WithField("partis", partis).Error("already in some queue")
	} else {
		p.referees[req.Mode].EnqueuePlayer(partis)
	}
	return nil
}

func (p *Participator) CancelRequest(ws *websocket.Conn) {
	p.logger.Info("CancelRequest")

	for _, referee := range p.referees {
		referee.RemovePlayerIf(func(player any) bool {
			return player.(*participant).ws == ws
		})
	}
}

func (p *Participator) onTeamReady(team *list.List) {
	logger := p.logger.WithField("team", team)
	logger.Info("onTeamReady")

	// создать новый поход
	reqs := extractPlayerRequests(team)
	compaignInfo, err := p.compaignFactory(reqs)
	if err != nil {
		logger.Error(err)
		return
	}

	// отправить каждому участнику ответ с успехом и разорвать соединение
	for node := team.Front(); node != nil; node = node.Next() {
		participant := node.Value.(*participant)
		participant.ws.WriteJSON(compaignInfo)

		go func(ws *websocket.Conn) {
			ws.ReadMessage()
			ws.Close()
		}(participant.ws)
	}
}

func (p *Participator) playerInAnyQueue(profileID string) bool {
	for _, referee := range p.referees {
		if referee.HasPlayer(func(player any) bool {
			return player.(*participant).req.ProfileID == profileID
		}) {
			return true
		}
	}

	return false
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

func extractPlayerRequests(team *list.List) []ParticipationRequest {
	reqs := make([]ParticipationRequest, 0, team.Len())
	for node := team.Front(); node != nil; node = node.Next() {
		reqs = append(reqs, node.Value.(*participant).req)
	}
	return reqs
}

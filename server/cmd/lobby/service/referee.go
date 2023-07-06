package service

import (
	"container/list"
	"sync"

	"github.com/sirupsen/logrus"
)

type TeamReadyFunc func(*list.List)

type RefereeConfig struct {
	MinPlayersInTeam int
	MaxPlayersInTeam int
}

type Referee struct {
	logger *logrus.Entry
	config RefereeConfig

	waitingPlayers *list.List
	onTeamReady    TeamReadyFunc
	mutex          sync.Mutex
}

func NewReferee(logger *logrus.Entry, config RefereeConfig, onTeamReady TeamReadyFunc) *Referee {
	service := &Referee{
		config:         config,
		waitingPlayers: list.New(),
		onTeamReady:    onTeamReady,
	}
	service.logger = logger.WithFields(logrus.Fields{
		"service": service,
	})
	return service
}

func (service *Referee) EnqueuePlayer(player any) {
	service.logger.WithField("player", player).Info("EnqueuePlayer")

	service.mutex.Lock()
	defer service.mutex.Unlock()

	service.waitingPlayers.PushBack(player)

	if service.waitingPlayers.Len() >= service.config.MinPlayersInTeam {
		go service.extractTeam()
	}
}

func (service *Referee) RemovePlayer(player any) {
	logger := service.logger.WithField("player", player)
	logger.Info("RemovePlayer")

	service.mutex.Lock()
	defer service.mutex.Unlock()

	elem := service.waitingPlayers.Front()
	for elem != nil && elem.Value != player {
		elem = elem.Next()
	}

	if elem != nil {
		service.waitingPlayers.Remove(elem)
	} else {
		logger.Error("player not found in waiting list")
	}
}

func (service *Referee) RemovePlayerIf(predicate func(player any) bool) bool {
	service.logger.Info("RemovePlayerIf")

	service.mutex.Lock()
	defer service.mutex.Unlock()

	elem := service.waitingPlayers.Front()
	for elem != nil {
		if predicate(elem.Value) {
			service.waitingPlayers.Remove(elem)
			return true
		}
		elem = elem.Next()
	}

	service.logger.Info("player not found by predicate")
	return false
}

func (service *Referee) HasPlayer(predicate func(player any) bool) bool {
	service.logger.Info("HasPlayer")

	service.mutex.Lock()
	defer service.mutex.Unlock()

	elem := service.waitingPlayers.Front()
	for elem != nil {
		if predicate(elem.Value) {
			return true
		}
		elem = elem.Next()
	}

	service.logger.Info("player not found by predicate")
	return false
}

// Выделяет группу игроков для похода в подземелье
//
// Алгоритм использует обычную очередь
func (service *Referee) extractTeam() {
	service.logger.Info("extractTeam")

	service.mutex.Lock()
	defer service.mutex.Unlock()

	if service.waitingPlayers.Len() < service.config.MinPlayersInTeam {
		service.logger.WithField("waitingLen", service.waitingPlayers.Len()).Error("too few players to make a team")
		return
	}

	team := list.New()
	for service.waitingPlayers.Len() > 0 && team.Len() < service.config.MaxPlayersInTeam {
		elem := service.waitingPlayers.Front()
		team.PushBack(elem.Value)
		service.waitingPlayers.Remove(elem)
	}

	go service.onTeamReady(team)
}

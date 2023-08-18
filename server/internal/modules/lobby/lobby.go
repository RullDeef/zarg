package lobby

import (
	"errors"
	"server/domain"
	"time"

	"go.uber.org/zap"
)

var ErrParticipationModeInvalid = errors.New("invalid participation mode") // некорректный режим похода

// Lobby - структура, отвечающая за обработку запросов на участие в походах и формирование групп
type Lobby struct {
	logger        *zap.SugaredLogger
	requestQueues map[ParticipationMode]*requestQueue // очереди запросов на участие в походе
}

type ParticipationMode string

const (
	ParticipationSingle ParticipationMode = "single" // поход для одного игрока с ботами
	ParticipationGuild  ParticipationMode = "guild"  // поход для игроков одной гильдии
	ParticipationRandom ParticipationMode = "random" // поход для случайно подобранных игроков
)

// New - создает новое лобби
func New(logger *zap.SugaredLogger) *Lobby {
	return &Lobby{
		logger: logger,
		requestQueues: map[ParticipationMode]*requestQueue{
			ParticipationSingle: newQueue(singlePartyBuilder, 5*time.Second),
			ParticipationGuild:  newQueue(guildPartyBuilder, 5*time.Second),
			ParticipationRandom: newQueue(randomPartyBuilder, 5*time.Second),
		},
	}
}

// Close - завершает процесс формирования команд во всех очередях
func (l *Lobby) Close() (err error) {
	l.logger.Info("Close")

	for _, queue := range l.requestQueues {
		err = errors.Join(err, queue.Close())
	}
	return
}

// AddRequest - добавляет запрос на участие в походе
func (l *Lobby) AddRequest(mode ParticipationMode, profile *domain.Profile) error {
	l.logger.Infow("AddRequest", "mode", mode, "profile", profile)

	if l.hasProfileInQueues(profile) {
		return ErrProfileAlreadyInQueue
	}

	if queue, ok := l.requestQueues[mode]; !ok {
		return ErrParticipationModeInvalid
	} else {
		queue.AddRequest(profile)
	}

	return nil
}

func (l *Lobby) hasProfileInQueues(profile *domain.Profile) bool {
	for _, queue := range l.requestQueues {
		if queue.hasProfile(profile) {
			return true
		}
	}
	return false
}

// CancelRequest - отмена запроса на участие в походе
func (l *Lobby) CancelRequest(profile *domain.Profile) {
	l.logger.Infow("CancelRequest", "profile", profile)

	for _, queue := range l.requestQueues {
		queue.CancelRequest(profile)
	}
}

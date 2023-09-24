package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNoParticipants = errors.New("no participants") // нет участников
	ErrNoItems        = errors.New("no items")        // нет предметов
	ErrTimeoutExpired = errors.New("timeout expired") // время распределения вышло
)

// Distributor - абстракция универсального "распределителя" вещей
// между игроками, который может учитывать предпочтения самих игроков.
// Его задача - определить соответствие между предметами и игроками
type Distributor interface {
	Distribute(context.Context, []*Player, []*PickableItem, DistributionConfig) (PlayerItemDistribution, error)
}

// DistributionConfig - структура, представляющая конфигурацию распределения
type DistributionConfig struct {
	// Timeout - ограничение по времени на распределение.
	// По истечению данного времени и(!) отсутствии корректного
	// распределения необходимо вернуть ошибку ErrTimeoutExpired.
	// В отсутствии временных ограничений Timeout = 0
	Timeout time.Duration

	MaxItemsPerPlayer     int  // максимальное число предметов, которые может выбрать один игрок
	MultipleOwnersAllowed bool // могут ли несколько игроков выбрать один и тот же предмет
}

// PlayerItemDistribution - тип распределения предметов между игроками
type PlayerItemDistribution map[*Player][]*PickableItem

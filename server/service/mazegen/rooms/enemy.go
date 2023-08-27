package rooms

import (
	"context"
	"math/rand"
	"server/domain"
)

// EnemyRoom - комната с врагами. При входе запускается бой
type EnemyRoom struct {
	enemies []*domain.Entity // список врагов в комнате

	randSource rand.Source // источник для генератора псевдослучайных чисел
}

var _ domain.Room = &EnemyRoom{}

// NewEnemyRoom - создает новую комнату с врагами
func NewEnemyRoom(enemies []*domain.Entity, randSource rand.Source) *EnemyRoom {
	return &EnemyRoom{
		enemies:    enemies,
		randSource: randSource,
	}
}

// Visit - запускает бой с врагами в комнате
func (e *EnemyRoom) Visit(ctx context.Context, c *domain.Compaign) error {
	fight, err := domain.NewFightRandomOrder(e.randSource, c.Participators, e.enemies)
	if err == nil {
		err = fight.PerformFight(ctx)
	}
	return err
}

package rooms

import (
	"context"
	"math/rand"
	"server/domain"
)

type EnemyRoom struct {
	enemies []*domain.Entity

	randSource rand.Source
}

var _ domain.Room = &EnemyRoom{}

func NewEnemyRoom(enemies []*domain.Entity, randSource rand.Source) *EnemyRoom {
	return &EnemyRoom{
		enemies: enemies,
	}
}

func (e *EnemyRoom) Visit(ctx context.Context, c *domain.Compaign) error {
	players := make([]domain.Fightable, len(c.Participators))
	for i, p := range c.Participators {
		players[i] = p
	}
	enemies := make([]domain.Fightable, len(e.enemies))
	for i, e := range e.enemies {
		enemies[i] = e
	}

	fight, err := domain.NewFightRandomOrder(e.randSource, players, enemies)
	if err == nil {
		err = fight.PerformFight(ctx)
	}
	return err
}

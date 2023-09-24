package rooms

import (
	"context"
	"log"
	"math/rand"
	"server/domain"
)

// BossRoom - комната с боссом
type BossRoom struct {
	bosses     []*domain.Entity
	randSource rand.Source
}

var _ domain.Room = &BossRoom{}

func NewBossRoom(bosses []*domain.Entity, randSource rand.Source) *BossRoom {
	return &BossRoom{
		bosses:     bosses,
		randSource: randSource,
	}
}

func (br *BossRoom) Visit(ctx context.Context, c *domain.Compaign) error {
	fight, err := domain.NewFightRandomOrder(br.randSource, c.Participators, br.bosses)
	if err == nil {
		err = fight.PerformFight(ctx)
	}
	return err
}

type BossRoomGenerator struct {
	logger *log.Logger
}

func NewBossRoomGenerator() *BossRoomGenerator {
	// TODO: implement
	return &BossRoomGenerator{
		logger: log.New(log.Writer(), "boss-room", 0),
	}
}

func (brg *BossRoomGenerator) Generate(randSource rand.Source) (domain.Room, error) {
	// TODO: implement
	brg.logger.Println("boss room generated")
	return NewBossRoom(nil, randSource), nil
}

package rooms

import (
	"context"
	"log"
	"math/rand"
	"server/domain"
)

type TrapRoom struct {
	logger *log.Logger
}

func (t *TrapRoom) Visit(ctx context.Context, c *domain.Compaign) error {
	// TODO: implement
	t.logger.Printf("trap room visited")
	return nil
}

type TrapRoomGenerator struct {
	logger *log.Logger
}

func NewTrapRoomGenerator() *TrapRoomGenerator {
	// TODO: implement
	return &TrapRoomGenerator{
		logger: log.New(log.Writer(), "trap-room", 0),
	}
}

func (t *TrapRoomGenerator) Generate(randSource rand.Source) (domain.Room, error) {
	// TODO: implement
	t.logger.Printf("trap room generated")
	return &TrapRoom{
		logger: t.logger,
	}, nil
}

package rooms

import (
	"context"
	"log"
	"math/rand"
	"server/domain"
)

type NPCRoom struct {
	logger *log.Logger
}

type NPCRoomGenerator struct {
	logger *log.Logger
}

func NewNPCRoom(logger *log.Logger) *NPCRoom {
	return &NPCRoom{
		logger: logger,
	}
}

func (n *NPCRoom) Visit(ctx context.Context, c *domain.Compaign) error {
	// TODO: implement
	n.logger.Println("npc room visited")
	return nil
}

func NewNPCRoomGenerator() *NPCRoomGenerator {
	// TODO: implement
	return &NPCRoomGenerator{
		logger: log.New(log.Writer(), "npc-room", 0),
	}
}

func (n *NPCRoomGenerator) Generate(randSource rand.Source) (domain.Room, error) {
	// TODO: implement
	n.logger.Println("npc room generated")
	return NewNPCRoom(n.logger), nil
}

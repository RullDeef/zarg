package rooms

import (
	"context"
	"log"
	"math/rand"
	"server/domain"
)

type RestRoom struct {
	logger *log.Logger
}

func NewRestRoom(logger *log.Logger) *RestRoom {
	return &RestRoom{
		logger: logger,
	}
}

func (r *RestRoom) Visit(ctx context.Context, c *domain.Compaign) error {
	// TODO: implement
	r.logger.Println("rest room visited")
	return nil
}

type RestRoomGenerator struct {
	logger *log.Logger
}

func NewRestRoomGenerator() *RestRoomGenerator {
	// TODO: implement
	return &RestRoomGenerator{
		logger: log.New(log.Writer(), "rest-room", 0),
	}
}

func (r *RestRoomGenerator) Generate(randSource rand.Source) (domain.Room, error) {
	// TODO: implement
	r.logger.Println("rest room generated")
	return &RestRoom{
		logger: r.logger,
	}, nil
}

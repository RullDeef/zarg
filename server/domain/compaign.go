package domain

import (
	"context"
	"errors"
	"time"
)

var ErrCompaignCannotBeReused = errors.New("compaign cannot be reused") // попытка повторного похода

// Compaign - структура представляющая активный поход
type Compaign struct {
	ID            CompaignID
	Participators []*Player     // участники похода с учетом строя
	StartTime     time.Time     // время начала похода
	Duration      time.Duration // длительность похода (0, если он еще не завершен)

	// wayChooser - функция выбора направления (при перемещении)
	wayChooser func(context.Context, []DungeonWay) (DungeonWay, error)
}

type CompaignID string

type WayChooserFunc func(context.Context, []DungeonWay) (DungeonWay, error)

// NewCompaign - создает новый поход по списку профилей с учетом строя
func NewCompaign(participators []*Profile, wayChooser WayChooserFunc) *Compaign {
	players := make([]*Player, len(participators))
	for i, p := range participators {
		players[i] = NewPlayer(p)
	}
	return &Compaign{
		Participators: players,
		wayChooser:    wayChooser,
	}
}

// VisitDungeon - выполняет поход по подземелью
func (c *Compaign) VisitDungeon(ctx context.Context, d Dungeon) error {
	if c.Duration != 0 {
		return ErrCompaignCannotBeReused
	}
	c.StartTime = time.Now()
	room := d.Entrance()
	for !d.IsFinalRoom(room) {
		if err := room.Visit(ctx, c); err != nil {
			return err
		}

		ways := d.NextWays(room)
		if way, err := c.wayChooser(ctx, ways); err != nil {
			return err
		} else {
			room = way.Room
		}
	}
	c.Duration = time.Since(c.StartTime)
	return nil
}

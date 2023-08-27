package mazegen

import (
	"server/domain"
)

// RegularDungeon - обычное подземелье в несколько этажей
type RegularDungeon struct {
	floors [][]domain.Room
}

var _ domain.Dungeon = &RegularDungeon{}

func (d *RegularDungeon) Entrance() domain.Room {
	return d.floors[0][0]
}

func (d *RegularDungeon) IsFinalRoom(room domain.Room) bool {
	lastFloor := d.floors[len(d.floors)-1]
	lastRoom := lastFloor[len(lastFloor)-1]
	return room == lastRoom
}

func (d *RegularDungeon) NextWays(room domain.Room) []domain.DungeonWay {
	for i, floor := range d.floors {
		for j, r := range floor {
			if r == room {
				var nextRoom domain.Room
				if j == len(floor)-1 {
					nextRoom = d.floors[i+1][0]
				} else {
					nextRoom = floor[j+1]
				}
				return []domain.DungeonWay{
					{
						Name:          "Темный корридор",
						VisitedBefore: false,
						Room:          nextRoom,
					},
				}
			}
		}
	}
	panic("must never happen")
}

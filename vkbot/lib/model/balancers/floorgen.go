package balancers

import (
	"log"
	"math/rand"
	I "zarg/lib/model/interfaces"
	"zarg/lib/model/player/squad"
)

type FloorGenBalancer struct {
	floorNumber int // starts from 1
	players     *squad.PlayerSquad
}

func NewFloorGenBalancer(floorNumber int, players *squad.PlayerSquad) *FloorGenBalancer {
	return &FloorGenBalancer{
		floorNumber: floorNumber,
		players:     players,
	}
}

// FloorGenBalancer interface implementation
func (b *FloorGenBalancer) TreasureRoomsCount() int {
	switch b.floorNumber {
	case 1:
		return 2 + rand.Intn(2)
	case 2:
		return 2 + rand.Intn(3)
	case 3:
		return 3 + rand.Intn(3)
	default:
		log.Panic("FloorGenBalancer floors under 3 not accounted!")
		return 0
	}
}

// FloorGenBalancer interface implementation
func (b *FloorGenBalancer) EnemyRoomsCount() int {
	switch b.floorNumber {
	case 1:
		return 3 + rand.Intn(2)
	case 2:
		return 3 + rand.Intn(2)
	case 3:
		return 4 + rand.Intn(2)
	default:
		log.Panic("FloorGenBalancer floors under 3 not accounted!")
		return 0
	}
}

// FloorGenBalancer interface implementation
func (b *FloorGenBalancer) TrapRoomsCount() int {
	switch b.floorNumber {
	case 1:
		return 2
	case 2:
		return 2 + rand.Intn(2)
	case 3:
		return 2 + rand.Intn(2)
	default:
		log.Panic("FloorGenBalancer floors under 3 not accounted!")
		return 0
	}
}

// FloorGenBalancer interface implementation
func (b *FloorGenBalancer) ItemsInTreasureRoomCount() int {
	return int((1.0 + 0.5*float32(b.floorNumber)) * float32(b.players.LenAlive()))
}

// FloorGenBalancer interface implementation
func (b *FloorGenBalancer) EnemyBalancer() I.EnemyBalancer {
	return NewEnemyBalancer(b.floorNumber, b.players)
}

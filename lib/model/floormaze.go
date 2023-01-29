package model

import (
	"math/rand"
)

type FloorMaze struct {
	name       string
	roomsCount int
	tierMin    int
	tierMax    int
}

func NewFloorMaze(name string, tierMin, tierMax int) *FloorMaze {
	roomsMin := 5 + tierMin
	roomsMax := 5 + tierMax
	roomsCount := roomsMin + rand.Intn(roomsMax-roomsMin+1)

	return &FloorMaze{
		name:       name,
		roomsCount: roomsCount,
		tierMin:    tierMin,
		tierMax:    tierMax,
	}
}

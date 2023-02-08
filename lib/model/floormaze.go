package model

import (
	"log"
	"math/rand"
	"zarg/lib/utils"
)

type FloorMaze struct {
	name     string
	rooms    []string
	currRoom int
}

func NewFloorMaze(name string) *FloorMaze {
	var rooms []string

	enemyRoomsCount := 3
	treasureRoomsCount := 2 + rand.Intn(2)

	for enemyRoomsCount > 0 || treasureRoomsCount > 0 {
		pm := utils.NewPropMap()
		pm.Add("enemy", enemyRoomsCount)
		pm.Add("treasure", treasureRoomsCount)

		switch pm.Choose().(string) {
		case "enemy":
			rooms = append(rooms, "enemy")
			enemyRoomsCount -= 1
		case "treasure":
			if rand.Float32() < 0.5 {
				rooms = append(rooms, "trap", "treasure")
			} else {
				rooms = append(rooms, "enemy", "treasure")
			}
			treasureRoomsCount -= 1
		}
	}

	// add rest room and boss room
	rooms = append(rooms, "rest", "boss")

	return &FloorMaze{
		name:  name,
		rooms: rooms,
	}
}

func (fm *FloorMaze) RoomsCount() int {
	return len(fm.rooms)
}

func (fm *FloorMaze) Room(index int) string {
	return fm.rooms[index]
}

func (fm *FloorMaze) NextRoom() string {
	if fm.currRoom >= len(fm.rooms) {
		log.Panic("NextRoom but there is no rooms left!")
	}
	room := fm.rooms[fm.currRoom]
	fm.currRoom += 1
	return room
}

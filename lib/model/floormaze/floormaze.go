package floormaze

import (
	"log"
	"zarg/lib/model/enemy/boss"
	"zarg/lib/model/enemy/squad"
	I "zarg/lib/model/interfaces"
	"zarg/lib/model/trap"
)

type FloorMaze struct {
	name     string
	rooms    []any
	currRoom int
}

type Room struct {
}

type EnemyRoom struct {
	Room
	Enemies *squad.EnemySquad
}

type TrapRoom struct {
	Room
	Trap *trap.Trap
}

type TreasureRoom struct {
	Room
	Items []I.Pickable
}

type RestRoom struct {
	Room
}

type BossRoom struct {
	Room
	Boss *boss.Boss
}

func newFloorMaze(name string, rooms []any) *FloorMaze {
	return &FloorMaze{
		name:     name,
		rooms:    rooms,
		currRoom: 0,
	}
}

func (fm *FloorMaze) RoomsCount() int {
	return len(fm.rooms)
}

func (fm *FloorMaze) Room(index int) any {
	return fm.rooms[index]
}

func (fm *FloorMaze) NextRoom() any {
	if fm.currRoom >= len(fm.rooms) {
		log.Panic("NextRoom but there is no rooms left!")
	}
	room := fm.rooms[fm.currRoom]
	fm.currRoom += 1
	return room
}

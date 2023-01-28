package model

import (
	"container/list"
	"fmt"
)

const maxPlayerHealth = 100

type Player struct {
	userID int
	name   string
	Health int
	Weapon *Weapon
}

type PlayerList struct {
	list *list.List
}

func NewPlayer(userID int, name string) *Player {
	return &Player{
		userID: userID,
		name:   name,
		Health: maxPlayerHealth,
		Weapon: nil,
	}
}

func (p *Player) UserID() int {
	return p.userID
}

func (p *Player) Name() string {
	return p.name
}

func NewPlayerList() *PlayerList {
	return &PlayerList{
		list: list.New(),
	}
}

func (pl *PlayerList) Empty() bool {
	return pl.list.Len() == 0
}

func (pl *PlayerList) Len() int {
	return pl.list.Len()
}

func (pl *PlayerList) Add(p *Player) {
	pl.list.PushBack(p)
}

func (pl *PlayerList) GetByID(userID int) *Player {
	for node := pl.list.Front(); node != nil; node = node.Next() {
		if node.Value.(*Player).userID == userID {
			return node.Value.(*Player)
		}
	}
	return nil
}

func (pl *PlayerList) RemoveByID(userID int) *Player {
	for node := pl.list.Front(); node != nil; node = node.Next() {
		if node.Value.(*Player).userID == userID {
			return pl.list.Remove(node).(*Player)
		}
	}
	return nil
}

func (pl *PlayerList) Foreach(f func(int, *Player)) {
	for node, i := pl.list.Front(), 0; node != nil; node, i = node.Next(), i+1 {
		f(i, node.Value.(*Player))
	}
}

func (pl *PlayerList) StatsInfo() string {
	inf := "Статы игроков:\n\n"

	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(*Player)
		inf += fmt.Sprintf("%s: %s\nHP: %d\n\n", p.Name(), p.Weapon.SummaryShort(), p.Health)
	}

	return inf
}

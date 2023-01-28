package model

import (
	"container/list"
	"fmt"
	"log"
	"strings"
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

func (pl *PlayerList) SetOrdering(ids []int) {
	if len(ids) != pl.Len() {
		log.Fatalf("len(players list) = %d, but len(ordering) = %d\n", pl.Len(), len(ids))
	}

	newList := list.New()
	for _, id := range ids {
		p := pl.GetByID(id)
		if p == nil {
			log.Fatalf("ordering: %+v is invalid for players list: %+v\n", ids, pl)
		}
		newList.PushBack(p)
	}
	pl.list = newList
}

func (pl *PlayerList) OrderingString() string {
	ordering := make([]string, pl.Len())
	pl.Foreach(func(i int, p *Player) {
		ordering[i] = p.Name()
	})
	return strings.Join(ordering, " -> ")
}

func (pl *PlayerList) PhantomOrderingString(ids []int) string {
	ordering := ""
	for _, id := range ids {
		p := pl.GetByID(id)
		ordering += p.Name() + " -> "
	}
	return ordering + "..."
}

func (pl *PlayerList) ListString() string {
	res := ""
	for i, node := 1, pl.list.Front(); node != nil; i, node = i+1, node.Next() {
		res += fmt.Sprintf("  %d. %s\n", i, node.Value.(*Player).Name())
	}
	return res
}

func (pl *PlayerList) StatsInfo() string {
	inf := "Статы игроков:\n\n"

	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(*Player)
		inf += fmt.Sprintf("%s: %s\nHP: %d\n\n", p.Name(), p.Weapon.SummaryShort(), p.Health)
	}

	return inf
}

package player

import (
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"zarg/lib/model/user"
)

type PlayerSquad struct {
	list *list.List
	iter *list.Element
}

func NewPlayerSquad() *PlayerSquad {
	return &PlayerSquad{
		list: list.New(),
	}
}

func (pl *PlayerSquad) Empty() bool {
	return pl.list.Len() == 0
}

func (pl *PlayerSquad) Len() int {
	return pl.list.Len()
}

func (pl *PlayerSquad) LenAlive() int {
	res := 0

	for i, node := 1, pl.list.Front(); node != nil; i, node = i+1, node.Next() {
		p := node.Value.(*Player)
		if p.Health > 0 {
			res += 1
		}
	}

	return res
}

func (pl *PlayerSquad) Add(p *Player) {
	pl.list.PushBack(p)
}

func (pl *PlayerSquad) GetByID(userID int) *Player {
	for node := pl.list.Front(); node != nil; node = node.Next() {
		if node.Value.(*Player).user.ID == userID {
			return node.Value.(*Player)
		}
	}
	return nil
}

func (pl *PlayerSquad) GetByUser(u *user.User) *Player {
	for node := pl.list.Front(); node != nil; node = node.Next() {
		if node.Value.(*Player).user == u {
			return node.Value.(*Player)
		}
	}
	return nil
}

func (pl *PlayerSquad) RemoveByID(userID int) *Player {
	for node := pl.list.Front(); node != nil; node = node.Next() {
		if node.Value.(*Player).user.ID == userID {
			return pl.list.Remove(node).(*Player)
		}
	}
	return nil
}

func (pl *PlayerSquad) Foreach(f func(int, *Player)) {
	for node, i := pl.list.Front(), 0; node != nil; node, i = node.Next(), i+1 {
		f(i, node.Value.(*Player))
	}
}

func (pl *PlayerSquad) ForeachAlive(f func(int, *Player)) {
	for node, i := pl.list.Front(), 0; node != nil; node = node.Next() {
		p := node.Value.(*Player)
		if p.Alive() {
			f(i, p)
			i += 1
		}
	}
}

func (pl *PlayerSquad) SetOrdering(order []*user.User) {
	if len(order) != pl.Len() {
		log.Panicf("len(players list) = %d, but len(ordering) = %d\n", pl.Len(), len(order))
	}

	newList := list.New()
	for _, u := range order {
		p := pl.GetByUser(u)
		if p == nil {
			log.Panicf("ordering: %+v is invalid for players list: %+v\n", order, pl)
		}
		newList.PushBack(p)
	}
	pl.list = newList
}

func (pl *PlayerSquad) OrderingString() string {
	ordering := make([]string, pl.LenAlive())
	pl.ForeachAlive(func(i int, p *Player) {
		ordering[i] = p.user.FullName()
	})
	return strings.Join(ordering, " -> ")
}

func (pl *PlayerSquad) ListString() string {
	res := ""
	for i, node := 1, pl.list.Front(); node != nil; i, node = i+1, node.Next() {
		res += fmt.Sprintf("  %d. %s\n", i, node.Value.(*Player).user.FullName())
	}
	return res
}

func (pl *PlayerSquad) Info() string {
	inf := ""
	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(*Player)
		if p.Health == 0 {
			inf += fmt.Sprintf("%s: мертв\n", p.user.FullName())
		} else {
			inf += fmt.Sprintf("%s: HP=%d, %s\n", p.user.FullName(), p.Health, p.Weapon.SummaryShort())
		}
	}
	return inf
}

func (pl *PlayerSquad) CompactInfo() string {
	inf := ""
	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(*Player)
		if p.Health == 0 {
			inf += fmt.Sprintf("%s: мертв\n", p.user.FullName())
		} else {
			inf += fmt.Sprintf("%s: HP=%d, %s\n", p.user.FullName(), p.Health, p.Weapon.SummaryShort())
		}
	}
	return inf
}

func (pl *PlayerSquad) EndFight() {
	pl.iter = nil
}

func (pl *PlayerSquad) ChooseNext() *Player {
	if pl.LenAlive() == 0 {
		log.Fatal("all players dead!")
	}

	if pl.iter == nil {
		pl.iter = pl.list.Front()
	} else {
		pl.iter = pl.iter.Next()
		if pl.iter == nil {
			pl.iter = pl.list.Front()
		}
	}

	for pl.iter.Value.(*Player).Health == 0 {
		pl.iter = pl.iter.Next()
		if pl.iter == nil {
			pl.iter = pl.list.Front()
		}
	}

	return pl.iter.Value.(*Player)
}

func (pl *PlayerSquad) ChooseFirstAlive() *Player {
	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(*Player)
		if p.Alive() {
			return p
		}
	}

	log.Panic("everyone dead, nothing to choose from!")
	return nil
}

func (pl *PlayerSquad) ChooseRandomAlive() *Player {
	n := rand.Intn(pl.LenAlive())

	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(*Player)
		if p.Alive() {
			if n == 0 {
				return p
			}
			n -= 1
		}
	}

	log.Panic("everyone dead, nothing to choose from!")
	return nil
}

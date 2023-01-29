package model

import (
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"strings"
)

const maxPlayerHealth = 100

type Player struct {
	user   *User
	Health int
	Weapon *Weapon
}

type PlayerList struct {
	list *list.List
	iter *list.Element
}

func NewPlayer(user *User) *Player {
	return &Player{
		user:   user,
		Health: maxPlayerHealth,
		Weapon: nil,
	}
}

func (p *Player) User() *User {
	return p.user
}

func (p *Player) MakeDamage(val int) {
	p.Health -= val
	if p.Health < 0 {
		p.Health = 0
	}
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

func (pl *PlayerList) LenAlive() int {
	res := 0

	for i, node := 1, pl.list.Front(); node != nil; i, node = i+1, node.Next() {
		p := node.Value.(*Player)
		if p.Health > 0 {
			res += 1
		}
	}

	return res
}

func (pl *PlayerList) Add(p *Player) {
	pl.list.PushBack(p)
}

func (pl *PlayerList) GetByID(userID int) *Player {
	for node := pl.list.Front(); node != nil; node = node.Next() {
		if node.Value.(*Player).user.ID == userID {
			return node.Value.(*Player)
		}
	}
	return nil
}

func (pl *PlayerList) GetByUser(u *User) *Player {
	for node := pl.list.Front(); node != nil; node = node.Next() {
		if node.Value.(*Player).user == u {
			return node.Value.(*Player)
		}
	}
	return nil
}

func (pl *PlayerList) RemoveByID(userID int) *Player {
	for node := pl.list.Front(); node != nil; node = node.Next() {
		if node.Value.(*Player).user.ID == userID {
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

func (pl *PlayerList) SetOrdering(users []*User) {
	if len(users) != pl.Len() {
		log.Fatalf("len(players list) = %d, but len(ordering) = %d\n", pl.Len(), len(users))
	}

	newList := list.New()
	for _, u := range users {
		p := pl.GetByUser(u)
		if p == nil {
			log.Fatalf("ordering: %+v is invalid for players list: %+v\n", users, pl)
		}
		newList.PushBack(p)
	}
	pl.list = newList
}

func (pl *PlayerList) OrderingString() string {
	ordering := make([]string, pl.Len())
	pl.Foreach(func(i int, p *Player) {
		ordering[i] = p.user.FullName()
	})
	return strings.Join(ordering, " -> ")
}

func (pl *PlayerList) PhantomOrderingString(users []*User) string {
	ordering := ""
	for _, u := range users {
		ordering += u.FullName() + " -> "
	}
	return ordering + "..."
}

func (pl *PlayerList) ListString() string {
	res := ""
	for i, node := 1, pl.list.Front(); node != nil; i, node = i+1, node.Next() {
		res += fmt.Sprintf("  %d. %s\n", i, node.Value.(*Player).user.FullName())
	}
	return res
}

func (pl *PlayerList) Info() string {
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

func (pl *PlayerList) CompactInfo() string {
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

func (pl *PlayerList) EndFight() {
	pl.iter = nil
}

func (pl *PlayerList) ChooseNext() *Player {
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

func (pl *PlayerList) ChooseRandomAlive() *Player {
	n := rand.Intn(pl.LenAlive())

	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(*Player)
		if p.Health > 0 {
			if n == 0 {
				return p
			}
			n -= 1
		}
	}

	log.Panic("everyone dead, nothing to choose from!")
	return nil
}

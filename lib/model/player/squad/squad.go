package squad

import (
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"strings"
	I "zarg/lib/model/interfaces"
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

// EntityList interface implementation
func (pl PlayerSquad) Len() int {
	return pl.list.Len()
}

// EntityList interface implementation
func (pl PlayerSquad) LenAlive() int {
	res := 0

	for i, node := 1, pl.list.Front(); node != nil; i, node = i+1, node.Next() {
		p := node.Value.(I.Player)
		if p.Alive() {
			res += 1
		}
	}

	return res
}

// EntityList interface implementation
func (pl *PlayerSquad) ForEach(f func(I.Entity)) {
	for node := pl.list.Front(); node != nil; node = node.Next() {
		f(node.Value.(I.Player))
	}
}

// PlayerList interface implementation
func (pl *PlayerSquad) ForEachAlive(f func(I.Entity)) {
	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(I.Player)
		if p.Alive() {
			f(p)
		}
	}
}

func (es *PlayerSquad) Has(player I.Entity) bool {
	for node := es.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(I.Entity)
		if p == player {
			return true
		}
	}
	return false
}

func (pl *PlayerSquad) Add(p I.Player) {
	pl.list.PushBack(p)
}

func (pl *PlayerSquad) GetByID(userID int) I.Player {
	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(I.Player)
		if p.ID() == userID {
			return p
		}
	}
	return nil
}

func (pl *PlayerSquad) RemoveByID(userID int) I.Player {
	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(I.Player)
		if p.ID() == userID {
			return pl.list.Remove(node).(I.Player)
		}
	}
	return nil
}

func (pl *PlayerSquad) SetOrdering(order []int) {
	if len(order) != pl.Len() {
		log.Panicf("len(players list) = %d, but len(ordering) = %d\n", pl.Len(), len(order))
	}

	newList := list.New()
	for _, id := range order {
		p := pl.GetByID(id)
		if p == nil {
			log.Panicf("ordering: %+v is invalid for players list: %+v\n", order, pl)
		}
		newList.PushBack(p)
	}
	pl.list = newList
}

func (pl *PlayerSquad) OrderingString() string {
	ordering := make([]string, 0, pl.LenAlive())
	pl.ForEachAlive(func(p I.Entity) {
		ordering = append(ordering, p.(I.Player).FullName())
	})
	return strings.Join(ordering, " -> ")
}

func (pl *PlayerSquad) ListString() string {
	res := ""
	for i, node := 1, pl.list.Front(); node != nil; i, node = i+1, node.Next() {
		res += fmt.Sprintf("  %d. %s\n", i, node.Value.(I.Player).FullName())
	}
	return res
}

func (pl *PlayerSquad) Info() string {
	inf := ""
	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(I.Player)
		if p.Alive() {
			effects := pl.EffectsCompactInfo(p)
			inf += fmt.Sprintf("%s (%d❤ %d🗡) %s\n", p.FullName(), p.Health(), p.AttackStats().TypedDamages()[I.DamageType1], effects)
			inf += fmt.Sprintf("оружие: %s. %s\n", p.Weapon().Name(), p.Weapon().Description())
			var items []string
			p.ForEachItem(func(p I.Pickable) {
				switch p := p.(type) {
				case I.Consumable:
					items = append(items, fmt.Sprintf("%s [x%d] (%s)", p.Name(), p.UsesLeft(), p.Description()))
				default:
					items = append(items, fmt.Sprintf("%s (%s)", p.Name(), p.Description()))
				}
			})
			if len(items) > 0 {
				inf += fmt.Sprintf("%s.\n\n", strings.Join(items, ", "))
			} else {
				inf += "\n"
			}
		} else {
			inf += fmt.Sprintf("%s: 💀\n\n", p.FullName())
		}
	}
	return inf
}

func (pl *PlayerSquad) CompactInfo() string {
	inf := ""
	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(I.Player)
		if p.Alive() {
			effects := pl.EffectsCompactInfo(p)
			inf += fmt.Sprintf("- %s (%d❤ %d🗡 %s)\n", p.FullName(), p.Health(), p.AttackStats().TypedDamages()[I.DamageType1], effects)
		} else {
			inf += fmt.Sprintf("- %s 💀\n", p.FullName())
		}
	}
	return inf
}

func (pl *PlayerSquad) EffectsCompactInfo(p I.Player) string {
	var effects []string
	for _, eff := range p.StatusEffects() {
		effects = append(effects, fmt.Sprintf("%sx%d", eff.Name, eff.TimeLeft))
	}
	if len(effects) == 0 {
		return ""
	} else {
		return "[" + strings.Join(effects, ",") + "]"
	}
}

func (pl *PlayerSquad) EndFight() {
	pl.iter = nil
}

func (pl *PlayerSquad) ChooseNext() I.Player {
	if pl.LenAlive() == 0 {
		log.Panic("all players dead!")
	}

	if pl.iter == nil {
		pl.iter = pl.list.Front()
	} else {
		pl.iter = pl.iter.Next()
		if pl.iter == nil {
			pl.iter = pl.list.Front()
		}
	}

	for !pl.iter.Value.(I.Player).Alive() {
		pl.iter = pl.iter.Next()
		if pl.iter == nil {
			pl.iter = pl.list.Front()
		}
	}

	return pl.iter.Value.(I.Player)
}

func (pl *PlayerSquad) ChooseFirstAlive() I.Player {
	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(I.Player)
		if p.Alive() {
			return p
		}
	}

	log.Panic("everyone dead, nothing to choose from!")
	return nil
}

func (pl *PlayerSquad) ChooseRandomAlive() I.Player {
	n := rand.Intn(pl.LenAlive())

	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(I.Player)
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

func (pl *PlayerSquad) ChooseRandomAlivePreferBlocking() I.Player {
	blocks := pl.countBlocks()
	if blocks == 0 {
		return pl.ChooseRandomAlive()
	}

	n := rand.Intn(blocks)
	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(I.Player)
		if p.Alive() && p.IsBlocking() {
			if n == 0 {
				return p
			}
			n -= 1
		}
	}

	log.Panic("must never happen!")
	return nil
}

func (pl *PlayerSquad) countBlocks() int {
	blocks := 0
	for node := pl.list.Front(); node != nil; node = node.Next() {
		p := node.Value.(I.Player)
		if p.Alive() && p.IsBlocking() {
			blocks += 1
		}
	}
	return blocks
}

package trap

import (
	"log"
	"zarg/lib/model/player"
)

type DamageType int

const (
	DamageEveryone = DamageType(iota)
	DamageRandom   = DamageType(iota)
	DamageFirst    = DamageType(iota)
)

type Trap struct {
	description string
	damageType  DamageType
	attack      int
}

func New(desc string, damageType DamageType, attack int) *Trap {
	return &Trap{
		description: desc,
		damageType:  damageType,
		attack:      attack,
	}
}

func (t *Trap) Description() string {
	return t.description
}

func (t *Trap) DamagesEveryone() bool {
	return t.damageType == DamageEveryone
}

// returns players damaged by this trap
func (t *Trap) Activate(pl *player.PlayerSquad) []*player.Player {
	switch t.damageType {
	case DamageFirst:
		return t.damageFirst(pl)
	case DamageRandom:
		return t.damageRandom(pl)
	case DamageEveryone:
		return t.damageEveryone(pl)
	default:
		log.Panicf("invalid trap damage type: %+v", t)
		return nil
	}
}

func (t *Trap) damageFirst(pl *player.PlayerSquad) []*player.Player {
	p := pl.ChooseFirstAlive()
	p.MakeDamage(t.attack)
	return append(make([]*player.Player, 0, 1), p)
}

func (t *Trap) damageRandom(pl *player.PlayerSquad) []*player.Player {
	p := pl.ChooseRandomAlive()
	p.MakeDamage(t.attack)
	return append(make([]*player.Player, 0, 1), p)
}

func (t *Trap) damageEveryone(pl *player.PlayerSquad) []*player.Player {
	var res []*player.Player

	pl.ForeachAlive(func(_ int, p *player.Player) {
		p.MakeDamage(t.attack)
		res = append(res, p)
	})

	return res
}

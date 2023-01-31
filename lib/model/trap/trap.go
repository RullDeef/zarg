package trap

import (
	"log"
	I "zarg/lib/model/interfaces"
	"zarg/lib/model/player/squad"
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
	attack      I.DamageStats
}

func New(desc string, damageType DamageType, attack int) *Trap {
	t := &Trap{
		description: desc,
		damageType:  damageType,
		attack: I.DamageStats{
			Base:       attack,
			Crit:       attack,
			CritChance: 0.0,
		},
	}

	t.attack.Producer = t
	return t
}

// DamageProducer interface implementation
func (t *Trap) Name() string {
	return t.description
}

func (t *Trap) DamagesEveryone() bool {
	return t.damageType == DamageEveryone
}

// returns players damaged by this trap
func (t *Trap) Activate(pl *squad.PlayerSquad) []I.Player {
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

func (t *Trap) damageFirst(pl *squad.PlayerSquad) []I.Player {
	p := pl.ChooseFirstAlive()
	p.Damage(t.attack)
	return append(make([]I.Player, 0, 1), p)
}

func (t *Trap) damageRandom(pl *squad.PlayerSquad) []I.Player {
	p := pl.ChooseRandomAlive()
	p.Damage(t.attack)
	return append(make([]I.Player, 0, 1), p)
}

func (t *Trap) damageEveryone(pl *squad.PlayerSquad) []I.Player {
	var res []I.Player

	pl.ForEachAlive(func(p I.Player) {
		p.Damage(t.attack)
		res = append(res, p)
	})

	return res
}

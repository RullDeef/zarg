package trap

import (
	"log"
	I "zarg/lib/model/interfaces"
	"zarg/lib/model/player/squad"
)

const (
	DamageEveryone = I.DamageType(iota)
	DamageRandom   = I.DamageType(iota)
	DamageFirst    = I.DamageType(iota)
)

type Trap struct {
	description string
	damageType  I.DamageType
	attack      I.DamageStats
}

type TrapAttack struct {
	typedDamages map[I.DamageType]int

	critChance float64
	critFactor float64
	isCrit     bool
}

func New(desc string, damageType I.DamageType, attack int) *Trap {
	t := &Trap{
		description: desc,
		damageType:  damageType,
		attack: &TrapAttack{
			typedDamages: map[I.DamageType]int{
				damageType: attack,
			},
			critChance: 0.5,
			critFactor: 1.2,
		},
	}
	return t
}

// DamageEmitor interface implementation
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
	p.Damage(t.attack.(*TrapAttack))
	return append(make([]I.Player, 0, 1), p)
}

func (t *Trap) damageRandom(pl *squad.PlayerSquad) []I.Player {
	p := pl.ChooseRandomAlive()
	p.Damage(t.attack.(*TrapAttack))
	return append(make([]I.Player, 0, 1), p)
}

func (t *Trap) damageEveryone(pl *squad.PlayerSquad) []I.Player {
	var res []I.Player

	pl.ForEachAlive(func(p I.Entity) {
		p.Damage(t.attack.(*TrapAttack))
		res = append(res, p.(I.Player))
	})

	return res
}

func (t *TrapAttack) TypedDamages() map[I.DamageType]int {
	return t.typedDamages
}

func (t *TrapAttack) CritChance() float64 {
	return t.critChance
}

func (t *TrapAttack) CritFactor() float64 {
	return t.critFactor
}

func (t *TrapAttack) IsCrit() bool {
	return t.isCrit
}

func (t *TrapAttack) StatusEffectChances() map[I.StatusEffect]float64 {
	return make(map[I.StatusEffect]float64)
}

func (t *TrapAttack) StatusEffects() []I.StatusEffect {
	return nil
}

package armor

import (
	"fmt"
	"math/rand"
	I "zarg/lib/model/interfaces"
)

type ArmorItem struct {
	name    string
	defence int
	owner   I.Entity

	// chance to produce a status effect on its owner
	statusEffectChances map[I.StatusEffect]float64
}

var armorNames = []string{
	"Стальные латы",
	"Кольчуга",
	"Маскировочный костюм",
}

func New(
	name string,
	defence int,
	statusEffectChances map[I.StatusEffect]float64,
) *ArmorItem {
	return &ArmorItem{
		name:                name,
		defence:             defence,
		owner:               nil,
		statusEffectChances: statusEffectChances,
	}
}

func Random() *ArmorItem {
	name := armorNames[rand.Intn(len(armorNames))]

	defenceMin := 1
	defenceMax := 10
	defence := defenceMin + rand.Intn(defenceMax-defenceMin+1)

	statusEffectChances := make(map[I.StatusEffect]float64)
	if rand.Float64() < 0.5 {
		time := rand.Intn(2) + 2
		chance := float64(rand.Intn(3)+1) / 20
		statusEffectChances[I.StatusEffectAgility(time)] = chance
	}
	if rand.Float64() < 0.2 {
		time := rand.Intn(2) + 2
		chance := float64(rand.Intn(3)+1) / 20
		statusEffectChances[I.StatusEffectRegeneration(time)] = chance
	}

	return New(name, defence, statusEffectChances)
}

// Pickable interface implementation
func (a ArmorItem) Name() string {
	return a.name
}

func (a ArmorItem) Description() string {
	msg := fmt.Sprintf("%d🛡", a.defence)

	for effect, chance := range a.statusEffectChances {
		msg += fmt.Sprintf("[%dx%s:%.0f%%]", effect.TimeLeft, effect.Name, chance*100)
	}

	return msg
}

// Pickable interface implementation
func (a ArmorItem) Owner() I.Entity {
	return a.owner
}

// Pickable interface implementation
func (a *ArmorItem) SetOwner(p I.Entity) {
	a.owner = p
}

// Pickable interface implementation
func (a *ArmorItem) ModifyOngoingDamage(dmg I.Damage) I.Damage {
	dmg.TypedDamages()[I.DamageType1] -= a.defence
	if dmg.TypedDamages()[I.DamageType1] < 0 {
		dmg.TypedDamages()[I.DamageType1] = 0
	}

	// try luck with status effects
	for effect, chance := range a.statusEffectChances {
		if rand.Float64() < chance {
			a.owner.AddStatusEffect(effect)
		}
	}

	return dmg
}

// Pickable interface implementation
func (a *ArmorItem) ModifyOutgoingDamage(dmg I.Damage) I.Damage {
	return dmg
}

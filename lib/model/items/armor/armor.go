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
}

var armorNames = []string{
	"Стальные латы",
	"Кольчуга",
	"Маскировочный костюм",
}

func New(name string, defence int) *ArmorItem {
	return &ArmorItem{
		name:    name,
		defence: defence,
		owner:   nil,
	}
}

func Random() *ArmorItem {
	name := armorNames[rand.Intn(len(armorNames))]

	defenceMin := 1
	defenceMax := 10
	defence := defenceMin + rand.Intn(defenceMax-defenceMin+1)

	return New(name, defence)
}

// Pickable interface implementation
func (a ArmorItem) Name() string {
	return a.name
}

func (a ArmorItem) Description() string {
	return fmt.Sprintf("%d🛡", a.defence)
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
	return dmg
}

// Pickable interface implementation
func (a *ArmorItem) ModifyOutgoingDamage(dmg I.Damage) I.Damage {
	return dmg
}

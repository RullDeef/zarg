package armor

import (
	"fmt"
	"math/rand"
	I "zarg/lib/model/interfaces"
)

type ArmorItem struct {
	name    string
	defence int
	owner   I.Player
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
func (a ArmorItem) Owner() I.Player {
	return a.owner
}

// Pickable interface implementation
func (a *ArmorItem) SetOwner(p I.Player) {
	a.owner = p
}

// Pickable interface implementation
func (a *ArmorItem) ModifyOngoingDamage(dmg I.DamageStats) I.DamageStats {
	dmg.Base -= a.defence
	if dmg.Base < 0 {
		dmg.Base = 0
	}
	dmg.Crit -= a.defence
	if dmg.Crit < 0 {
		dmg.Crit = 0
	}
	return dmg
}

// Pickable interface implementation
func (a *ArmorItem) ModifyOutgoingDamage(dmg I.DamageStats) I.DamageStats {
	return dmg
}

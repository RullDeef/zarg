package potion

import (
	"fmt"
	I "zarg/lib/model/interfaces"
	"zarg/lib/utils"
)

type HealingPotion struct {
	name     string
	value    int
	owner    I.Player
	usesLeft int
}

func Random() I.Pickable {
	pm := utils.NewPropMap()

	pm.Add(NewHealingPotion("Зелье восстановления I", 20, 1), 3)
	pm.Add(NewHealingPotion("Зелье восстановления I", 20, 3), 2)
	pm.Add(NewHealingPotion("Зелье восстановления I", 20, 5), 1)
	pm.Add(NewHealingPotion("Зелье восстановления II", 35, 1), 2)
	pm.Add(NewHealingPotion("Зелье восстановления II", 35, 2), 1)

	return pm.Choose().(I.Pickable)
}

func NewHealingPotion(name string, value int, amount int) *HealingPotion {
	return &HealingPotion{
		name:     name,
		value:    value,
		owner:    nil,
		usesLeft: amount,
	}
}

// Pickable interface implementation
func (hp *HealingPotion) Name() string {
	return hp.name
}

// Pickable interface implementation
func (hp *HealingPotion) Owner() I.Player {
	return hp.owner
}

// Pickable interface implementation
func (hp *HealingPotion) SetOwner(p I.Player) {
	hp.owner = p
}

// Pickable interface implementation
func (hp *HealingPotion) ModifyOngoingDamage(ds I.DamageStats) I.DamageStats {
	return ds
}

// Pickable interface implementation
func (hp *HealingPotion) ModifyOutgoingDamage(ds I.DamageStats) I.DamageStats {
	return ds
}

// Consumable interface implementation
func (hp *HealingPotion) Description() string {
	return fmt.Sprintf("+%dHP", hp.value)
}

// Consumable interface implementation
func (hp *HealingPotion) UsesLeft() int {
	return hp.usesLeft
}

// Consumable interface implementation
func (hp *HealingPotion) Consume() {
	if hp.usesLeft > 0 {
		hp.owner.Heal(hp.value)
		hp.usesLeft -= 1
	}
}

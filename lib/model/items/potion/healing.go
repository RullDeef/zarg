package potion

import (
	"fmt"
	I "zarg/lib/model/interfaces"
)

type HealingPotion struct {
	name     string
	value    int
	owner    I.Entity
	usesLeft int
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
func (hp *HealingPotion) Owner() I.Entity {
	return hp.owner
}

// Pickable interface implementation
func (hp *HealingPotion) SetOwner(p I.Entity) {
	hp.owner = p
}

// Pickable interface implementation
func (hp *HealingPotion) ModifyOngoingDamage(dmg I.Damage) I.Damage {
	return dmg
}

// Pickable interface implementation
func (hp *HealingPotion) ModifyOutgoingDamage(dmg I.Damage) I.Damage {
	return dmg
}

func (hp *HealingPotion) Stack(item I.Pickable) bool {
	switch item := item.(type) {
	case *HealingPotion:
		if hp.name == item.name {
			hp.usesLeft += item.usesLeft
			return true
		}
	default:
	}
	return false
}

// Consumable interface implementation
func (hp *HealingPotion) Description() string {
	return fmt.Sprintf("+%d❤", hp.value)
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

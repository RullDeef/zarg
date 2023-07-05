package potion

import (
	"fmt"
	I "zarg/lib/model/interfaces"
)

type StrengthPotion struct {
	name     string
	owner    I.Entity
	turns    int
	usesLeft int
}

func NewStrengthPotion(name string, turns int, amount int) *StrengthPotion {
	return &StrengthPotion{
		name:     name,
		owner:    nil,
		turns:    turns,
		usesLeft: amount,
	}
}

// Pickable interface implementation
func (hp *StrengthPotion) Name() string {
	return hp.name
}

// Pickable interface implementation
func (hp *StrengthPotion) Owner() I.Entity {
	return hp.owner
}

// Pickable interface implementation
func (hp *StrengthPotion) SetOwner(p I.Entity) {
	hp.owner = p
}

// Pickable interface implementation
func (hp *StrengthPotion) ModifyOngoingDamage(dmg I.Damage) I.Damage {
	return dmg
}

// Pickable interface implementation
func (hp *StrengthPotion) ModifyOutgoingDamage(dmg I.Damage) I.Damage {
	return dmg
}

func (hp *StrengthPotion) Stack(item I.Pickable) bool {
	switch item := item.(type) {
	case *StrengthPotion:
		if hp.name == item.name && hp.turns == item.turns {
			hp.usesLeft += item.usesLeft
			return true
		}
	default:
	}
	return false
}

// Consumable interface implementation
func (hp *StrengthPotion) Description() string {
	return fmt.Sprintf("x1.25🗡(%d)", hp.turns)
}

// Consumable interface implementation
func (hp *StrengthPotion) UsesLeft() int {
	return hp.usesLeft
}

// Consumable interface implementation
func (hp *StrengthPotion) Consume() {
	if hp.usesLeft > 0 {
		hp.owner.AddStatusEffect(I.StatusEffectStrongness(hp.turns))
		hp.usesLeft -= 1
	}
}

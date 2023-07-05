package potion

import (
	"testing"
	I "zarg/lib/model/interfaces"
)

func TestInterfaces(t *testing.T) {
	var _ I.Pickable = &HealingPotion{}
	var _ I.Consumable = &HealingPotion{}
}

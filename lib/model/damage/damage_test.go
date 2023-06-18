package damage

import (
	"testing"
	I "zarg/lib/model/interfaces"
)

func TestInterfaces(t *testing.T) {
	var _ I.DamageStats = &BaseDamageStats{}
	var _ I.Damage = &BaseDamage{}
}

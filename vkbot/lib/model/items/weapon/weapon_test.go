package weapon

import (
	"testing"
	I "zarg/lib/model/interfaces"
)

func TestInterfaces(t *testing.T) {
	var _ I.Weapon = &Weapon{}
}

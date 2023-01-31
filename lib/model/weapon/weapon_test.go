package weapon

import (
	I "zarg/lib/model/interfaces"
)

func TestInterfaces() {
	var _ I.Weapon = &Weapon{}
}

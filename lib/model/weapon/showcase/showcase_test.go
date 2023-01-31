package showcase

import (
	I "zarg/lib/model/interfaces"
)

func TestInterfaces() {
	var _ I.WeaponShowcase = &WeaponShowcase{}
}

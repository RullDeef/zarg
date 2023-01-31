package enemy

import (
	I "zarg/lib/model/interfaces"
)

func TestInerfaces() {
	var _ I.Enemy = &Enemy{}
}

package enemy

import (
	"testing"
	I "zarg/lib/model/interfaces"
)

func TestInerfaces(t *testing.T) {
	var _ I.Enemy = &Enemy{}
}

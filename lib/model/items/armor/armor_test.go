package armor

import (
	"testing"
	"zarg/lib/model/interfaces"
)

func TestInterfaces(t *testing.T) {
	var _ interfaces.Pickable = &ArmorItem{}
}

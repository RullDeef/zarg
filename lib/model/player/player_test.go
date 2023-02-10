package player

import (
	"testing"
	i "zarg/lib/model/interfaces"
)

// check interface implementation
func TestInterfaces(t *testing.T) {
	var _ i.Player = &Player{}
}

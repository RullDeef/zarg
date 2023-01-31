package player

import (
	i "zarg/lib/model/interfaces"
)

// check interface implementation
func TestInterfaces() {
	var _ i.Player = &Player{}
}

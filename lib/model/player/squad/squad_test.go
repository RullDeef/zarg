package squad

import (
	I "zarg/lib/model/interfaces"
)

func TestInterfaces() {
	var _ I.PlayerList = &PlayerSquad{}
}

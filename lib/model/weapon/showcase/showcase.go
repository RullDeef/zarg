package showcase

import (
	"fmt"
	"sync"
	I "zarg/lib/model/interfaces"
)

// starting weapons that players may choose from
type WeaponShowcase struct {
	weapons  []I.Weapon
	pickedBy []I.Player
	lock     sync.Mutex
}

func NewWeaponShowcase(n int, maker func() I.Weapon) *WeaponShowcase {
	weapons := make([]I.Weapon, n)
	pickedBy := make([]I.Player, n)

	for i := 0; i < n; i++ {
		weapons[i] = maker()
	}

	return &WeaponShowcase{
		weapons:  weapons,
		pickedBy: pickedBy,
		lock:     sync.Mutex{},
	}
}

// WeaponShowcase interface implementation
func (ws *WeaponShowcase) WeaponsInfo() string {
	res := ""
	for i, weapon := range ws.weapons {
		res += fmt.Sprintf("  %d) %s. %s\n", i+1, weapon.Name(), weapon.Description())
	}
	return res
}

// WeaponShowcase interface implementation
func (ws *WeaponShowcase) HasMadePick(player I.Player) bool {
	for _, p := range ws.pickedBy {
		if p == player {
			return true
		}
	}
	return false
}

// WeaponShowcase interface implementation
func (ws *WeaponShowcase) TryPick(p I.Player, i int) (bool, I.Weapon, I.Player) {
	ws.lock.Lock()
	defer ws.lock.Unlock()

	w := ws.weapons[i]
	if ws.pickedBy[i] == nil {
		ws.unpick(p)
		ws.pickedBy[i] = p
		return true, w, p
	}

	if ws.pickedBy[i] == p {
		return true, w, p
	}

	return false, w, ws.pickedBy[i]
}

// WeaponShowcase interface implementation
func (ws *WeaponShowcase) ConfirmPick() {
	ws.lock.Lock()
	defer ws.lock.Unlock()

	for i, p := range ws.pickedBy {
		if p != nil {
			p.PickWeapon(ws.weapons[i])
		}
	}
}

func (ws *WeaponShowcase) unpick(player I.Player) {
	for i, p := range ws.pickedBy {
		if p == player {
			ws.pickedBy[i] = nil
		}
	}
}

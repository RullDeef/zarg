package weaponshowcase

import (
	"fmt"
	"sync"
	"zarg/lib/model/player"
	"zarg/lib/model/weapon"
)

// starting weapons that players may choose from
type WeaponShowcase struct {
	weapons  []*weapon.Weapon
	pickedBy []*player.Player
	lock     sync.Mutex
}

func NewWeaponShowcase(n int, maker func() *weapon.Weapon) *WeaponShowcase {
	weapons := make([]*weapon.Weapon, n)
	pickedBy := make([]*player.Player, n)

	for i := 0; i < n; i++ {
		weapons[i] = maker()
	}

	return &WeaponShowcase{
		weapons:  weapons,
		pickedBy: pickedBy,
		lock:     sync.Mutex{},
	}
}

func (ws *WeaponShowcase) Weapons() []*weapon.Weapon {
	return ws.weapons
}

func (ws *WeaponShowcase) WeaponsInfo() string {
	res := ""
	for i, weapon := range ws.weapons {
		res += fmt.Sprintf("  %d) %s\n", i+1, weapon.Summary())
	}
	return res
}

func (ws *WeaponShowcase) ConfirmPick() {
	ws.lock.Lock()
	defer ws.lock.Unlock()

	for i, p := range ws.pickedBy {
		if p != nil {
			p.Weapon = ws.weapons[i]
		}
	}
}

func (ws *WeaponShowcase) HasMadePick(player *player.Player) bool {
	for _, p := range ws.pickedBy {
		if p == player {
			return true
		}
	}
	return false
}

func (ws *WeaponShowcase) TryPick(p *player.Player, i int) (bool, *weapon.Weapon, *player.Player) {
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

func (ws *WeaponShowcase) unpick(player *player.Player) {
	for i, p := range ws.pickedBy {
		if p == player {
			ws.pickedBy[i] = nil
		}
	}
}

package model

import (
	"fmt"
	"sync"
)

// starting weapons that players may choose from
type WeaponShowcase struct {
	weapons  []*Weapon
	pickedBy []*Player
	lock     sync.Mutex
}

func NewWeaponShowcase(n int, maker func() *Weapon) *WeaponShowcase {
	weapons := make([]*Weapon, n)
	pickedBy := make([]*Player, n)

	for i := 0; i < n; i++ {
		weapons[i] = maker()
	}

	return &WeaponShowcase{
		weapons:  weapons,
		pickedBy: pickedBy,
		lock:     sync.Mutex{},
	}
}

func (ws *WeaponShowcase) Weapons() []*Weapon {
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

func (ws *WeaponShowcase) HasMadePick(player *Player) bool {
	for _, p := range ws.pickedBy {
		if p == player {
			return true
		}
	}
	return false
}

func (ws *WeaponShowcase) TryPick(p *Player, i int) (bool, *Weapon, *Player) {
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

func (ws *WeaponShowcase) unpick(player *Player) {
	for i, p := range ws.pickedBy {
		if p == player {
			ws.pickedBy[i] = nil
		}
	}
}

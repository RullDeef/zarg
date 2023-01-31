package player

import (
	"log"
	"math/rand"
	I "zarg/lib/model/interfaces"
)

const maxHealth = 100

type Player struct {
	user   I.User
	health int
	weapon I.Weapon
	items  []I.Pickable
}

func NewPlayer(user I.User) *Player {
	return &Player{
		user:   user,
		health: maxHealth,
		weapon: nil,
		items:  nil,
	}
}

// User interface implementation
func (p Player) ID() int {
	return p.user.ID()
}

// User interface implementation
func (p Player) FirstName() string {
	return p.user.FirstName()
}

// User interface implementation
func (p Player) LastName() string {
	return p.user.LastName()
}

// User interface implementation
func (p Player) FullName() string {
	return p.user.FullName()
}

// Entity interface implementation
func (p Player) Name() string {
	return p.user.FullName()
}

// Entity interface implementation
func (p Player) Health() int {
	return p.health
}

// Entity interface implementation
func (p *Player) Heal(value int) {
	p.health += value
	if p.health > maxHealth {
		p.health = maxHealth
	}
}

// Entity interface implementation
func (p *Player) Damage(dmg I.DamageStats) int {
	for _, it := range p.items {
		dmg = it.ModifyOngoingDamage(dmg)
	}

	val := dmg.Base
	if rand.Float32() < dmg.CritChance {
		val = dmg.Crit
	}

	p.health -= val
	if p.health < 0 {
		p.health = 0
	}
	return val
}

// Entity interface implementation
func (p Player) Alive() bool {
	return p.health > 0
}

// Player interface implementation
func (p Player) Weapon() I.Weapon {
	return p.weapon
}

// Player interface implementation
func (p *Player) PickWeapon(w I.Weapon) {
	p.weapon = w
	w.SetOwner(p)
}

// Player interface implementation
func (p *Player) Attack() I.DamageStats {
	dmg := p.weapon.Attack()
	for _, it := range p.items {
		dmg = it.ModifyOutgoingDamage(dmg)
	}
	return dmg
}

// Player interface implementation
func (p *Player) PickItem(item I.Pickable) {
	p.items = append(p.items, item)
}

// Player interface implementation
func (p *Player) DropItem(item I.Pickable) {
	for i, it := range p.items {
		if it == item {
			p.items = append(p.items[:i], p.items[i+1:]...)
			return
		}
	}
	log.Panicf("no such item (%+v) with player (%+v)", item, p)
}

// Player interface implementation
func (p *Player) ForEachItem(f func(I.Pickable)) {
	for _, item := range p.items {
		f(item)
	}
}

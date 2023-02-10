package player

import (
	"log"
	"math/rand"
	I "zarg/lib/model/interfaces"
)

const maxHealth = 100

type Player struct {
	user       I.User
	health     int
	weapon     I.Weapon
	items      []I.Pickable
	isBlocking bool
}

func NewPlayer(user I.User) *Player {
	return &Player{
		user:       user,
		health:     maxHealth,
		weapon:     nil,
		items:      nil,
		isBlocking: false,
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

	// if is blocking - reduce to 80%
	if p.isBlocking {
		val = int(0.8 * float32(val))
		p.isBlocking = false
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
	if p.weapon != nil {
		p.weapon.SetOwner(nil)
	}
	p.weapon = w
	w.SetOwner(p)
}

// Player interface implementation
func (p *Player) Attack() I.DamageStats {
	dmg := p.weapon.Attack()
	for _, it := range p.items {
		dmg = it.ModifyOutgoingDamage(dmg)
	}
	p.isBlocking = false
	return dmg
}

// Player interface implementation
func (p *Player) BlockAttack() {
	p.isBlocking = true
}

// Player interface implementation
func (p *Player) IsBlocking() bool {
	return p.isBlocking
}

// Player interface implementation
func (p *Player) PickItem(item I.Pickable) {
	p.items = append(p.items, item)
	item.SetOwner(p)
}

// Player interface implementation
func (p *Player) DropItem(item I.Pickable) {
	for i, it := range p.items {
		if it == item {
			item.SetOwner(nil)
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

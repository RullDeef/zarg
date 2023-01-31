package player

import (
	i "zarg/lib/model/interfaces"
)

const maxHealth = 100

type Player struct {
	user   i.User
	health int
	weapon i.Weapon
}

func NewPlayer(user i.User) *Player {
	return &Player{
		user:   user,
		health: maxHealth,
		weapon: nil,
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
func (p *Player) Damage(value int) {
	p.health -= value
	if p.health < 0 {
		p.health = 0
	}
}

// Entity interface implementation
func (p Player) Alive() bool {
	return p.health > 0
}

// Player interface implementation
func (p Player) Weapon() i.Weapon {
	return p.weapon
}

// Player interface implementation
func (p *Player) PickWeapon(w i.Weapon) {
	p.weapon = w
}

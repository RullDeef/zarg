package player

import (
	"zarg/lib/model/user"
	"zarg/lib/model/weapon"
)

const maxHealth = 100

type Player struct {
	user   *user.User
	Health int
	Weapon *weapon.Weapon
}

func NewPlayer(user *user.User) *Player {
	return &Player{
		user:   user,
		Health: maxHealth,
		Weapon: nil,
	}
}

func (p *Player) User() *user.User {
	return p.user
}

func (p *Player) Alive() bool {
	return p.Health > 0
}

func (p *Player) MakeDamage(val int) {
	p.Health -= val
	if p.Health < 0 {
		p.Health = 0
	}
}

func (p *Player) Heal(val int) {
	p.Health += val
	if p.Health > maxHealth {
		p.Health = maxHealth
	}
}

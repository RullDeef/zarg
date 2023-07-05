package player

import (
	"zarg/lib/model/entity"
	I "zarg/lib/model/interfaces"
)

const (
	maxHealth      = 100
	blockingFactor = 0.6
)

type Player struct {
	entity.BaseEntity
	user       I.User
	weapon     I.Weapon
	isBlocking bool
}

func NewPlayer(user I.User) *Player {
	var p *Player
	p = &Player{
		BaseEntity: entity.NewBase(user.FullName(), maxHealth, func() I.DamageStats {
			return p.AttackStats()
		}),
		user:       user,
		weapon:     nil,
		isBlocking: false,
	}
	return p
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
func (p *Player) Damage(dmg I.Damage) int {
	p.ForEachItem(func(item I.Pickable) {
		dmg = item.ModifyOngoingDamage(dmg)
	})

	totalDmg := p.CalcTotalDmg(dmg)

	if p.isBlocking {
		totalDmg = int(blockingFactor * float64(totalDmg))
	}

	// apply merged damage
	p.ApplyPureDamage(totalDmg)
	return totalDmg
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
func (p *Player) AttackStats() I.DamageStats {
	return p.weapon.AttackStats()
}

// Player interface implementation
func (p *Player) BlockAttack() {
	p.isBlocking = true
}

// Player interface implementation
func (p *Player) StopBlocking() {
	p.isBlocking = false
}

// Player interface implementation
func (p *Player) IsBlocking() bool {
	return p.isBlocking
}

func (p *Player) BeforeStartFight(interactor I.Interactor, friends I.EntityList, enemies I.EntityList) {

}

func (p *Player) AfterEndFight(interactor I.Interactor, friends I.EntityList, enemies I.EntityList) {

}

func (p *Player) BeforeDeath(Interactor I.Interactor, friends I.EntityList, enemies I.EntityList) {

}

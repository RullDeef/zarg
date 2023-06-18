package boss

import (
	I "zarg/lib/model/interfaces"
)

type Boss struct {
	currPhase *BossPhase

	interactor I.Interactor
}

func New(phases ...*BossPhase) *Boss {
	for i := 1; i < len(phases); i++ {
		phases[i-1].nextPhase = phases[i]
	}
	return &Boss{
		currPhase: phases[0],
	}
}

// Entity interface implementation
func (b *Boss) Name() string {
	return b.currPhase.Name()
}

// Entity interface implementation
func (b *Boss) Health() int {
	totalHealth := 0
	for bf := b.currPhase; bf != nil; bf = bf.nextPhase {
		totalHealth += bf.Health()
	}
	return totalHealth
}

// Entity interface implementation
func (b *Boss) Heal(value int) {
	b.currPhase.Heal(value)
}

// Entity interface implementation
func (b *Boss) Damage(dmg I.Damage) (res int) {
	res, b.currPhase = b.currPhase.Damage(dmg, b.interactor)
	return
}

// Entity interface implementation
func (b *Boss) Alive() bool {
	return b.currPhase.Alive()
}

// Enemy interface implementation
func (b *Boss) Attack(r float64) I.Damage {
	return b.currPhase.Attack(r)
}

// Enemy interface implementation
func (b *Boss) AttackStats() I.DamageStats {
	return b.currPhase.AttackStats()
}

func (b *Boss) CanDropItem(item I.Pickable) bool {
	return false
}

func (b *Boss) CanPickItem(item I.Pickable) bool {
	return false
}

func (b *Boss) PickItem(item I.Pickable) {

}

func (b *Boss) DropItem(item I.Pickable) {

}

func (b *Boss) ForEachItem(f func(I.Pickable)) {
}

func (b *Boss) ItemsCount() int {
	return 0
}

func (b *Boss) BeforeStartFight(interactor I.Interactor, friends I.EntityList, enemies I.EntityList) {
	b.interactor = interactor
}

func (b *Boss) AfterEndFight(I.Interactor, I.EntityList, I.EntityList) {

}

func (b *Boss) BeforeDeath(interactor I.Interactor, friends I.EntityList, enemies I.EntityList) {

}

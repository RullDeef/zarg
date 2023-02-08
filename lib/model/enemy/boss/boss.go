package boss

import I "zarg/lib/model/interfaces"

type Boss struct {
	currPhase *BossPhase
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
	return b.currPhase.name
}

// Entity interface implementation
func (b *Boss) Health() int {
	totalHealth := 0
	for bf := b.currPhase; bf != nil; bf = bf.nextPhase {
		totalHealth += bf.health
	}
	return totalHealth
}

// Entity interface implementation
func (b *Boss) Heal(value int) {
	b.currPhase.Heal(value)
}

// Entity interface implementation
func (b *Boss) Damage(ds I.DamageStats) (res int) {
	res, b.currPhase = b.currPhase.Damage(ds)
	return
}

// Entity interface implementation
func (b *Boss) Alive() bool {
	return b.currPhase.health > 0
}

// Enemy interface implementation
func (b *Boss) Attack() I.DamageStats {
	return b.currPhase.atack()
}

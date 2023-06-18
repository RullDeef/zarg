package boss

import (
	"zarg/lib/model/entity"
	I "zarg/lib/model/interfaces"
)

type BossPhase struct {
	entity.BaseEntity
	// attack        func() I.DamageStats
	nextPhase     *BossPhase
	onPhaseSwitch func(currPhase, nextPhase *BossPhase, interactor I.Interactor)
}

func NewPhase(name string, health int, attack func() I.DamageStats, phaseSwitch func(*BossPhase, *BossPhase, I.Interactor)) *BossPhase {
	return &BossPhase{
		BaseEntity: entity.NewBase(name, health, attack),
		// attack:        attack,
		onPhaseSwitch: phaseSwitch,
	}
}

func (bf *BossPhase) Damage(dmg I.Damage, interactor I.Interactor) (int, *BossPhase) {
	res := bf.BaseEntity.Damage(dmg)

	if !bf.Alive() && bf.nextPhase != nil {
		bf.onPhaseSwitch(bf, bf.nextPhase, interactor)
		bf = bf.nextPhase
	}
	return res, bf
}

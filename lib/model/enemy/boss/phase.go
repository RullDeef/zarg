package boss

import (
	"math/rand"
	I "zarg/lib/model/interfaces"
)

type BossPhase struct {
	name          string
	health        int
	maxHealth     int
	atack         func() I.DamageStats
	nextPhase     *BossPhase
	onPhaseSwitch func(currPhase, nextPhase *BossPhase)
}

func NewPhase(name string, health int, atack func() I.DamageStats, phaseSwitch func(*BossPhase, *BossPhase)) *BossPhase {
	return &BossPhase{
		name:          name,
		health:        health,
		maxHealth:     health,
		atack:         atack,
		onPhaseSwitch: phaseSwitch,
	}
}

func (bf *BossPhase) Heal(value int) {
	bf.health += value
	if bf.health > bf.maxHealth {
		bf.health = bf.maxHealth
	}
}

func (bf *BossPhase) Damage(ds I.DamageStats) (int, *BossPhase) {
	res := ds.Base
	if rand.Float32() < ds.CritChance {
		res = ds.Crit
	}
	bf.health -= res
	if bf.health < 0 {
		bf.health = 0
		if bf.nextPhase != nil {
			bf.onPhaseSwitch(bf, bf.nextPhase)
			bf = bf.nextPhase
		}
	}
	return res, bf
}

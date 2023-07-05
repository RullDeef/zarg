package damage

import I "zarg/lib/model/interfaces"

type BaseDamage struct {
	BaseDamageStats

	isCrit bool

	statusEffects []I.StatusEffect
}

func NewDamage(
	stats *BaseDamageStats,
	isCrit bool,
) *BaseDamage {
	return &BaseDamage{
		BaseDamageStats: *stats,
		isCrit:          isCrit,
		statusEffects:   nil,
	}
}

func NewDamageWithEffects(
	stats *BaseDamageStats,
	isCrit bool,
	statusEffects []I.StatusEffect,
) *BaseDamage {
	return &BaseDamage{
		BaseDamageStats: *stats,
		isCrit:          isCrit,
		statusEffects:   statusEffects,
	}
}

func (d *BaseDamage) IsCrit() bool {
	return d.isCrit
}

func (d *BaseDamage) StatusEffects() []I.StatusEffect {
	return d.statusEffects
}

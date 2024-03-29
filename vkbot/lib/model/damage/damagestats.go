package damage

import I "zarg/lib/model/interfaces"

type BaseDamageStats struct {
	typedDamages map[I.DamageType]int
	critChance   float64
	critFactor   float64

	statusEffectChances map[I.StatusEffect]float64
}

func NewStats(typedDamages map[I.DamageType]int, critChance, critFactor float64) *BaseDamageStats {
	return &BaseDamageStats{
		typedDamages,
		critChance,
		critFactor,

		make(map[I.StatusEffect]float64),
	}
}

func NewStatsWithEffects(
	typedDamages map[I.DamageType]int,
	critChance float64,
	critFactor float64,
	statusEffectChances map[I.StatusEffect]float64,
) *BaseDamageStats {
	return &BaseDamageStats{
		typedDamages,
		critChance,
		critFactor,
		statusEffectChances,
	}
}

// DamageStats interface implementation
func (d *BaseDamageStats) TypedDamages() map[I.DamageType]int {
	return d.typedDamages
}

// DamageStats interface implementation
func (d *BaseDamageStats) CritChance() float64 {
	return d.critChance
}

// DamageStats interface implementation
func (d *BaseDamageStats) CritFactor() float64 {
	return d.critFactor
}

func (d *BaseDamageStats) StatusEffectChances() map[I.StatusEffect]float64 {
	return d.statusEffectChances
}

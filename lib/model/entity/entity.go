package entity

import (
	"fmt"
	"math"
	"zarg/lib/model/damage"
	I "zarg/lib/model/interfaces"
	"zarg/lib/utils"
)

type BaseEntity struct {
	name      string
	health    int
	maxHealth int
	items     []I.Pickable

	// damage type weaknesses/resists
	damageFactors map[I.DamageType]float64 // +x weak | -x resist

	statusEffects []I.StatusEffect
	damageStats   func() I.DamageStats
	// attack        func() I.Damage

	// special actions
	BeforeStartFightFunc func(I.Interactor, I.EntityList, I.EntityList)
	AfterEndFightFunc    func(I.Interactor, I.EntityList, I.EntityList)
	BeforeDeathFunc      func(I.Interactor, I.EntityList, I.EntityList)
}

func New(name string, maxHealth int, damageStats func() I.DamageStats) *BaseEntity {
	return &BaseEntity{
		name:        name,
		health:      maxHealth,
		maxHealth:   maxHealth,
		damageStats: damageStats,
		// attack:      attack,
	}
}

func NewBase(name string, maxHealth int, damageStats func() I.DamageStats) BaseEntity {
	return BaseEntity{
		name:        name,
		health:      maxHealth,
		maxHealth:   maxHealth,
		damageStats: damageStats,
	}
}

// Entity interface implementation
func (e BaseEntity) Name() string {
	return e.name
}

// Entity interface implementation
func (e BaseEntity) Health() int {
	return e.health
}

// Entity interface implementation
func (e *BaseEntity) Heal(value int) {
	if value < 0 {
		panic(fmt.Sprintf("heal value must not be negative. Got: %d", value))
	}

	e.health += value
	if e.health > e.maxHealth {
		e.health = e.maxHealth
	}
}

// Entity interface implementation
func (e BaseEntity) Alive() bool {
	return e.health > 0
}

// Entity interface implementation
func (e *BaseEntity) Damage(dmg I.Damage) int {
	e.ForEachItem(func(item I.Pickable) {
		dmg = item.ModifyOngoingDamage(dmg)
	})

	totalDmg := e.CalcTotalDmg(dmg)

	// apply merged damage
	e.ApplyPureDamage(totalDmg)

	return totalDmg
}

func (e *BaseEntity) ApplyPureDamage(val int) {
	e.health -= val
	if e.health < 0 {
		e.health = 0
	} else if e.health > e.maxHealth {
		e.health = e.maxHealth
	}
}

func (e *BaseEntity) CalcTotalDmg(dmg I.Damage) int {
	// apply weaknesses/resists. Merge all types of damage
	var totalDmgFloat = 0.0
	for dmgType, val := range dmg.TypedDamages() {
		factor := 1.0 + e.damageFactors[dmgType]
		totalDmgFloat += factor * float64(val)
	}

	// apply crit
	if dmg.IsCrit() {
		totalDmgFloat *= dmg.CritFactor()
	}

	return int(math.Round(float64(totalDmgFloat)))
}

// Entity interface implementation
func (e *BaseEntity) AttackStats() I.DamageStats {
	return e.damageStats()
}

// Entity interface implementation
func (e *BaseEntity) Attack(r float64) I.Damage {
	stats := e.AttackStats().(*damage.BaseDamageStats)
	var dmg I.Damage = damage.NewDamage(stats, r < stats.CritChance())
	e.ForEachItem(func(item I.Pickable) {
		dmg = item.ModifyOutgoingDamage(dmg)
	})
	return dmg
}

// Entity interface implementation
func (e BaseEntity) CanPickItem(item I.Pickable) bool {
	return true
}

// Entity interface implementation
func (e *BaseEntity) CanDropItem(item I.Pickable) bool {
	return item.Owner() == e
}

// Entity interface implementation
func (e *BaseEntity) PickItem(item I.Pickable) {
	if !e.CanPickItem(item) {
		panic(fmt.Sprintf("tried to pick item %v, but cannot", item))
	}

	item.SetOwner(e)
	e.items = append(e.items, item)
}

// Entity interface implementation
func (e *BaseEntity) DropItem(item I.Pickable) {
	if !e.CanDropItem(item) {
		panic(fmt.Sprintf("tried to drop item %v, but cannot", item))
	}

	totalItems := len(e.items)
	index := utils.FindFirstOrPanic(totalItems, func(i int) bool {
		return e.items[i] == item
	})

	item.SetOwner(nil)
	e.items = append(e.items[:index], e.items[index+1:]...)
}

// Entity interface implementation
func (e BaseEntity) ForEachItem(f func(I.Pickable)) {
	for _, item := range e.items {
		f(item)
	}
}

// Entity interface implementation
func (e BaseEntity) ItemsCount() int {
	return len(e.items)
}

func (e *BaseEntity) BeforeStartFight(interactor I.Interactor, friends I.EntityList, enemies I.EntityList) {
	if e.BeforeStartFightFunc != nil {
		e.BeforeStartFightFunc(interactor, friends, enemies)
	}
}
func (e *BaseEntity) AfterEndFight(interactor I.Interactor, friends I.EntityList, enemies I.EntityList) {
	if e.AfterEndFightFunc != nil {
		e.AfterEndFightFunc(interactor, friends, enemies)
	}
}
func (e *BaseEntity) BeforeDeath(interactor I.Interactor, friends I.EntityList, enemies I.EntityList) {
	if e.BeforeDeathFunc != nil {
		e.BeforeDeathFunc(interactor, friends, enemies)
	}
}

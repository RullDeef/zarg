package enemy

import (
	"math/rand"
	I "zarg/lib/model/interfaces"
)

var enemyNames = []string{
	"Гигантская крыса", "Выползень", "Грызяк", "Горгулья",
	"Злой огонёк", "Летающая черепушка", "Зубастый ящер",
	"Чёрт",
}

type Enemy struct {
	name      string
	health    int
	maxHealth int
	attack    func() I.DamageStats
}

func New(name string, health int, attack func() I.DamageStats) *Enemy {
	return &Enemy{
		name:      name,
		health:    health,
		maxHealth: health,
		attack:    attack,
	}
}

func Random(health int, attack func() I.DamageStats) *Enemy {
	name := enemyNames[rand.Intn(len(enemyNames))]
	return New(name, health, attack)
}

// Entity interface implementation
func (e Enemy) Name() string {
	return e.name
}

// Entity interface implementation
func (e *Enemy) Health() int {
	return e.health
}

// Entity interface implementation
func (e *Enemy) Heal(value int) {
	e.health += value
	if e.health > e.maxHealth {
		e.health = e.maxHealth
	}
}

// Entity interface implementation
func (e *Enemy) Damage(dmg I.DamageStats) int {
	val := dmg.Base
	if rand.Float32() < dmg.CritChance {
		val = dmg.Crit
	}
	e.health -= val
	if e.health < 0 {
		e.health = 0
	}
	return val
}

// Entity interface implementation
func (e *Enemy) Alive() bool {
	return e.health > 0
}

// Enemy interface implementation
func (e *Enemy) Attack() I.DamageStats {
	dmg := e.attack()
	dmg.Producer = e
	return dmg
}

package enemy

import (
	"math/rand"
	"zarg/lib/model/entity"
	I "zarg/lib/model/interfaces"
)

var enemyNames = []string{
	"Гигантская крыса", "Выползень", "Грызяк", "Горгулья",
	"Злой огонёк", "Летающая черепушка", "Зубастый ящер",
	"Чёрт",
}

type Enemy struct {
	entity.BaseEntity
}

func New(name string, health int, attack func() I.DamageStats) *Enemy {
	return &Enemy{
		BaseEntity: entity.NewBase(name, health, attack),
	}
}

func Random(health int, attack func() I.DamageStats) *Enemy {
	name := enemyNames[rand.Intn(len(enemyNames))]
	return New(name, health, attack)
}

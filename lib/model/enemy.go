package model

import "math/rand"

var enemyNames = []string{
	"Гигантская крыса", "Выползень", "Грызяк",
}

type Enemy struct {
	Name   string
	Health int
	Attack int
	Tier   int
}

func NewEnemy(name string, health int, attack int, tier int) *Enemy {
	return &Enemy{
		Name:   name,
		Health: health,
		Attack: attack,
		Tier:   tier,
	}
}

func RandomEnemy(tier int, attackMean, attackDiff int) *Enemy {
	name := enemyNames[rand.Intn(len(enemyNames))]
	health := 20

	attackMin := attackMean - attackDiff
	attackMax := attackMean + attackDiff
	attack := attackMin + rand.Intn(attackMax-attackMin+1)

	return NewEnemy(name, health, attack, tier)
}

func (e *Enemy) MakeDamage(val int) {
	e.Health -= val
	if e.Health < 0 {
		e.Health = 0
	}
}

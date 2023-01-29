package model

import "math/rand"

var enemyNames = map[int][]string{
	1: {"Гигантская крыса", "Выползень", "Грызяк"},
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
	name := enemyNames[tier][rand.Intn(len(enemyNames[tier]))]
	health := 20

	attackMin := attackMean - attackDiff
	attackMax := attackMean + attackDiff
	attack := attackMin + rand.Intn(attackMax-attackMin+1)

	return NewEnemy(name, health, attack, tier)
}

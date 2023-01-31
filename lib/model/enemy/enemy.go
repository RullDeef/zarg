package enemy

import "math/rand"

var enemyNames = []string{
	"Гигантская крыса", "Выползень", "Грызяк",
}

type Enemy struct {
	name        string
	health      int
	maxHealth   int
	attackPower int
	onAttack    func(*Enemy)
}

func New(name string, health int, attackPower int, attack func(*Enemy)) *Enemy {
	return &Enemy{
		name:      name,
		health:    health,
		maxHealth: health,
		onAttack:  attack,
	}
}

func Random(attackMean int, attackDiff int, attack func(*Enemy)) *Enemy {
	name := enemyNames[rand.Intn(len(enemyNames))]
	health := 20

	attackMin := attackMean - attackDiff
	attackMax := attackMean + attackDiff
	attackPower := attackMin + rand.Intn(attackMax-attackMin+1)

	return New(name, health, attackPower, attack)
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
func (e *Enemy) Damage(value int) {
	e.health -= value
	if e.health < 0 {
		e.health = 0
	}
}

// Entity interface implementation
func (e *Enemy) Alive() bool {
	return e.health > 0
}

// Enemy interface implementation
func (e *Enemy) Attack() {
	e.onAttack(e)
}

func (e *Enemy) AttackPower() int {
	return e.attackPower
}

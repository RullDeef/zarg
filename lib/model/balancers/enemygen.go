package balancers

import (
	"math"
	"math/rand"
	I "zarg/lib/model/interfaces"
	"zarg/lib/model/player/squad"
)

type EnemyBalancer struct {
	floorNumber int
	players     *squad.PlayerSquad
}

func NewEnemyBalancer(floorNumber int, players *squad.PlayerSquad) *EnemyBalancer {
	return &EnemyBalancer{
		floorNumber: floorNumber,
		players:     players,
	}
}

// EnemyBalancer interface implementation
func (eb *EnemyBalancer) EnemiesCount() int {
	return 1 + eb.players.LenAlive()/2 + rand.Intn(2)
}

// EnemyBalancer interface implementation
func (eb *EnemyBalancer) Health() (min, max int) {
	power := 0.0
	eb.players.ForEachAlive(func(player I.Player) {
		atk := player.Weapon().Attack()
		// p := float64(atk.CritChance)
		power += float64(atk.Base) // *(1.0-p) + float64(atk.Crit)*p
	})

	meanPower := power / float64(eb.players.LenAlive())

	// enemy must die in 2-3 attacks from player
	min = int(math.Ceil(meanPower * 2.0))
	max = int(math.Ceil(meanPower * 3.0))
	return
}

// EnemyBalancer interface implementation
func (eb *EnemyBalancer) Attack() (min, max int) {
	// floor 1: 9-11 attacks to player death
	// floor 2: 7-9 attacks to player death
	// floor 3: 5-7 attacks to player death

	maxHealth := 50.0
	eb.players.ForEachAlive(func(player I.Player) {
		maxHealth = math.Max(maxHealth, float64(player.Health()))
	})

	switch eb.floorNumber {
	case 1:
		min = int(math.Ceil(maxHealth / 11.0))
		max = int(math.Ceil(maxHealth / 9.0))
	case 2:
		min = int(math.Ceil(maxHealth / 9.0))
		max = int(math.Ceil(maxHealth / 7.0))
	case 3:
		min = int(math.Ceil(maxHealth / 7.0))
		max = int(math.Ceil(maxHealth / 5.0))
	default:
		panic("attack for floor under 3 not accounted")
	}
	return
}

// EnemyBalancer interface implementation
func (eb *EnemyBalancer) ExtraCrit() float32 {
	switch eb.floorNumber {
	case 1:
		return 1.2
	case 2:
		return 1.5
	case 3:
		return 1.8
	default:
		panic("extra crit for floor under 3 not accounted")
	}
}

// EnemyBalancer interface implementation
func (eb *EnemyBalancer) CritChance() float32 {
	switch eb.floorNumber {
	case 1:
		return 0.1
	case 2:
		return 0.1 * float32(1+rand.Intn(2))
	case 3:
		return 0.1 * float32(2+rand.Intn(2))
	default:
		panic("crit chance for floor under 3 not accounted")
	}
}

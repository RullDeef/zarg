package interfaces

import (
	"context"
	"fmt"
)

type User interface {
	ID() int
	FirstName() string
	LastName() string
	FullName() string
}

type Entity interface {
	Name() string

	Health() int
	Heal(value int)
	Alive() bool

	// returns actual damage taken
	Damage(dmg Damage) int

	AttackStats() DamageStats

	// controllable crits
	Attack(rand float64) Damage

	// pickup management

	CanPickItem(Pickable) bool
	CanDropItem(Pickable) bool
	PickItem(Pickable)
	DropItem(Pickable)

	ForEachItem(func(Pickable))
	ItemsCount() int

	// special action informing
	BeforeStartFight(interactor Interactor, friends EntityList, enemies EntityList)
	AfterEndFight(interactor Interactor, friends EntityList, enemies EntityList)
	BeforeDeath(Interactor Interactor, friends EntityList, enemies EntityList)

	// status effects
	StatusEffects() []StatusEffect
	AddStatusEffect(StatusEffect)

	// returns amount of turns to be proceed
	ApplyStatusEffectsBeforeMyTurn(Interactor Interactor, friends EntityList, enemies EntityList) int
}

type EntityList interface {
	Len() int
	LenAlive() int

	ForEach(func(Entity))
	ForEachAlive(func(Entity))

	Has(Entity) bool
}

type Player interface {
	User
	Entity

	Weapon() Weapon
	PickWeapon(Weapon)

	BlockAttack()
	StopBlocking()
	IsBlocking() bool
}

// weapon types
const (
	WeaponTypeSlicing  = "режущее"
	WeaponTypeStabbing = "колющее"
	WeaponTypeCrushing = "дробящее"
	WeaponTypeMagical  = "магическое"
	WeaponTypeSpecial  = "особое"
)

type Weapon interface {
	Pickable

	Kind() string
	AttackStats() DamageStats
}

// type DamageStats struct {
// 	Base       int
// 	Crit       int
// 	CritChance float32
// }

type DamageType int

const (
	DamageType1 DamageType = iota
	DamageType2
	DamageType3
	DamageType4
)

type DamageStats interface {
	// returns mapping from damage type to damage value
	TypedDamages() map[DamageType]int

	CritChance() float64
	CritFactor() float64 // > 1

	StatusEffectChances() map[StatusEffect]float64
}

type Damage interface {
	DamageStats

	IsCrit() bool
	StatusEffects() []StatusEffect
}

type StatusEffect struct {
	Name        string
	Description string

	// for now it just amount of turns left
	TimeLeft int // TODO: add custom ~ActionTime struct
}

func (s StatusEffect) ShortDescribeWithChance(chance float64) string {
	return fmt.Sprintf("%sx%d:%.0f%%", s.Name, s.TimeLeft, chance)
}

func (s StatusEffect) DescribeWithChance(chance float64) string {
	return fmt.Sprintf("%sx%d (%s) %.0f%%", s.Name, s.TimeLeft, s.Description, chance)
}

var (
	StatusEffectStun = func(n int) StatusEffect {
		return StatusEffect{
			Name:        "🌀",
			Description: "оглушение. Пропуск хода",
			TimeLeft:    n,
		}
	}
	StatusEffectAgility = func(n int) StatusEffect {
		return StatusEffect{
			Name:        "⚡",
			Description: "проворность. Дополнительный ход",
			TimeLeft:    n,
		}
	}
	StatusEffectRegeneration = func(n int) StatusEffect {
		return StatusEffect{
			Name:        "💞",
			Description: "регенерация. +1❤",
			TimeLeft:    n,
		}
	}
	StatusEffectBleeding = func(n int) StatusEffect {
		return StatusEffect{
			Name:        "❣",
			Description: "кровотечение. -1❤",
			TimeLeft:    n,
		}
	}
)

type Pickable interface {
	Name() string
	Description() string

	Owner() Entity
	SetOwner(Entity)

	ModifyOngoingDamage(Damage) Damage
	ModifyOutgoingDamage(Damage) Damage
}

type Usable interface {
	Pickable

	Use()
	IsUsed() bool
}

type Consumable interface {
	Pickable

	UsesLeft() int
	Consume()
}

type Enemy interface {
	Entity
}

type WeaponShowcase interface {
	WeaponsInfo() string
	HasMadePick(Player) bool

	// returns theese values:
	// bool - wether or not pick was successfull
	// Weapon - that was tried to pick
	// Player - who now owns this weapon
	TryPick(p Player, opt int) (bool, Weapon, Player)

	ConfirmPick()
}

type UserMessage interface {
	User() User
	Message() string
}

type Interactor interface {
	Printf(fmt string, args ...any)

	// gets a messages from chat.
	Receive(ctx context.Context, f func(UserMessage)) error
}

type FloorGenBalancer interface {
	TreasureRoomsCount() int
	EnemyRoomsCount() int
	TrapRoomsCount() int

	ItemsInTreasureRoomCount() int

	EnemyBalancer() EnemyBalancer
}

type EnemyBalancer interface {
	EnemiesCount() int

	Health() (min, max int)
	Attack() (min, max int)
	ExtraCrit() float64 // percent from attack (> 1.0)
	CritChance() float64
}

type FightManager interface {
	// returns 1, if sqiad1 wins, 2 if squad2 wins, 0 - draw
	PerformFight(squad1, squad2 EntityList) int
}

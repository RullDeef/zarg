package interfaces

import "context"

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
	// returns actual damage taken
	Damage(DamageStats) int
	Alive() bool
}

type Player interface {
	User
	Entity

	Weapon() Weapon
	PickWeapon(Weapon)

	Attack() DamageStats

	BlockAttack()
	IsBlocking() bool

	PickItem(Pickable)
	DropItem(Pickable)
	ForEachItem(func(Pickable))
}

type PlayerList interface {
	Len() int
	LenAlive() int

	ForEach(func(Player))
	ForEachAlive(func(Player))
}

type Weapon interface {
	Pickable

	Description() string

	Attack() DamageStats
}

type DamageEmitor interface {
	Name() string
}

type DamageStats struct {
	Producer   DamageEmitor
	Base       int
	Crit       int
	CritChance float32
}

type Pickable interface {
	Name() string

	Owner() Player
	SetOwner(Player)

	ModifyOngoingDamage(DamageStats) DamageStats
	ModifyOutgoingDamage(DamageStats) DamageStats
}

type Usable interface {
	Pickable

	Description() string
	Use()
	IsUsed() bool
}

type Consumable interface {
	Pickable

	Description() string

	UsesLeft() int
	Consume()
}

type Enemy interface {
	Entity

	Attack() DamageStats
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

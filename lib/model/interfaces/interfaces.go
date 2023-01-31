package interfaces

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
	Damage(value int)
	Alive() bool
}

type Player interface {
	User
	Entity

	Weapon() Weapon
	PickWeapon(Weapon)
}

type PlayerList interface {
	Len() int
	LenAlive() int

	ForEach(func(Player))
	ForEachAlive(func(Player))
}

type Weapon interface {
	Title() string
	Description() string

	Attack() int
}

type Pickable interface {
	Name() string

	PickUp(Player)
	Owner() Player
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

	AttackPower() int
	Attack()
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

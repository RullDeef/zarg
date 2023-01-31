package weapon

import (
	"fmt"
	"log"
	"math/rand"
	I "zarg/lib/model/interfaces"
)

type WeaponKind int

type Weapon struct {
	name   string
	attack int
	Kind   WeaponKind
	owner  I.Player
}

var (
	weaponKinds = []string{
		"режущее",
		"колющее",
		"дробящее",
		"магическое",
		"особое",
	}
	weaponNames = map[WeaponKind][]string{
		0: {"Ржавый меч", "Окровавленная Лопата"},
		1: {"Садовые вилы"},
		2: {"Монтировка"},
		3: {"Файербол"},
		4: {"Хлыст боли"},
	}
)

func FistsWeapon(attackMean, attackDiff int) *Weapon {
	return &Weapon{
		name:   "Кулаки",
		attack: attackMean - attackDiff + 2*rand.Int()%(attackDiff+1),
		Kind:   4,
		owner:  nil,
	}
}

func RandomWeapon(tier int, attackMean, attackDiff int) *Weapon {
	kind := WeaponKind(rand.Int() % len(weaponKinds))
	name := weaponNames[kind][rand.Int()%len(weaponNames[kind])]

	return &Weapon{
		name:   name,
		attack: attackMean - attackDiff + 2*rand.Int()%(attackDiff+1),
		Kind:   kind,
		owner:  nil,
	}
}

func RandomWeapons(n int, tier int, attackMean, attackDiff int) []*Weapon {
	var w []*Weapon
	for i := 0; i < n; i++ {
		w = append(w, RandomWeapon(tier, attackMean, attackDiff))
	}
	return w
}

// Weapon interface implementation
func (w Weapon) Description() string {
	return fmt.Sprintf("Атака - %d", w.attack)
}

// Weapon interface implementation
func (w Weapon) Title() string {
	return w.name
}

// Weapon interface implementation
func (w *Weapon) SetOwner(p I.Player) {
	w.owner = p
}

// Weapon interface implementation
func (w Weapon) Attack() I.DamageStats {
	if w.owner == nil {
		log.Panicf("owner for weapon %+v is not set!", w)
	}

	return I.DamageStats{
		Producer:   w.owner,
		Base:       w.attack,
		Crit:       w.attack + 10,
		CritChance: 0.07,
	}
}

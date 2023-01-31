package weapon

import (
	"fmt"
	"math/rand"
)

type WeaponKind int

type Weapon struct {
	name   string
	attack int
	Kind   WeaponKind
	Tier   int
}

var (
	tierStrings = []string{
		"I",
		"II",
		"III",
		"IV",
		"V",
	}
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
		Tier:   -1,
	}
}

func RandomWeapon(tier int, attackMean, attackDiff int) *Weapon {
	kind := WeaponKind(rand.Int() % len(weaponKinds))
	name := weaponNames[kind][rand.Int()%len(weaponNames[kind])]

	return &Weapon{
		name:   name,
		attack: attackMean - attackDiff + 2*rand.Int()%(attackDiff+1),
		Kind:   kind,
		Tier:   tier,
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
	name := w.name
	if w.Tier >= 0 {
		name += " " + tierStrings[w.Tier]
	}
	return name
}

// Weapon interface implementation
func (w Weapon) Attack() int {
	return w.attack
}

package model

import (
	"fmt"
	"math/rand"
)

type WeaponKind int

type Weapon struct {
	Name   string
	Attack int
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
		Name:   "Кулаки",
		Attack: attackMean - attackDiff + 2*rand.Int()%(attackDiff+1),
		Kind:   4,
		Tier:   -1,
	}
}

func RandomWeapon(tier int, attackMean, attackDiff int) *Weapon {
	kind := WeaponKind(rand.Int() % len(weaponKinds))
	name := weaponNames[kind][rand.Int()%len(weaponNames[kind])]

	return &Weapon{
		Name:   name,
		Attack: attackMean - attackDiff + 2*rand.Int()%(attackDiff+1),
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

func (w Weapon) Summary() string {
	name := w.Name
	if w.Tier >= 0 {
		name += " " + tierStrings[w.Tier]
	}
	return fmt.Sprintf("%s. Атака - %d", name, w.Attack)
}

func (w Weapon) SummaryShort() string {
	name := w.Name
	if w.Tier >= 0 {
		name += " " + tierStrings[w.Tier]
	}
	return name
}

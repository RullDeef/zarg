package weapon

import (
	"fmt"
	"math/rand"
	"zarg/lib/model/damage"
	I "zarg/lib/model/interfaces"
)

type WeaponKind int

const (
	defaultCritFactor = 1.2
	fistsBaseAttack   = 5
)

type Weapon struct {
	name         string
	typedDamages map[I.DamageType]int
	critChance   float64
	critFactor   float64
	kind         WeaponKind
	owner        I.Entity
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
		0: {"Ржавый меч", "Окровавленная Лопата", "Топор лесоруба"},
		1: {"Садовые вилы", "Кривая пика", "Лук в колючку"},
		2: {"Монтировка", "Строительная кувалда", "Дубина с зубами грызяка"},
		3: {"Файербол", "Молния", "Посох колдуна"},
		4: {"Хлыст боли", "Шипастый цеп"},
	}
)

func FistsWeapon() *Weapon {
	typedDamages := make(map[I.DamageType]int)
	typedDamages[I.DamageType1] = fistsBaseAttack

	return &Weapon{
		name:         "Кулаки",
		typedDamages: typedDamages,
		critChance:   0.05,
		critFactor:   2.0,
		kind:         4,
		owner:        nil,
	}
}

func RandomWeapon(attackMean, attackDiff int) *Weapon {
	kind := WeaponKind(rand.Int() % len(weaponKinds))
	name := weaponNames[kind][rand.Int()%len(weaponNames[kind])]

	typedDamages := make(map[I.DamageType]int)
	typedDamages[I.DamageType1] = attackMean - attackDiff + rand.Intn(2*attackDiff+1)

	return &Weapon{
		name:         name,
		typedDamages: typedDamages,
		critChance:   0.1,
		critFactor:   defaultCritFactor,
		kind:         kind,
		owner:        nil,
	}
}

func RandomWeapons(n int, attackMean, attackDiff int) []*Weapon {
	var w []*Weapon
	for i := 0; i < n; i++ {
		w = append(w, RandomWeapon(attackMean, attackDiff))
	}
	return w
}

// Pickable interface implementation
func (w *Weapon) Name() string {
	return w.name
}

// Pickable interface implementation
func (w Weapon) Description() string {
	return fmt.Sprintf("%d🗡", w.typedDamages[I.DamageType1])
}

// Pickable interface implementation
func (w *Weapon) Owner() I.Entity {
	return w.owner
}

// Pickable interface implementation
func (w *Weapon) SetOwner(entity I.Entity) {
	w.owner = entity
}

// Pickable interface implementation
func (w *Weapon) ModifyOngoingDamage(dmg I.Damage) I.Damage {
	return dmg
}

// Pickable interface implementation
func (w *Weapon) ModifyOutgoingDamage(dmg I.Damage) I.Damage {
	return dmg
}

// Weapon interface implementation
func (w Weapon) AttackStats() I.DamageStats {
	return damage.NewStats(w.typedDamages, w.critChance, w.critFactor)
}

// Weapon interface implementation
func (w Weapon) Kind() string {
	return weaponKinds[w.kind]
}

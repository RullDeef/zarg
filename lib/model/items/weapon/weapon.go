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

	statusEffectChances map[I.StatusEffect]float64
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
		statusEffectChances: map[I.StatusEffect]float64{
			I.StatusEffectStun(1): 0.5,
		},
	}
}

func RandomWeapon(attackMean, attackDiff int) *Weapon {
	kind := WeaponKind(rand.Int() % len(weaponKinds))
	name := weaponNames[kind][rand.Int()%len(weaponNames[kind])]

	typedDamages := make(map[I.DamageType]int)
	typedDamages[I.DamageType1] = attackMean - attackDiff + rand.Intn(2*attackDiff+1)

	statusEffects := make(map[I.StatusEffect]float64)

	// TODO: move constants out into status-effect-balancer kind of
	if rand.Float64() < 0.2 {
		time := rand.Intn(2) + 1
		chance := float64(rand.Intn(5)+3) / 20
		statusEffects[I.StatusEffectStun(time)] = chance
	}
	if rand.Float64() < 0.2 {
		time := rand.Intn(2) + 1
		chance := float64(rand.Intn(5)+3) / 20
		statusEffects[I.StatusEffectBleeding(time)] = chance
	}
	if rand.Float64() < 0.2 {
		time := rand.Intn(2) + 1
		chance := float64(rand.Intn(5)+3) / 20
		statusEffects[I.StatusEffectBurning(time)] = chance
	}
	if rand.Float64() < 0.2 {
		time := rand.Intn(2) + 1
		chance := float64(rand.Intn(5)+3) / 20
		statusEffects[I.StatusEffectFreezing(time)] = chance
	}
	if rand.Float64() < 0.2 {
		time := rand.Intn(2) + 1
		chance := float64(rand.Intn(5)+3) / 20
		statusEffects[I.StatusEffectWeakness(time)] = chance
	}

	return &Weapon{
		name:                name,
		typedDamages:        typedDamages,
		critChance:          0.1,
		critFactor:          defaultCritFactor,
		kind:                kind,
		owner:               nil,
		statusEffectChances: statusEffects,
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
	msg := fmt.Sprintf("%d🗡", w.typedDamages[I.DamageType1])

	for effect, chance := range w.statusEffectChances {
		msg += fmt.Sprintf("[%dx%s:%.0f%%]", effect.TimeLeft, effect.Name, chance*100)
	}

	return msg
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
	return damage.NewStatsWithEffects(w.typedDamages, w.critChance, w.critFactor, w.statusEffectChances)
}

// Weapon interface implementation
func (w Weapon) Kind() string {
	return weaponKinds[w.kind]
}

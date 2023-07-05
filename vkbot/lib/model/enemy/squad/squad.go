package squad

import (
	"container/list"
	"fmt"
	"log"
	"strings"
	I "zarg/lib/model/interfaces"
)

type EnemySquad struct {
	list *list.List
	iter *list.Element
}

func New(n int, builder func() I.Enemy) *EnemySquad {
	list := list.New()

	for i := 0; i < n; i++ {
		list.PushBack(builder())
	}

	return &EnemySquad{
		list: list,
		iter: nil,
	}
}

func (es *EnemySquad) Len() int {
	return es.list.Len()
}

func (es *EnemySquad) ForEach(f func(I.Entity)) {
	for node := es.list.Front(); node != nil; node = node.Next() {
		f(node.Value.(I.Enemy))
	}
}

func (es *EnemySquad) ForEachAlive(f func(I.Entity)) {
	for node := es.list.Front(); node != nil; node = node.Next() {
		e := node.Value.(I.Enemy)
		if e.Alive() {
			f(e)
		}
	}
}

func (es *EnemySquad) Has(enemy I.Entity) bool {
	for node := es.list.Front(); node != nil; node = node.Next() {
		e := node.Value.(I.Entity)
		if e == enemy {
			return true
		}
	}
	return false
}

func (es *EnemySquad) CompactInfo() string {
	res := ""

	for node := es.list.Front(); node != nil; node = node.Next() {
		e := node.Value.(I.Enemy)
		if e.Alive() {
			atk := e.AttackStats()
			effects := es.EffectsCompactInfo(e)
			res += fmt.Sprintf("- %s (%d❤ %d🗡) %s\n", e.Name(), e.Health(), atk.TypedDamages()[I.DamageType1], effects)
		} else {
			res += fmt.Sprintf("- %s 💀\n", e.Name())
		}
	}

	return res
}

func (es *EnemySquad) EffectsCompactInfo(e I.Enemy) string {
	var effects []string
	for _, eff := range e.StatusEffects() {
		effects = append(effects, fmt.Sprintf("%sx%d", eff.Name, eff.TimeLeft))
	}
	if len(effects) == 0 {
		return ""
	} else {
		return "[" + strings.Join(effects, "|") + "]"
	}
}

func (es *EnemySquad) LenAlive() int {
	res := 0

	for node := es.list.Front(); node != nil; node = node.Next() {
		e := node.Value.(I.Enemy)
		if e.Alive() {
			res += 1
		}
	}

	return res
}

func (es *EnemySquad) ChooseNext() I.Enemy {
	if es.LenAlive() == 0 {
		log.Panic("all enemies dead!")
	}

	if es.iter == nil {
		es.iter = es.list.Front()
	} else {
		es.iter = es.iter.Next()
		if es.iter == nil {
			es.iter = es.list.Front()
		}
	}

	for !es.iter.Value.(I.Enemy).Alive() {
		es.iter = es.iter.Next()
		if es.iter == nil {
			es.iter = es.list.Front()
		}
	}

	return es.iter.Value.(I.Enemy)
}

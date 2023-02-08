package squad

import (
	"container/list"
	"fmt"
	"log"
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

func (es *EnemySquad) ForEach(f func(I.Enemy)) {
	for node := es.list.Front(); node != nil; node = node.Next() {
		f(node.Value.(I.Enemy))
	}
}

func (es *EnemySquad) ForEachAlive(f func(I.Enemy)) {
	for node := es.list.Front(); node != nil; node = node.Next() {
		e := node.Value.(I.Enemy)
		if e.Alive() {
			f(e)
		}
	}
}

func (es *EnemySquad) CompactInfo() string {
	res := ""

	for node := es.list.Front(); node != nil; node = node.Next() {
		e := node.Value.(I.Enemy)
		if e.Alive() {
			atk := e.Attack()
			res += fmt.Sprintf("- %s (%d❤ %d🗡)\n", e.Name(), e.Health(), atk.Base)
		} else {
			res += fmt.Sprintf("- %s 💀\n", e.Name())
		}
	}

	return res
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

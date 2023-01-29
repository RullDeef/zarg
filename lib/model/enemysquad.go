package model

import (
	"container/list"
	"fmt"
	"log"
)

type EnemySquad struct {
	list *list.List
	iter *list.Element
}

func NewEnemySquad(n int, builder func() *Enemy) *EnemySquad {
	list := list.New()

	for i := 0; i < n; i++ {
		list.PushBack(builder())
	}

	return &EnemySquad{
		list: list,
		iter: nil,
	}
}

func (es *EnemySquad) Info() string {
	res := ""

	for i, node := 1, es.list.Front(); node != nil; i, node = i+1, node.Next() {
		e := node.Value.(*Enemy)
		res += fmt.Sprintf("%d) %s (HP:%d, Атака:%d)\n", i, e.Name, e.Health, e.Attack)
	}

	return res
}

func (es *EnemySquad) LenAlive() int {
	res := 0

	for i, node := 1, es.list.Front(); node != nil; i, node = i+1, node.Next() {
		e := node.Value.(*Enemy)
		if e.Health > 0 {
			res += 1
		}
	}

	return res
}

func (es *EnemySquad) ChooseNext() *Enemy {
	if es.LenAlive() == 0 {
		log.Fatal("all enemies dead!")
	}

	if es.iter == nil {
		es.iter = es.list.Front()
	} else {
		es.iter = es.iter.Next()
		if es.iter == nil {
			es.iter = es.list.Front()
		}
	}

	for es.iter.Value.(*Enemy).Health == 0 {
		es.iter = es.iter.Next()
		if es.iter == nil {
			es.iter = es.list.Front()
		}
	}

	return es.iter.Value.(*Enemy)
}

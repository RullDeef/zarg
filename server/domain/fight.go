package domain

import (
	"container/ring"
	"context"
	"errors"
	"math/rand"
	"slices"
)

var (
	ErrEntityIsDead   = errors.New("entity is dead")     // сущность мертва на начало боя
	ErrEntityRepeated = errors.New("entity is repeated") // сущность повторяется
	ErrOrderInvalid   = errors.New("order is invalid")   // порядок существ некорректен

	ErrFightIsOver    = errors.New("fight is over")
	ErrFightIsNotOver = errors.New("fight is not over")

	// ErrEntityEscaped - существо сбежало во время своего хода.
	// Данную ошибку необходимо возвращать во время хода
	// при реализации контроллера сущности, чтобы вывести существо из боя
	ErrEntityEscaped = errors.New("entity escaped")
)

type FightResult int

const (
	FightIsNotOver FightResult = iota // бой еще не закончен
	FightWinLeft                      // победа первой группы существ
	FightWinRight                     // победа второй группы существ
	FightAllDead                      // все существа оказались мертвы после очередного хода
)

// Fight - структура, представляющая состояние активного боя
type Fight struct {
	left   []Fightable
	right  []Fightable
	order  *ring.Ring // очередность ходов среди всех сущностей (кольцевой буфер из Fightable)
	Result FightResult
}

type Fightable interface {
	IsAlive() bool
	MakeMove(context.Context) error
}

// NewFight - конструктор боя
func NewFight[T Fightable](left []T, right []T, order []T) (*Fight, error) {
	all := make([]Fightable, 0, len(left)+len(right))
	for _, e := range left {
		all = append(all, e)
	}
	for _, e := range right {
		all = append(all, e)
	}

	if !allAlive(all) {
		return nil, ErrEntityIsDead
	}

	if !allDifferent(all) {
		return nil, ErrEntityRepeated
	}

	// check that order is valid
	for _, e := range order {
		if !slices.Contains(all, Fightable(e)) {
			return nil, ErrOrderInvalid
		}
	}

	leftF := make([]Fightable, len(left))
	for i, e := range left {
		leftF[i] = e
	}

	rightF := make([]Fightable, len(right))
	for i, e := range right {
		rightF[i] = e
	}

	ring := ring.New(len(order))
	for _, e := range order {
		ring.Value = e
		ring = ring.Next()
	}

	return &Fight{
		left:   leftF,
		right:  rightF,
		order:  ring,
		Result: FightIsNotOver,
	}, nil
}

// NewFightRandomOrder - конструктор боя со случайным порядком ходов
func NewFightRandomOrder(src rand.Source, left []Fightable, right []Fightable) (*Fight, error) {
	order := make([]Fightable, 0, len(left)+len(right))
	order = append(order, left...)
	order = append(order, right...)

	rand.New(src).Shuffle(len(order), func(i, j int) {
		order[i], order[j] = order[j], order[i]
	})

	return NewFight(left, right, order)
}

// NewFightSemiRandomOrder - конструктор боя с сохранением относительного порядка команд
func NewFightSemiRandomOrder(src rand.Source, left []Fightable, right []Fightable) (*Fight, error) {
	order := make([]Fightable, 0, len(left)+len(right))
	rand := rand.New(src)

	i, j := 0, 0
	for i < len(left) || j < len(right) {
		if rand.Intn(2) == 0 && i < len(left) {
			order = append(order, left[i])
			i++
		} else if j < len(right) {
			order = append(order, right[j])
			j++
		} else {
			order = append(order, left[i])
			i++
		}
	}

	return NewFight(left, right, order)
}

func allAlive(entities []Fightable) bool {
	for _, e := range entities {
		if !e.IsAlive() {
			return false
		}
	}
	return true
}

func allDifferent(entities []Fightable) bool {
	for i := 0; i < len(entities)-1; i++ {
		for j := i + 1; j < len(entities); j++ {
			if entities[i] == entities[j] {
				return false
			}
		}
	}
	return true
}

// PerformFight - выполняет бой поочередно давая сущностям возможность ходить
func (f *Fight) PerformFight(ctx context.Context) error {
	if f.Result != FightIsNotOver {
		return ErrFightIsOver
	}

	var err error
	for err == nil && hasAlive(f.left) && hasAlive(f.right) && f.order.Len() > 1 {
		ent := f.order.Value.(Fightable)
		if ent.IsAlive() {
			err = ent.MakeMove(ctx)
			// проверить, сбежало ли существо
			if err == ErrEntityEscaped {
				f.handleEscape()
				err = nil
			}
		} else {
			// удаляем мертвое существо из буфера очередности
			f.handleEscape()
		}

		f.order = f.order.Next()
	}

	// определить результат боя
	if err == nil {
		if hasAlive(f.left) {
			f.Result = FightWinLeft
		} else if hasAlive(f.right) {
			f.Result = FightWinRight
		} else {
			f.Result = FightAllDead
		}
	}
	return err
}

func hasAlive(entities []Fightable) bool {
	for _, e := range entities {
		if e.IsAlive() {
			return true
		}
	}
	return false
}

// handleEscape - полностью удаляет текущее существо в кольцевом буфере из битвы.
// После этого кольцевой буфер будет указывать на предыдущее существо
func (f *Fight) handleEscape() {
	ent := f.order.Value

	for i, e := range f.left {
		if e == ent {
			f.left = append(f.left[:i], f.left[i+1:]...)
			break
		}
	}

	for i, e := range f.right {
		if e == ent {
			f.right = append(f.right[:i], f.right[i+1:]...)
			break
		}
	}

	// удалить существо из кольцевого буфера
	for node := f.order.Next(); node != f.order; node = node.Next() {
		if node.Value == ent {
			node = node.Prev()
			node.Unlink(1)
		}
	}

	// удалить последний (первый) элемент
	f.order = f.order.Prev()
	f.order.Unlink(1)
}

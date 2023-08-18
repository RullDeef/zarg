package domain

import (
	"context"
	"errors"
	"math/rand"
	"slices"
)

var (
	// ErrEntityIsDead - сущность мертва на начало боя
	ErrEntityIsDead = errors.New("entity is dead")

	// ErrEntityRepeated - сущность повторяется
	ErrEntityRepeated = errors.New("entity is repeated")

	ErrOrderInvalid = errors.New("order is invalid")

	ErrFightIsOver    = errors.New("fight is over")
	ErrFightIsNotOver = errors.New("fight is not over")
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
	left  []Fightable
	right []Fightable

	// order - очередность ходов среди всех сущностей
	order []Fightable

	Result FightResult
}

type Fightable interface {
	IsAlive() bool
	MakeMove(context.Context) error
}

// NewFight - конструктор боя
func NewFight(left []Fightable, right []Fightable, order []Fightable) (*Fight, error) {
	all := make([]Fightable, 0, len(left)+len(right))
	all = append(all, left...)
	all = append(all, right...)

	if !allAlive(all) {
		return nil, ErrEntityIsDead
	}

	if !allDifferent(all) {
		return nil, ErrEntityRepeated
	}

	// check that order is valid
	for _, e := range order {
		if !slices.Contains(all, e) {
			return nil, ErrOrderInvalid
		}
	}

	return &Fight{
		left:   left,
		right:  right,
		order:  order,
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
	for i := 0; err == nil && hasAlive(f.left) && hasAlive(f.right); i++ {
		i %= len(f.order)
		if ent := f.order[i]; ent.IsAlive() {
			err = ent.MakeMove(ctx)
		}
	}

	// determine fight result
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

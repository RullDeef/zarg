package domain

import (
	"errors"
)

var (
	ErrEffectExpired = errors.New("effect expired")         // эффект истек и не может быть применен
	ErrMergeWithSelf = errors.New("cannot merge with self") // эффект не может быть объединен с самим собой
)

// Effect - активный накладываемый на существо эффект.
// Существом может быть как игрок так и враг
type Effect struct {
	Name        string           // название эффекта
	Description string           // описание эффекта
	SourceType  EffectSourceType // тип источника эффекта.

	// Duration - длительность эффекта в оставшихся ходах данного существа.
	// Если эффект длится до конца боя, то Duration == -1
	Duration int

	// activator - функция активации эффекта.
	// Вызывается при передачи хода данному игроку
	activator func() error

	// merger - функция объединения эффектов данного типа
	merger func(*Effect, *Effect) (*Effect, error)
}

type EffectSourceType int

const (
	EffectFromOponent  EffectSourceType = iota // источник эффекта - опонент
	EffectFromTeammate                         // источник эффекта - союзник
	EffectFromSelf                             // источник эффекта - само существо (кратковременный баф)
	EffectFromCharm                            // источник эффекта - амулет у существа (на весь бой)
)

// Activate - функция, активирующая эффект.
// Возвращает true, если эффект был исчерпан и должет быть удален
func (e *Effect) Activate() (bool, error) {
	if !e.isLimited() {
		return false, e.activator()
	}

	if e.Duration > 0 {
		e.Duration--
		return e.Duration == 0, e.activator()
	}

	return false, ErrEffectExpired
}

func (e *Effect) isLimited() bool {
	return e.Duration != -1
}

func (e *Effect) MergeWith(other *Effect) (*Effect, error) {
	if e == other {
		return nil, ErrMergeWithSelf
	}
	return e.merger(e, other)
}

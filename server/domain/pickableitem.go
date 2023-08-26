package domain

import "context"

// PickableItem - структура, представляющая один подбираемый предмет
type PickableItem struct {
	Title       string  // краткое название предмета
	Description string  // полное описание предмета
	Weight      float64 // вес предмета

	// Cost - стоимость при продаже.
	// Если продажа невозможна, то стоимость равна 0
	Cost int

	IsStoryline       bool // является ли предмет сюжетным (например, инструмент плотника)
	IsWeapon          bool // является ли предмет оружием
	IsEquipable       bool // можно ли надеть предмет (броня, аксессуары, амулеты)
	KeepOnEscape      bool // остается ли предмет в инвентаре при побеге
	KeepOnDeath       bool // остается ли предмет после смерти
	KeepAfterCompaign bool // остается ли предмет после прохождения похода

	// UseCases - список вариантов использования предмета во время боя.
	// Если предмет не может быть использован, список пуст
	UseCases []*ItemUseCase
}

type ItemUseCase struct {
	Title            string // краткое описание варианта использования
	Description      string // подробное описание варианта использования
	CanBeUsedOnRest  bool   // может ли быть использован во время отдыха
	CanBeUsedInFight bool   // может ли быть использован во время боя
	UseSpendsMove    bool   // тратит ли использование ход в бою
	IsDestructive    bool   // предмет уничтожится после использования?

	// UsesLeft - оставшееся количество использований (если IsDestructive == false).
	// Если предмет может быть использован неограниченное число раз, UsesLeft == -1
	UsesLeft int

	Action func(context.Context, *PickableItem) error // само действие
}

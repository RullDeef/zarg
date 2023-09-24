package domain

import "context"

// PickableItem - структура, представляющая один подбираемый предмет
type PickableItem struct {
	Title       string   // краткое название предмета
	Description string   // полное описание предмета
	Kind        ItemKind // тип предмета (оружие, броня, зелье и тп)
	Weight      float64  // вес предмета

	// Cost - стоимость при продаже.
	// Если продажа невозможна, то стоимость равна 0
	Cost int

	IsStoryline       bool // является ли предмет сюжетным (например, инструмент плотника)
	IsWeapon          bool // является ли предмет оружием
	IsEquipable       bool // можно ли надеть предмет (броня, аксессуары, амулеты)
	KeepOnEscape      bool // остается ли предмет в инвентаре при побеге
	KeepOnDeath       bool // остается ли предмет после смерти
	KeepAfterCompaign bool // остается ли предмет после прохождения похода

	Rarity float64 // вероятность найти данный предмет в сокровищнице (от 0 до 1)

	// UseCases - список вариантов использования предмета во время боя.
	// Если предмет не может быть использован, список пуст
	UseCases []*ItemUseCase
}

// ItemUseCase - способ использования предмета
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

// ItemKind - тип предмета
type ItemKind int

const (
	ItemKindWeaponMelee           = iota // оружие / ближний бой
	ItemKindWeaponRanged                 // оружие / дальний бой
	ItemKindWeaponMagic                  // оружие / магическое
	ItemKindWeaponThrowable              // оружие / метательное
	ItemKindArmorCuirass                 // броня / кираса
	ItemKindArmorTassets                 // броня / налядвенники
	ItemKindArmorHelmet                  // броня / шлем
	ItemKindArmorGloves                  // броня / перчатки
	ItemKindArmorBoots                   // броня / ботинки
	ItemKindPotionHealing                // зелье / восстанавливающее
	ItemKindPotionEffect                 // зелье / накладывающее эффекты
	ItemKindPotionBoost                  // зелье / повышающее характеристики
	ItemKindSpecialCarpenterTool         // особое / инструмент плотника
	ItemKindSpecialJewel                 // особое / драгоценные камни
	ItemKindSpecialEnchantedStone        // особое / зачарованные камни
	ItemKindCharm                        // амулет
)

// IsWeapon - является ли предмет оружием
func (k ItemKind) IsWeapon() bool {
	return k == ItemKindWeaponMelee ||
		k == ItemKindWeaponRanged ||
		k == ItemKindWeaponMagic ||
		k == ItemKindWeaponThrowable
}

// IsArmor - является ли предмет броней
func (k ItemKind) IsArmor() bool {
	return k == ItemKindArmorCuirass ||
		k == ItemKindArmorTassets ||
		k == ItemKindArmorHelmet ||
		k == ItemKindArmorGloves ||
		k == ItemKindArmorBoots
}

// IsPotion - является ли предмет зельем
func (k ItemKind) IsPotion() bool {
	return k == ItemKindPotionHealing ||
		k == ItemKindPotionEffect ||
		k == ItemKindPotionBoost
}

// IsSpecial - является ли предмет особым
func (k ItemKind) IsSpecial() bool {
	return k == ItemKindSpecialCarpenterTool ||
		k == ItemKindSpecialJewel ||
		k == ItemKindSpecialEnchantedStone
}

// Clone - полностью дублирует предмет. Полезно для
// инстанцирования предметов при генерации лута или сокровищ.
func (pi *PickableItem) Clone() *PickableItem {
	return &PickableItem{
		Title:             pi.Title,
		Description:       pi.Description,
		Kind:              pi.Kind,
		Weight:            pi.Weight,
		Cost:              pi.Cost,
		IsStoryline:       pi.IsStoryline,
		IsWeapon:          pi.IsWeapon,
		IsEquipable:       pi.IsEquipable,
		KeepOnEscape:      pi.KeepOnEscape,
		KeepOnDeath:       pi.KeepOnDeath,
		KeepAfterCompaign: pi.KeepAfterCompaign,
		Rarity:            pi.Rarity,
		UseCases:          cloneUseCases(pi.UseCases),
	}
}

func cloneUseCases(useCases []*ItemUseCase) []*ItemUseCase {
	new := make([]*ItemUseCase, len(useCases))
	for i, uc := range useCases {
		new[i] = &ItemUseCase{
			Title:            uc.Title,
			Description:      uc.Description,
			CanBeUsedOnRest:  uc.CanBeUsedOnRest,
			CanBeUsedInFight: uc.CanBeUsedInFight,
			UseSpendsMove:    uc.UseSpendsMove,
			IsDestructive:    uc.IsDestructive,
			UsesLeft:         uc.UsesLeft,
			Action:           uc.Action,
		}
	}
	return new
}

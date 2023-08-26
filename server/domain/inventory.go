package domain

import "errors"

var (
	ErrOverweight   = errors.New("overweight")     // перегруз
	ErrItemNotFound = errors.New("item not found") // предмет не найден в инвентаре
)

type Inventory struct {
	Items []*PickableItem `json:"items"`

	// constraints - ограничения на комбинации предметов при попытке подбора.
	// Новый предмет будет успешно побран, если он удовлетворяет всем ограничениям
	constraints []InventoryConstraint
}

// InventoryConstraint - ограничение на подбираемые предметы.
// error должен описывать причину несоответствия условию.
//
// Типичный пример ограничения - ограничение по весу предметов. Новый предмет
// не может быть подобран, если суммарный вес инвентаря превзойдет некоторое значение
type InventoryConstraint interface {
	CanHoldIn(inv *Inventory, newItem *PickableItem) error
}

type InventoryConstraintFunc func(inv *Inventory, newItem *PickableItem) error

func (f InventoryConstraintFunc) CanHoldIn(inv *Inventory, newItem *PickableItem) error {
	return f(inv, newItem)
}

func NewEmptyInventory(constraints ...InventoryConstraint) *Inventory {
	return &Inventory{
		Items:       nil,
		constraints: constraints,
	}
}

// Weight - вычисляет суммарный вес инвентаря (всех предметов)
func (inv *Inventory) Weight() float64 {
	var weight float64
	for _, item := range inv.Items {
		weight += item.Weight
	}
	return weight
}

// Cost - вычисляет суммарную стоимость всех предметов инвентаря
func (inv *Inventory) Cost() int {
	var cost int
	for _, item := range inv.Items {
		cost += item.Cost
	}
	return cost
}

// HasItem - проверяет, есть ли данный предмет в инвентаре (по ссылке)
func (inv *Inventory) HasItem(item *PickableItem) bool {
	for _, i := range inv.Items {
		if i == item {
			return true
		}
	}
	return false
}

// Pickup - подбирает предмет, если он удовлетворяет всем ограничениям
func (inv *Inventory) Pickup(item *PickableItem) error {
	if err := inv.checkConstraints(item); err != nil {
		return err
	}

	inv.Items = append(inv.Items, item)
	return nil
}

func (inv *Inventory) checkConstraints(item *PickableItem) error {
	for _, constraint := range inv.constraints {
		if err := constraint.CanHoldIn(inv, item); err != nil {
			return err
		}
	}
	return nil
}

// Drop - удаляет предмет из инвентаря.
// Можно использовать при выбрасывании или расходовании предметов
func (inv *Inventory) Drop(item *PickableItem) error {
	for i, itemInInventory := range inv.Items {
		if itemInInventory == item {
			inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			return nil
		}
	}
	return ErrItemNotFound
}

// NewMaxWeightConstraint - ограничение по максимальному весу
func NewMaxWeightConstraint(maxWeight func() float64) InventoryConstraint {
	return InventoryConstraintFunc(func(inv *Inventory, newItem *PickableItem) error {
		if maxWeight() < inv.Weight()+newItem.Weight {
			return ErrOverweight
		}
		return nil
	})
}

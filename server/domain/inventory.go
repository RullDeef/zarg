package domain

import "errors"

var (
	ErrOverweight         = errors.New("overweight")          // перегруз
	ErrConstraintMismatch = errors.New("constraint mismatch") // невозможно взять предмет из-за ограничений
	ErrItemNotFound       = errors.New("item not found")      // предмет не найден в инвентаре
)

type Inventory struct {
	Items []*PickableItem `json:"items"`

	// constraints - ограничения на комбинации предметов при попытке подбора.
	// Новый предмет будет успешно побран, если он удовлетворяет всем ограничениям
	constraints []InventoryConstraint
}

type InventoryConstraint interface {
	CanHoldWith(items []*PickableItem, newItem *PickableItem) bool
}

type InventoryConstraintFunc func(items []*PickableItem, newItem *PickableItem) bool

func (f InventoryConstraintFunc) CanHoldWith(items []*PickableItem, newItem *PickableItem) bool {
	return f(items, newItem)
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

// Pickup - подбирает предмет, если для него хватает веса в инвентаре
func (inv *Inventory) Pickup(item *PickableItem, maxWeight float64) error {
	if inv.Weight()+item.Weight > maxWeight {
		return ErrOverweight
	}

	if !inv.checkConstraints(item) {
		return ErrConstraintMismatch
	}

	inv.Items = append(inv.Items, item)
	return nil
}

func (inv *Inventory) checkConstraints(item *PickableItem) bool {
	for _, constraint := range inv.constraints {
		if !constraint.CanHoldWith(inv.Items, item) {
			return false
		}
	}
	return true
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

package domain

// PickableItem - структура, представляющая один подбираемый предмет
type PickableItem struct {
	Title       string  // краткое название предмета
	Description string  // полное описание предмета
	Weight      float64 // вес предмета

	// Cost - стоимость при продаже.
	// Если продажа невозможна, то стоимость равна 0
	Cost int
}

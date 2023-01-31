package armor

type ArmorItem struct {
	name string
}

func New(name string) *ArmorItem {
	return &ArmorItem{
		name: name,
	}
}

// Pickup interface implementation
func (a *ArmorItem) Name() string {
	return a.name
}

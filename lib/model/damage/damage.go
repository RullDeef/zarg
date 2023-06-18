package damage

type BaseDamage struct {
	BaseDamageStats

	isCrit bool
}

func NewDamage(stats *BaseDamageStats, isCrit bool) *BaseDamage {
	return &BaseDamage{
		BaseDamageStats: *stats,
		isCrit:          isCrit,
	}
}

func (d *BaseDamage) IsCrit() bool {
	return d.isCrit
}

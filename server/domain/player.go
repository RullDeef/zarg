package domain

// Player - структура представляющая игрока во время похода
type Player struct {
	Entity           // базовая игровая сущность
	profile *Profile // обратная ссылка на профиль игрока
}

// NewPlayer - создает нового игрока по его профилю
func NewPlayer(profile *Profile) *Player {
	maxHealth := profile.MaxHealth()
	effectGroup := NewEffectGroup(NoopEffectPolicy())

	return &Player{
		Entity:  NewEntityBase(profile.Nickname, maxHealth, profile.Inventory, effectGroup),
		profile: profile,
	}
}

package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type ProfileID string

// Profile is in-game player representation.
// Account is a description of ways to identify a player (log in by any means).
type Profile struct {
	ID      ProfileID `json:"id"` // ID игрока
	Account *Account  `json:"-"`

	Nickname string `json:"nickname"` // никнейм игрока
	Avatar   string `json:"avatar"`   // url-ссылка на изображение
	Money    int    `json:"money"`    // количество монет игрока

	Strength    int `json:"strength"`    // сила. Влияет на максимальный переносимый вес
	Endurance   int `json:"endurance"`   // выносливость. Влияет на максимальное здоровье
	Luck        int `json:"luck"`        // удача. Влияет на шанс критического удара
	Observation int `json:"observation"` // внимательность. Влияет на шанс найти дополнительное сокровище

	// GuildID - ID гильдии, в которой находится игрок в настоящий момент.
	// Если игрок не состоит ни в одной гильдии - данное поле пусто
	GuildID GuildID `json:"guild_id"`

	Inventory *Inventory `json:"inventory"`
}

func NewAnonymousProfile() *Profile {
	ID := uuid.New().String()
	Nickname := fmt.Sprintf("Игрок%s", ID[:4])

	return &Profile{
		ID:          ProfileID(ID),
		Account:     nil,
		Nickname:    Nickname,
		Money:       0,
		Strength:    0,
		Endurance:   0,
		Luck:        0,
		Observation: 0,
		Inventory:   NewEmptyInventory(),
	}
}

// MaxHealth - вычисляет максимальное здоровье
func (p *Profile) MaxHealth() int {
	return 100 + 5*p.Endurance
}

// MaxWeight - вычисляет максимальный переносимый вес
func (p *Profile) MaxWeight() int {
	return 30 + 2*p.Strength
}

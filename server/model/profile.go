package model

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

type Profile struct {
	ID        string `json:"id"`
	AccountID string `json:"-"`

	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"` // url
	Money    int    `json:"money"`

	Strength    int `json:"strength"`
	Endurance   int `json:"endurance"`
	Luck        int `json:"luck"`
	Observation int `json:"observation"`

	Inventory *Inventory `json:"inventory"`
}

func NewAnonymousProfile() *Profile {
	ID := uuid.New().String()
	AccountID := ""
	Nickname := fmt.Sprintf("Игрок%d", 1000+rand.Intn(8999))

	return &Profile{
		ID:        ID,
		AccountID: AccountID,
		Nickname:  Nickname,
		Money:     0,

		Strength:    0,
		Endurance:   0,
		Luck:        0,
		Observation: 0,
		Inventory:   NewEmptyInventory(),
	}
}

func (p *Profile) MaxHealth() int {
	return 100 + 5*p.Endurance
}

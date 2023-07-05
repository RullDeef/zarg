package model

import (
	"container/list"

	smodel "server/model"
)

type Player struct {
	Health        int
	Profile       *smodel.Profile
	Inventory     *smodel.Inventory // same as Profile.Inventory
	StatusEffects *list.List
}

func NewPlayer(profile *smodel.Profile) *Player {
	return &Player{
		Health:        profile.MaxHealth(),
		Profile:       profile,
		Inventory:     profile.Inventory,
		StatusEffects: list.New(),
	}
}

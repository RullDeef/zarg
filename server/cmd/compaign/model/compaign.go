package model

import "time"

type Compaign struct {
	ID        string     `json:"id"`
	StartTime time.Time  `json:"-"`
	EndTime   *time.Time `json:"-"`
	Players   []*Player  `json:"players"`
	Dungeon   *Dungeon   `json:"dungeon"`
}

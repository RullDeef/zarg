package model

import "time"

type Dungeon struct {
	ID           string
	CreationTime time.Time
	Seed         int64
	Rooms        []Room
}

type Room interface {
}

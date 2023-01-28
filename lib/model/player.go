package model

const maxPlayerHealth = 100

type Player struct {
	userID int
	name   string
	Health int
	Weapon *Weapon
}

func NewPlayer(userID int, name string) *Player {
	return &Player{
		userID: userID,
		name:   name,
		Health: maxPlayerHealth,
		Weapon: nil,
	}
}

func (p *Player) UserID() int {
	return p.userID
}

func (p *Player) Name() string {
	return p.name
}

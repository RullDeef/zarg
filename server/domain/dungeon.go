package domain

import "context"

// Dungeon - интерфейс подземелья любого типа.
// Предоставляет возможность перемещаться между комнатами подземелья
type Dungeon interface {
	Entrance() Room             // получить первую комнату подземелья
	NextWays(Room) []DungeonWay // возможные направления для перемещения из данной комнаты
	IsFinalRoom(Room) bool      // является ли комната финальной
}

type Room interface {
	Visit(context.Context, *Compaign) error // посещение комнаты членами похода
}

// DungeonWay - направление для передвижения между комнатами
type DungeonWay struct {
	Name          string // текстовое представление направления
	VisitedBefore bool   // был ли данный путь уже посещен (для лабиринта)
	Room          Room
}

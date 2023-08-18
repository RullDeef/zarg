package mazegen

import (
	"math/rand"
	"server/domain"
)

// RegularDungeonGenerator - генератор линейного (обычного) подземелья из нескольких этажей
type RegularDungeonGenerator struct {
}

func (g *RegularDungeonGenerator) Generate(src rand.Source) (domain.Dungeon, error) {
	panic("not implemented")
}

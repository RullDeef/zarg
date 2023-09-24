package rooms

import (
	"context"
	"math/rand"
	"server/domain"
	"slices"
)

// TreasureRoom - комната с сокровищами
type TreasureRoom struct {
	items             []*domain.PickableItem
	distributor       domain.Distributor
	maxItemsPerPlayer int
}

type TreasureRoomGenerator struct {
	itemsPool         []*domain.PickableItem // пул предметов
	distributor       domain.Distributor     // распределитель предметов
	maxItemsPerPlayer int                    // максимальное количество предметов для одного игрока

	kindMap map[itemKind]uint // сколько предметов одного типа может быть в сокровищнице
}

// itemKind - тип предмета (обобщенный)
type itemKind int

const (
	itemKindWeapon itemKind = iota
	itemKindArmor
	itemKindPotion
	itemKindSpecial
)

// newTreasureRoom - создает новую комнату с сокровищами
func newTreasureRoom(
	items []*domain.PickableItem,
	distributor domain.Distributor,
	maxItemsPerPlayer int,
) *TreasureRoom {
	return &TreasureRoom{
		items:             items,
		distributor:       distributor,
		maxItemsPerPlayer: maxItemsPerPlayer,
	}
}

// Visit - запускает механизм распределения сокровищ между игроками
func (t *TreasureRoom) Visit(ctx context.Context, c *domain.Compaign) error {
	distr, err := t.distributor.Distribute(ctx, c.Participators, t.items, domain.DistributionConfig{
		Timeout:               0,
		MaxItemsPerPlayer:     t.maxItemsPerPlayer,
		MultipleOwnersAllowed: false,
	})
	if err != nil {
		return err
	}

	for p, items := range distr {
		for _, item := range items {
			err := p.Pickup(item)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// NewTreasureRoomGenerator - создает стандартный генератор комнаты сокровищ
func NewTreasureRoomGenerator(
	items []*domain.PickableItem,
	distributor domain.Distributor,
	participantsCount uint,
	maxItemsPerPlayer int,
) *TreasureRoomGenerator {
	kindMap := make(map[itemKind]uint) // TODO: формализовать и зафиксировать этот блок
	if participantsCount < 3 {
		kindMap[itemKindWeapon] = 2
		kindMap[itemKindArmor] = 2
		kindMap[itemKindPotion] = 3
		kindMap[itemKindSpecial] = 1
	} else if participantsCount < 5 {
		kindMap[itemKindWeapon] = 3
		kindMap[itemKindArmor] = 3
		kindMap[itemKindPotion] = 4
		kindMap[itemKindSpecial] = 1
	} else {
		kindMap[itemKindWeapon] = 4
		kindMap[itemKindArmor] = 4
		kindMap[itemKindPotion] = 5
		kindMap[itemKindSpecial] = 1
	}

	return &TreasureRoomGenerator{
		itemsPool:         items,
		distributor:       distributor,
		maxItemsPerPlayer: maxItemsPerPlayer,
		kindMap:           kindMap,
	}
}

func (g *TreasureRoomGenerator) Generate(src rand.Source) (domain.Room, error) {
	// choose items randomly from the pool
	maxItems := 0
	for _, count := range g.kindMap {
		maxItems += int(count)
	}

	rnd := rand.New(src)
	items := make([]*domain.PickableItem, 0)

	weapons := g.kindMap[itemKindWeapon]
	armors := g.kindMap[itemKindArmor]
	potions := g.kindMap[itemKindPotion]
	specials := g.kindMap[itemKindSpecial]

	// перебираем предметы случайным образом. Если новый предмет не помещается в сокровищницу,
	// но имеет меньшую редкость, то он может заменить предметы с большей редкостью (меньшим rarity).
	for i := 0; i < maxItems; i++ {
		perm := rnd.Perm(len(g.itemsPool))
		item := g.itemsPool[perm[0]].Clone()
		needReplace := false

		if item.Kind.IsWeapon() {
			if weapons > 0 {
				weapons--
			} else {
				needReplace = true
			}
		} else if item.Kind.IsArmor() {
			if armors > 0 {
				armors--
			} else {
				needReplace = true
			}
		} else if item.Kind.IsPotion() {
			if potions > 0 {
				potions--
			} else {
				needReplace = true
			}
		} else if item.Kind.IsSpecial() {
			if specials > 0 {
				specials--
			} else {
				needReplace = true
			}
		} else {
			i--
			continue // no need to sort
		}

		if needReplace {
			for j := 0; j < len(items); j++ {
				if items[j].Rarity < item.Rarity {
					items[j] = item
					break
				}
			}
		} else {
			items = append(items, item)
		}

		slices.SortFunc(items, func(i, j *domain.PickableItem) int {
			return int(10000*i.Rarity) - int(10000*j.Rarity) // TODO: made a better comparator
		})
	}

	return newTreasureRoom(items, g.distributor, g.maxItemsPerPlayer), nil
}

package rooms

import (
	"context"
	"server/domain"
)

// TreasureRoom - комната с сокровищами
type TreasureRoom struct {
	items             []*domain.PickableItem
	distributor       domain.Distributor
	maxItemsPerPlayer int
}

type TreasureRoomGenerator struct {
	itemsPool map[poolType][]itemDescriptor // пулы дескрипторов предметов
}

// poolType тип пула предметов
type poolType int

const (
	PoolTreasureRoom poolType = iota // в комнате сокровищ
	PoolEnemyDrop                    // дроп с врага
)

// itemDescriptor - дескриптор предмета
type itemDescriptor struct {
	rarity map[poolType]float64 // редкость - вероятность обнаружить данный предмет в данном пуле
}

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

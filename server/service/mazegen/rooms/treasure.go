package rooms

import (
	"context"
	"server/domain"
)

type TreasureRoom struct {
	items             []*domain.PickableItem
	distributor       domain.Distributor
	maxItemsPerPlayer int
}

// NewTreasureRoom - создает новую комнату с сокровищами
func NewTreasureRoom(
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

package distributor

import (
	"context"
	"log"
	"server/domain"
	"slices"
	"sync"
)

// RegularDistributor - стандартный распредилитель вещей
// между игроками, который учитывает выбор самих игроков.
type RegularDistributor struct {
	onPlayerCanChooseFunc PlayerCanChooseFunc
}

// PlayerCanChooseFunc - функция, которая вызывается для выбора предмета
// игроком. В канал wantPickup игрок отправляет предметы из списка items
// которые он хотел бы поднять, а в канал wantDrop игрок отправляет
// предметы которые передумал брать.
//
// При получении пустого списка предметов функция должна немедленно завершиться,
// что будет означать подтверждение выбора игрока.
type PlayerCanChooseFunc func(
	ctx context.Context,
	player *domain.Player,
	items <-chan []*domain.PickableItem,
	wantPickup chan<- *domain.PickableItem,
	wantDrop chan<- *domain.PickableItem,
) error

func NewRegular(onPlayerCanChoose PlayerCanChooseFunc) *RegularDistributor {
	return &RegularDistributor{
		onPlayerCanChooseFunc: onPlayerCanChoose,
	}
}

// Distribute - распределяет вещи между игроками
func (d *RegularDistributor) Distribute(
	ctx context.Context,
	participants []*domain.Player,
	items []*domain.PickableItem,
	config domain.DistributionConfig,
) (domain.PlayerItemDistribution, error) {
	if len(participants) == 0 {
		return nil, domain.ErrNoParticipants
	}
	if len(items) == 0 {
		return nil, domain.ErrNoItems
	}

	if config.Timeout == 0 {
		return d.unlimitedDistribute(ctx, participants, items, config)
	} else {
		ctx, cancel := context.WithTimeout(ctx, config.Timeout)
		defer cancel()
		return d.unlimitedDistribute(ctx, participants, items, config)
	}
}

// unlimitedDistribute - распределение завершится, когда все предметы будут взяты участниками
func (d *RegularDistributor) unlimitedDistribute(
	ctx context.Context,
	participants []*domain.Player,
	items []*domain.PickableItem,
	config domain.DistributionConfig,
) (domain.PlayerItemDistribution, error) {
	distrib := make(domain.PlayerItemDistribution, len(participants))
	freeItems := slices.Clone(items) // предметы, еще не подобранные ни одним участником
	lock := sync.Mutex{}
	wg := sync.WaitGroup{}

	itemsChans := make([]chan []*domain.PickableItem, len(participants))
	wantPickupChans := make([]chan *domain.PickableItem, len(participants))
	wantDropChans := make([]chan *domain.PickableItem, len(participants))
	for i, player := range participants {
		i, player := i, player
		itemsChans[i] = make(chan []*domain.PickableItem, 16)
		wantPickupChans[i] = make(chan *domain.PickableItem, 16)
		wantDropChans[i] = make(chan *domain.PickableItem, 16)

		// wants handler
		go func() {
			for {
				select {
				case item := <-wantPickupChans[i]:
					lock.Lock()
					// check that item is free
					if slices.Contains(freeItems, item) {
						j := slices.Index(freeItems, item)
						freeItems = slices.Delete(slices.Clone(freeItems), j, j+1)
						distrib[player] = append(distrib[player], item)
						for _, c := range itemsChans {
							c <- freeItems
						}
					}
					lock.Unlock()
				case item := <-wantDropChans[i]:
					lock.Lock()
					if slices.Contains(items, item) {
						j := slices.Index(distrib[player], item)
						distrib[player] = slices.Delete(distrib[player], j, j+1)
						freeItems = append(slices.Clone(freeItems), item)
						for _, c := range itemsChans {
							if c != nil {
								c <- freeItems
							}
						}
					}
					lock.Unlock()
				case <-ctx.Done():
					return
				}
			}
		}()

		go func() {
			defer wg.Done()
			defer close(wantPickupChans[i])
			defer close(wantDropChans[i])

			err := d.onPlayerCanChooseFunc(
				ctx,
				player,
				itemsChans[i],
				wantPickupChans[i],
				wantDropChans[i],
			)
			if err != nil {
				log.Println(err) // TODO: handle error gracefully
			}
		}()

		wg.Add(1)
	}

	// send starting free items
	lock.Lock()
	for i := range participants {
		itemsChans[i] <- freeItems
	}
	lock.Unlock()

	wg.Wait()
	for i := range participants {
		close(itemsChans[i])
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return distrib, nil
}

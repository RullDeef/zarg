package main

import (
	"context"
	"math/rand"
	"server/domain"
	"server/internal/modules/logger"
	"server/service/distributor"
	"server/service/mazegen"
	"time"
)

// Тестирование обычного похода в подземелье
func testRegularDungeon() {
	log, _ := logger.New()
	src := rand.NewSource(7431118648)

	players := []*domain.Profile{
		newTestProfile("Mike"),
		newTestProfile("John"),
	}

	generator, err := mazegen.NewRegularGenerator(
		2,    // participants count
		0.0,  // guild score avg
		0.0,  // dungeons completed avg
		true, // from one guild
		[]*domain.PickableItem{
			{
				Title:             "sword",
				Description:       "Melee weapon",
				Kind:              domain.ItemKindWeaponMelee,
				Weight:            5.0,
				Cost:              20,
				IsWeapon:          true,
				DropAfterCompaign: true,
				Rarity:            1.0,
			},
			{
				Title:             "wand",
				Description:       "Magic weapon",
				Kind:              domain.ItemKindWeaponMagic,
				Weight:            2.0,
				Cost:              40,
				IsWeapon:          true,
				DropAfterCompaign: true,
				Rarity:            0.6,
			},
		}, // items pool
		distributor.NewRegular(func(
			ctx context.Context,
			player *domain.Player,
			items <-chan []*domain.PickableItem,
			wantPickup, wantDrop chan<- *domain.PickableItem,
		) error {
			rnd := rand.New(src)
			for freeItems := range items {
				if len(freeItems) == 0 {
					break
				}

				// random delay && try pick random item
				time.Sleep(time.Duration(rnd.Intn(200)) * time.Millisecond)
				item := freeItems[rnd.Intn(len(freeItems))]
				wantPickup <- item
			}
			return ctx.Err()
		}), // distributor
		2, // max items per player
	)
	if err != nil {
		log.Sugar().Fatalw("failed to create dungeon generator", "error", err)
	}

	dungeon, err := generator.Generate(src)
	if err != nil {
		log.Sugar().Fatalw("failed to generate dungeon", "error", err)
	}

	compaign := domain.NewCompaign(players, domain.TrivialWayChooser)
	err = compaign.VisitDungeon(context.Background(), dungeon)
	if err != nil {
		log.Sugar().Fatalw("failed to visit dungeon", "error", err)
	}

	log.Info("done")
}

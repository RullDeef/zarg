// Sandbox - модуль песочницы для тестирования локальных гипотез и отловли багов.
// Код в песочнице не должен следовать принципам чистой архитектуры,
// его задача - тестирование гипотез.

package main

import (
	"context"
	"math/rand"
	"server/domain"
	"server/internal/modules/logger"
	"server/service/mazegen"

	"go.uber.org/zap"
)

func main() {
	testRegularDungeon()
}

// Тестирование обычного похода в подземелье
func testRegularDungeon() {
	log := logger.New()

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
				KeepAfterCompaign: true,
				Rarity:            1.0,
			},
			{
				Title:             "wand",
				Description:       "Magic weapon",
				Kind:              domain.ItemKindWeaponMagic,
				Weight:            2.0,
				Cost:              40,
				IsWeapon:          true,
				KeepAfterCompaign: true,
				Rarity:            0.6,
			},
		}, // items pool
		newTestDistributor(log), // distributor
		2,                       // max items per player
	)
	if err != nil {
		log.Sugar().Fatalw("failed to create dungeon generator", "error", err)
	}

	dungeon, err := generator.Generate(rand.NewSource(7431118648))
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

// создает тестового игрока с пустым инвентарем
func newTestProfile(name string) *domain.Profile {
	profile := domain.NewAnonymousProfile()
	profile.Nickname = name
	return profile
}

type testDistributor struct {
	logger *zap.Logger
}

var _ domain.Distributor = &testDistributor{}

func (d *testDistributor) Distribute(
	ctx context.Context,
	players []*domain.Player,
	items []*domain.PickableItem,
	config domain.DistributionConfig,
) (domain.PlayerItemDistribution, error) {
	d.logger.Sugar().Infoln("distributor: items:")
	for i, item := range items {
		d.logger.Sugar().Infof("%d) %s", i, item.Title)
	}
	d.logger.Sugar().Infoln("players:")
	for i, player := range players {
		d.logger.Sugar().Infof("%d) %s", i, player.Name)
	}
	d.logger.Sugar().Infoln("distributing evenly...")
	distr := make(domain.PlayerItemDistribution, len(players))
	j := 0
	for _, item := range items {
		player := players[j]
		d.logger.Sugar().Infof("give %s to %s", item.Title, player.Name)
		distr[player] = append(distr[player], item)
		j++
		if j >= len(players) {
			j = 0
		}
	}

	return distr, nil
}

func newTestDistributor(logger *zap.Logger) domain.Distributor {
	return &testDistributor{
		logger: logger,
	}
}

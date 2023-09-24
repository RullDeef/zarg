package main

import (
	"context"
	"server/domain"

	"go.uber.org/zap"
)

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

package main

import (
	"context"
	"server/internal/modules/importer"
	"server/internal/modules/logger"
)

func testImporter() {
	log, err := logger.New()
	if err != nil {
		panic(err)
	}
	ii := importer.NewItemImporter("data/items/test.yml")
	ctx := logger.WithLogger(context.Background(), log)
	item, err := ii.Import(ctx, 1)
	if err != nil {
		panic(err)
	}
	log.Sugar().Infow("item", "item", item)

	// TODO: UseCases не импортируется
}

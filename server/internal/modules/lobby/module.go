// Package lobby - пакет отвечающий за логику распределения участников по группам и походам
package lobby

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("lobby",
	fx.Provide(NewLobbyManaged),
	fx.Provide(NewController),
)

func NewLobbyManaged(lc fx.Lifecycle, logger *zap.SugaredLogger) *Lobby {
	lobby := NewLobby(logger)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return lobby.Close()
		},
	})

	return lobby
}

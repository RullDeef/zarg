package logger

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("logger",
	fx.Provide(New),
	fx.Provide(func(logger *zap.Logger) *zap.SugaredLogger {
		return logger.Sugar()
	}),
)

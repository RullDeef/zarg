package logger

import (
	"context"

	"go.uber.org/zap"
)

type loggerKey struct{}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	if ctx.Value(loggerKey{}) != nil {
		panic("tried to set logger in context with logger")
	}
	return context.WithValue(ctx, loggerKey{}, logger)
}

func GetLogger(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*zap.Logger)
	if !ok {
		panic("logger requested, but not found in context")
	}
	return logger
}

func GetSugaredLogger(ctx context.Context) *zap.SugaredLogger {
	return GetLogger(ctx).Sugar()
}

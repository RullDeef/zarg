package logger

import "go.uber.org/zap"

func New() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

func Sugarize(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}

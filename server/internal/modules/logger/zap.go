package logger

import "go.uber.org/zap"

func New() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger
}

func Sugarize(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}

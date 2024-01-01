package logger

import (
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func InitLogger() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()

	return &sugar
}

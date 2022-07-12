package logger

import (
	"go.uber.org/zap"
)

func GetLogger() *zap.SugaredLogger {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return zapLogger.Sugar()
}

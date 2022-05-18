package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var log *zap.Logger

func InitLogging() {
	var err error
	log, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defer log.Sync()
}

func Info(msg string) {
	log.Info(msg)
}

func Infof(template string, args ...interface{}) {
	log.Info(fmt.Sprintf(template, args...))
}

func Warn(msg string) {
	log.Warn(msg)
}

func Warnf(template string, args ...interface{}) {
	log.Warn(fmt.Sprintf(template, args...))
}

func Error(msg string) {
	log.Error(msg)
}

func Errorf(template string, args ...interface{}) {
	log.Error(fmt.Sprintf(template, args...))
}

func Fatal(msg string) {
	log.Fatal(msg)
}

func Fatalf(template string, args ...interface{}) {
	log.Fatal(fmt.Sprintf(template, args...))
}

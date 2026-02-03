package logger

import (
	"log"

	"go.uber.org/zap"
)

type Logger interface {
	Infow(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
}

var logger Logger

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	logger.Fatalw(msg, keysAndValues...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

func Init(debug bool) {
	var l *zap.Logger
	var err error

	if debug {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatalf("failed to start logger: %s", err)
	}

	logger = l.Sugar()
}

func Get() Logger {
	if logger == nil {
		log.Fatalf("init logger first")
	}

	return logger
}

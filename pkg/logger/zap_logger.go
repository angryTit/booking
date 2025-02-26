package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type zapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() Logger {
	logger, _ := zap.NewProduction()
	return &zapLogger{
		logger: logger,
	}
}

func (l *zapLogger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Info(fmt.Sprintf(msg, args...))
	} else {
		l.logger.Info(msg)
	}
}

func (l *zapLogger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Error(fmt.Sprintf(msg, args...))
	} else {
		l.logger.Error(msg)
	}
}

func (l *zapLogger) Fatal(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Fatal(fmt.Sprintf(msg, args...))
	} else {
		l.logger.Fatal(msg)
	}
}

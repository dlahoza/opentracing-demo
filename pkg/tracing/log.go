package tracing

import (
	"log"
)

type jaegerLoggerAdapter struct {
	logger *log.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Println(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Printf(msg, args...)
}

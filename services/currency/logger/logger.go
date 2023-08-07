package logger

import "log"

type LogLevel string

var defaultLogger Logger

type Logger interface {
	Log(level LogLevel, message string)
}

func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

func DefaultLog(level LogLevel, message string) {
	if defaultLogger != nil {
		defaultLogger.Log(level, message)
		return
	}

	log.Printf("[%s]: %s\n", string(level), message)
}

const (
	DEBUG LogLevel = "DEBUG"
	INFO           = "INFO"
	ERROR          = "ERROR"
)

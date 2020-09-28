package durable

import (
	"log"
)

type LoggerClient struct{}

type Logger struct{}

func BuildLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) LogLine(v ...interface{}) {
	log.Println(v...)
}

func (logger *Logger) LogFormat(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (logger *Logger) LogPanic(v ...interface{}) {
	log.Panicln(v...)
}

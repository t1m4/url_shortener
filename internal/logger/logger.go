package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"url_shortener/configs"
)

const (
	DEBUG = 2
	INFO  = 1
	ERROR = 0
)

type Logger interface {
	Debug(messages ...interface{})
	Info(messages ...interface{})
	Error(messages ...interface{})
}
type logger struct {
	level  int
	logger *log.Logger
}

func New(config *configs.Config) Logger {
	l := log.New(os.Stdout, "", log.LstdFlags|log.Llongfile)
	var logLevelMap = map[string]int{
		configs.DebugLevel: DEBUG,
		configs.InfoLevel:  INFO,
	}
	return &logger{
		level:  logLevelMap[config.LOG_LEVEL],
		logger: l,
	}
}
func (l *logger) Debug(messages ...interface{}) {
	if l.level >= DEBUG {
		l.logger.Println(configs.DebugLevel, concatMessages(messages...))
	}
}
func (l *logger) Info(messages ...interface{}) {
	if l.level >= INFO {
		l.logger.Println(configs.InfoLevel, concatMessages(messages...))
	}
}
func (l *logger) Error(messages ...interface{}) {
	if l.level >= ERROR {
		l.logger.Println(configs.ErrorLevel, concatMessages(messages...))
	}
}

func concatMessages(messages ...interface{}) string {
	var resultBuffer bytes.Buffer
	for _, message := range messages {
		resultBuffer.WriteString(fmt.Sprintf("%v ", message))
	}
	return resultBuffer.String()
}

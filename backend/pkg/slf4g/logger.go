package slf4g

import (
	"fmt"
	"strings"
)

const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type Logger interface {
	Fatal(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Trace(msg string, args ...interface{})

	Sync()
	Name() string
}

var lm LoggerManager

type LoggerManager struct {
	rootLogger Logger
	loggerMap  map[string]Logger
}

func init() {
	lm = LoggerManager{&ConsoleLogger{}, make(map[string]Logger)}
}

func R() Logger {
	return lm.rootLogger
}

func N(name string) Logger {
	name = strings.ToLower(name)
	nl := lm.loggerMap[name]
	if nl == nil {
		return lm.rootLogger
	}
	return nl
}

func Set(name string, logger Logger) {
	name = strings.ToLower(name)
	if name == "root" || name == "" {
		lm.rootLogger = logger
	}
	lm.loggerMap[name] = logger
}

func Sync() {
	for _, l := range lm.loggerMap {
		if l != nil {
			l.Sync()
		}
	}
}

type ConsoleLogger struct {
}

func (l *ConsoleLogger) Fatal(msg string, args ...interface{}) {
	println("[F]", fmt.Sprintf(msg, args...))
}

func (l *ConsoleLogger) Error(msg string, args ...interface{}) {
	println("[E]", fmt.Sprintf(msg, args...))
}

func (l *ConsoleLogger) Warn(msg string, args ...interface{}) {
	println("[W]", fmt.Sprintf(msg, args...))
}

func (l *ConsoleLogger) Info(msg string, args ...interface{}) {
	println("[I]", fmt.Sprintf(msg, args...))
}

func (l *ConsoleLogger) Debug(msg string, args ...interface{}) {
	println("[D]", fmt.Sprintf(msg, args...))
}

func (l *ConsoleLogger) Trace(msg string, args ...interface{}) {
	println("[T]", fmt.Sprintf(msg, args...))
}

func (l *ConsoleLogger) Name() string {
	return "root"
}

func (l *ConsoleLogger) Sync() {
}

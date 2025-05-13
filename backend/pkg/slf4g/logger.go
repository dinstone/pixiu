package slf4g

import (
	"context"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Logger struct {
	ctx context.Context
}

var logger *Logger
var oncefn sync.Once

func Setup(ctx context.Context) *Logger {
	if logger == nil {
		oncefn.Do(func() {
			logger = &Logger{ctx}
		})
	}
	return logger
}

func Get() *Logger {
	return logger
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	runtime.LogFatalf(l.ctx, msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	runtime.LogErrorf(l.ctx, msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	runtime.LogWarningf(l.ctx, msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	runtime.LogInfof(l.ctx, msg, args...)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	runtime.LogDebugf(l.ctx, msg, args...)
}

func (l *Logger) Trace(msg string, args ...interface{}) {
	runtime.LogTracef(l.ctx, msg, args...)
}

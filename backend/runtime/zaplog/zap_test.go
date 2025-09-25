package zaplog

import (
	"os"
	"path/filepath"
	"pixiu/backend/pkg/constant"
	"pixiu/backend/pkg/slf4g"
	"testing"
)

func TestZapLogger(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	acd := filepath.Join(wd, constant.AppCode)
	l := setupLogger("dev", acd)
	l.Info("hello world")
}

func setupLogger(bt string, chd string) slf4g.Logger {
	ld := filepath.Join(chd, "logs")
	if bt == "dev" {
		hc := HandlerConfig{
			Name:  "console-handler",
			Type:  "console",
			Level: "debug",
		}
		fc := HandlerConfig{
			Name:     "file-handler",
			Type:     "file",
			Level:    "debug",
			LogDir:   ld,
			MaxAge:   2,
			FileName: "app.log",
		}
		rl := LoggerConfig{
			Name:     "root",
			Level:    "debug",
			ShowLine: true,
			Handlers: []string{"console-handler", "file-handler"},
		}
		Setup(ZapConfig{Loggers: []LoggerConfig{rl}, Handlers: []HandlerConfig{hc, fc}})
	} else {
		fc := HandlerConfig{
			Name:     "file-handler",
			Type:     "file",
			Level:    "debug",
			LogDir:   ld,
			MaxAge:   2,
			FileName: "app.log",
		}
		rl := LoggerConfig{
			Name:     "root",
			Level:    "info",
			ShowLine: true,
			Handlers: []string{"file-handler"},
		}
		Setup(ZapConfig{Loggers: []LoggerConfig{rl}, Handlers: []HandlerConfig{fc}})
	}

	return slf4g.R()
}

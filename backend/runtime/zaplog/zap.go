package zaplog

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"pixiu/backend/pkg/slf4g"
	"pixiu/backend/pkg/utils"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapConfig struct {
	Loggers  []LoggerConfig  `mapstructure:"loggers" json:"loggers" yaml:"loggers"`
	Handlers []HandlerConfig `mapstructure:"handlers" json:"handlers" yaml:"handlers"`
}

type LoggerConfig struct {
	Name     string   `mapstructure:"name" json:"name" yaml:"name"`
	Level    string   `mapstructure:"level" json:"level" yaml:"level"` // 日志级别
	ShowLine bool     `mapstructure:"show_line" json:"show_line" yaml:"show_line"`
	Handlers []string `mapstructure:"handlers" json:"handlers" yaml:"handlers"`
}

type HandlerConfig struct {
	Name  string `mapstructure:"name" json:"name" yaml:"name"`    // 名称标识
	Type  string `mapstructure:"type" json:"type" yaml:"type"`    // 输出类型：file、console
	Level string `mapstructure:"level" json:"level" yaml:"level"` // 日志级别

	Format   string `mapstructure:"format" json:"format" yaml:"format"` // 输出格式
	LogDir   string `mapstructure:"log_dir" json:"log_dir" yaml:"log_dir"`
	FileName string `mapstructure:"file_name" json:"file_name" yaml:"file_name"`
	MaxAge   int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"` // 日志留存时间
}

type zapLogger struct {
	name   string
	level  int8
	logger *zap.Logger
}

func (z *zapLogger) Fatal(msg string, args ...interface{}) {
	if checkLevel(z.level, slf4g.LevelFatal) {
		z.logger.Fatal(fmt.Sprintf(msg, args...))
	}
}

func (z *zapLogger) Error(msg string, args ...interface{}) {
	if checkLevel(z.level, slf4g.LevelError) {
		z.logger.Error(fmt.Sprintf(msg, args...))
	}
}

func (z *zapLogger) Warn(msg string, args ...interface{}) {
	if checkLevel(z.level, slf4g.LevelWarn) {
		z.logger.Warn(fmt.Sprintf(msg, args...))
	}
}

func (z *zapLogger) Info(msg string, args ...interface{}) {
	if checkLevel(z.level, slf4g.LevelInfo) {
		z.logger.Info(fmt.Sprintf(msg, args...))
	}
}

func (z *zapLogger) Debug(msg string, args ...interface{}) {
	if checkLevel(z.level, slf4g.LevelDebug) {
		z.logger.Debug(fmt.Sprintf(msg, args...))
	}
}

func (z *zapLogger) Trace(msg string, args ...interface{}) {
	if checkLevel(z.level, slf4g.LevelTrace) {
		z.logger.Debug(fmt.Sprintf(msg, args...))
	}
}

func (z *zapLogger) Name() string {
	return z.name
}

func (z *zapLogger) Sync() {
	z.logger.Sync()
}

func checkLevel(r int8, e int8) bool {
	return e >= r
}

func Setup(zc ZapConfig) {
	// init cores
	coreMap := make(map[string]zapcore.Core)
	for _, hc := range zc.Handlers {
		coreMap[hc.Name] = newZapCore(hc)
	}

	// parse logger
	lgm := make(map[string]*zapLogger)
	for _, lc := range zc.Loggers {
		ncs := make([]zapcore.Core, 0, 7)
		// find cores
		for _, cr := range lc.Handlers {
			c := coreMap[cr]
			if c != nil {
				ncs = append(ncs, c)
			}
		}
		if len(ncs) > 0 {
			zl := zap.New(zapcore.NewTee(ncs...))
			if lc.ShowLine {
				zl.WithOptions(zap.AddCaller())
			}
			ll := convertLevel(lc.Level)
			lgm[lc.Name] = &zapLogger{lc.Name, ll, zl}
		}
	}

	// set logger
	if len(lgm) > 0 {
		for n, l := range lgm {
			slf4g.Set(n, l)
		}
	}
}

func convertLevel(level string) int8 {
	level = strings.ToLower(level)
	switch level {
	case "trace":
		return slf4g.LevelTrace
	case "debug":
		return slf4g.LevelDebug
	case "info":
		return slf4g.LevelInfo
	case "warn":
		return slf4g.LevelWarn
	case "error":
		return slf4g.LevelError
	case "fatal":
		return slf4g.LevelFatal
	default:
		return slf4g.LevelTrace
	}
}

func newZapCore(cc HandlerConfig) zapcore.Core {
	switch cc.Type {
	case "file":
		return getRollingFileCore(cc)
	case "console":
		return getConsoleCore(cc)
	default:
		return nil
	}
}

func getConsoleCore(cc HandlerConfig) zapcore.Core {
	var ze zapcore.Encoder
	if cc.Format == "json" {
		ze = zapcore.NewJSONEncoder(GetConsoleEncoderConfig(cc))
	} else {
		ze = zapcore.NewConsoleEncoder(GetConsoleEncoderConfig(cc))
	}

	ws := zapcore.AddSync(os.Stdout)
	lf := getLevelPriority(transportLevel(cc.Level))

	return zapcore.NewCore(ze, ws, lf)
}

func getRollingFileCore(cc HandlerConfig) zapcore.Core {
	ws, err := GetFileWriteSyncer(cc)
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return nil
	}

	var ze zapcore.Encoder
	if cc.Format == "json" {
		ze = zapcore.NewJSONEncoder(GetFileEncoderConfig(cc))
	} else {
		ze = zapcore.NewConsoleEncoder(GetFileEncoderConfig(cc))
	}
	lf := getLevelPriority(transportLevel(cc.Level))

	return zapcore.NewCore(ze, ws, lf)
}

// GetConsoleEncoderConfig 获取控制台日志的 zapcore.EncoderConfig
func GetConsoleEncoderConfig(cc HandlerConfig) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 控制台使用颜色编码器
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 控制台使用短路径
	}
}

// GetFileEncoderConfig 获取文件日志的 zapcore.EncoderConfig
func GetFileEncoderConfig(cc HandlerConfig) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 文件使用小写编码器
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 文件使用全路径
	}
}

// CustomTimeEncoder 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// transportLevel 根据字符串转化为 zapcore.Level
func transportLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

// GetLevelPriority 根据 zapcore.Level 获取 zap.LevelEnablerFunc
func getLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	switch level {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool { // 调试级别
			return level >= zap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool { // 日志级别
			return level >= zap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool { // 警告级别
			return level >= zap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool { // 错误级别
			return level >= zap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool { // dpanic级别
			return level >= zap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool { // panic级别
			return level >= zap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool { // 终止级别
			return level >= zap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool { // 调试级别
			return level >= zap.DebugLevel
		}
	}
}

// GetFileWriteSyncer 获取 zapcore.WriteSyncer
func GetFileWriteSyncer(cc HandlerConfig) (zapcore.WriteSyncer, error) {
	// 判断是否有Director文件夹
	if ok, _ := utils.PathExists(cc.LogDir); !ok {
		_ = os.Mkdir(cc.LogDir, os.ModePerm)
	}

	// 使用file-rotatelogs进行日志分割
	fileWriter, err := rotatelogs.New(
		path.Join(cc.LogDir, "%Y-%m-%d", cc.FileName),
		rotatelogs.WithClock(rotatelogs.Local),
		rotatelogs.WithMaxAge(time.Duration(cc.MaxAge)*24*time.Hour), // 日志留存时间
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	return zapcore.AddSync(fileWriter), err
}

package dao

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type SqliteConfig struct {
	Type         string `mapstructure:"type" json:"type" yaml:"type"`
	Prefix       string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // 全局表前缀，单独定义TableName则不生效
	Singular     bool   `mapstructure:"singular" json:"singular" yaml:"singular"`                   // 是否开启全局禁用复数，true表示开启
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
	LogMode      string `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`                   // 是否开启Gorm全局日志
	LogZap       bool   `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"`
	Dsn          string `mapstructure:"dsn" json:"dsn" yaml:"dsn"`
}

func NewGormDB(config *SqliteConfig) (*gorm.DB, error) {
	// 提取 DSN 中的文件路径部分（去掉参数）
	dbPath := config.Dsn
	if idx := strings.Index(dbPath, "?"); idx != -1 {
		dbPath = dbPath[:idx]
	}
	// 创建目录（如果不存在）
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	// 打开数据库
	db, err := gorm.Open(sqlite.Open(config.Dsn), gormConfig(config))
	if err != nil {
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		return db, nil
	}
}

func gormConfig(sc *SqliteConfig) *gorm.Config {
	cfg := &gorm.Config{
		// 命名策略
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   sc.Prefix,   // 表前缀，在表名前添加前缀，如添加用户模块的表前缀 user_
			SingularTable: sc.Singular, // 是否使用单数形式的表名，如果设置为 true，那么 User 模型会对应 users 表
		},

		DisableForeignKeyConstraintWhenMigrating: true,
	}

	// 默认日志
	defaultLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})
	switch sc.LogMode {
	case "silent", "Silent":
		cfg.Logger = defaultLogger.LogMode(logger.Silent)
	case "error", "Error":
		cfg.Logger = defaultLogger.LogMode(logger.Error)
	case "warn", "Warn":
		cfg.Logger = defaultLogger.LogMode(logger.Warn)
	case "info", "Info":
		cfg.Logger = defaultLogger.LogMode(logger.Info)
	default:
		cfg.Logger = defaultLogger.LogMode(logger.Info)
	}

	return cfg
}

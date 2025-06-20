package container

import (
	"context"
	"os"
	"path/filepath"
	"pixiu/backend/adapter/dao"
	"pixiu/backend/adapter/storage"
	"pixiu/backend/business/stock"
	"pixiu/backend/business/system"
	"pixiu/backend/business/uaac"
	"pixiu/backend/pkg/constant"
	"pixiu/backend/pkg/gormer"
	"pixiu/backend/pkg/slf4g"
	"pixiu/backend/pkg/utils"
	"time"

	"github.com/vrischmann/userdir"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gorm.io/gorm"
)

// App struct
type App struct {
	ctx context.Context

	gdb *gorm.DB

	svs map[string]interface{}

	Info system.AppInfo
}

// NewApp creates a new App application struct
func NewApp() *App {
	info := system.AppInfo{
		AppName:   constant.AppName,
		AppCode:   constant.AppCode,
		Version:   "1.1.0",
		Comments:  "A modern lightweight cross-platform desktop system.",
		Copyright: "Copyright © 2025 dinstone all rights reserved",
	}
	return &App{svs: make(map[string]interface{}), Info: info}
}

func (a *App) Service(name string) interface{} {
	return a.svs[name]
}

func (a *App) Context() context.Context {
	return a.ctx
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// setup logger
	logger := slf4g.Setup(ctx)

	gcf := loadConfig()
	logger.Info("sqlite db config: %+v", gcf)

	ei := runtime.Environment(ctx)
	if ei.BuildType != "dev" {
		gcf.Dsn = filepath.Join(userdir.GetConfigHome(), constant.AppCode, gcf.Dsn)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		gcf.Dsn = filepath.Join(wd, gcf.Dsn)
	}
	logger.Info("sqlite db dir is %s", gcf.Dsn)

	gdb, err := dao.NewGormDB(gcf)
	if err != nil {
		panic(err)
	}
	a.gdb = gdb

	// 检查表是否存在
	exists := gdb.Migrator().HasTable(&uaac.Account{})
	if !exists {
		err := gdb.AutoMigrate(
			uaac.Account{}, uaac.Profile{}, stock.StockInfo{}, stock.Investment{}, stock.Transaction{},
		)
		if err != nil {
			logger.Warn("注册数据库表失败: %v\n", err)
			panic(err)
		}

		// init admin user account
		pwd := utils.BcryptHash("admin@123")
		gdb.Save(&uaac.Account{Username: "admin", Password: pwd, Disabled: false})
		gdb.Save(&uaac.Profile{Username: "admin", NickName: "管理员", Avatar: "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"})

		logger.Info("数据库表创建成功")
	}

	gormer := gormer.NewGormer(gdb)
	uacs := uaac.NewUaacService(gormer, dao.NewUaacDao(gormer))
	ss := stock.NewStockService(gormer, dao.NewStockDao(gormer))
	a.svs["UaacService"] = uacs
	a.svs["StockService"] = ss

	pls := storage.NewLocalStorage(constant.AppCode, "preferences.yaml")
	pss := system.NewSystemService(pls)
	a.svs["SystemService"] = pss

	// start window event
	go loopWindowEvent(ctx)
}

// This is called just after the front-end dom has been completely rendered
func (a *App) DomReady(ctx context.Context) {
	runtime.WindowShow(ctx)
}

func (a *App) Shutdown(ctx context.Context) {
	// close db
	db, err := a.gdb.DB()
	if err != nil {
		db.Close()
	}
}

func loadConfig() *dao.SqliteConfig {
	return &dao.SqliteConfig{
		Dsn:          "stock.db?cache=shared&mode=rw",
		LogMode:      "info",
		LogZap:       false,
		MaxIdleConns: 1,
		MaxOpenConns: 5,
		Prefix:       "t_",
		Singular:     true,
		Type:         "sqlite3",
	}

}

func loopWindowEvent(ctx context.Context) {
	var fullscreen, maximised, minimised, normal bool
	var width, height int
	var dirty bool
	for {
		time.Sleep(300 * time.Millisecond)
		if ctx == nil {
			continue
		}

		dirty = false
		if f := runtime.WindowIsFullscreen(ctx); f != fullscreen {
			// full-screen switched
			fullscreen = f
			dirty = true
		}

		if w, h := runtime.WindowGetSize(ctx); w != width || h != height {
			// window size changed
			width, height = w, h
			dirty = true
		}

		if m := runtime.WindowIsMaximised(ctx); m != maximised {
			maximised = m
			dirty = true
		}

		if m := runtime.WindowIsMinimised(ctx); m != minimised {
			minimised = m
			dirty = true
		}

		if n := runtime.WindowIsNormal(ctx); n != normal {
			normal = n
			dirty = true
		}

		if dirty {
			runtime.EventsEmit(ctx, "window_changed", map[string]any{
				"fullscreen": fullscreen,
				"width":      width,
				"height":     height,
				"maximised":  maximised,
				"minimised":  minimised,
				"normal":     normal,
			})
		}
	}
}

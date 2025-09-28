package engine

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"pixiu/backend/adapter/dao"
	"pixiu/backend/adapter/ipc"
	"pixiu/backend/adapter/storage"
	"pixiu/backend/business/stock"
	"pixiu/backend/business/system"
	"pixiu/backend/business/uaac"
	"pixiu/backend/pkg/constant"
	"pixiu/backend/pkg/gormer"
	"pixiu/backend/pkg/slf4g"
	"pixiu/backend/pkg/utils"
	"pixiu/backend/runtime/zaplog"
	"strings"

	"github.com/vrischmann/userdir"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gorm.io/gorm"
)

// App Engine is a Container
type AppEngine struct {
	ctx context.Context

	env runtime.EnvironmentInfo

	acd string // app config directory

	gdb *gorm.DB

	lcs []ipc.LifeCycle

	ncmap map[string]interface{} // named component map

	appInfo system.AppInfo
}

// NewAppEngine creates a new AppEngine struct
func NewAppEngine() *AppEngine {
	info := system.AppInfo{
		AppName:   constant.AppName,
		AppCode:   constant.AppCode,
		Version:   "v1.2.0",
		Comments:  "A modern lightweight cross-platform desktop system.",
		Copyright: "Copyright © 2025 dinstone all rights reserved",
	}

	ae := &AppEngine{ncmap: make(map[string]interface{}), appInfo: info}

	uapi := ipc.NewUaacApi(ae)
	sapi := ipc.NewStockApi(ae)
	papi := ipc.NewSystemApi(ae)
	ae.lcs = append(ae.lcs, uapi, sapi, papi)

	return ae
}

func (a *AppEngine) AppInfo() *system.AppInfo {
	return &a.appInfo
}

func (a *AppEngine) GetComponent(name string) interface{} {
	return a.ncmap[name]
}

func (a *AppEngine) WailsContext() context.Context {
	return a.ctx
}

func (a *AppEngine) ConfigHome() string {
	return a.acd
}

func (a *AppEngine) BindAPI() []interface{} {
	binds := make([]interface{}, len(a.lcs))
	for i, v := range a.lcs {
		binds[i] = v
	}

	return binds
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *AppEngine) Startup(ctx context.Context) {
	a.ctx = ctx

	a.env = runtime.Environment(ctx)
	// setup config home
	if a.env.BuildType == "dev" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		a.acd = filepath.Join(wd, constant.AppCode)
	} else {
		a.acd = filepath.Join(userdir.GetConfigHome(), constant.AppCode)
	}
	// make directory
	if _, err := os.Stat(a.acd); os.IsNotExist(err) {
		if err := os.MkdirAll(a.acd, 0755); err != nil {
			panic(err)
		}
	}

	// setup logger
	logger := setupLogger(a.env.BuildType, a.acd)

	gcf := loadSqliteConfig(a.acd)
	logger.Info("sqlite db config: %+v", gcf)

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
	a.ncmap["UaacService"] = uacs
	a.ncmap["StockService"] = ss

	pls := storage.NewLocalStorage(a.acd, "preferences.yaml")
	pss := system.NewSystemService(pls)
	a.ncmap["SystemService"] = pss

	a.ncmap["AvatorStorage"] = storage.NewAvatorStorage(a.acd)

	for _, v := range a.lcs {
		v.Start()
	}

}

func setupLogger(bt string, chd string) slf4g.Logger {
	ld := filepath.Join(chd, "logs")
	if bt == "dev" {
		hc := zaplog.HandlerConfig{
			Name:  "console-handler",
			Type:  "console",
			Level: "debug",
		}
		fc := zaplog.HandlerConfig{
			Name:     "file-handler",
			Type:     "file",
			Level:    "debug",
			LogDir:   ld,
			MaxAge:   2,
			FileName: "app.log",
		}
		rl := zaplog.LoggerConfig{
			Name:     "root",
			Level:    "debug",
			ShowLine: true,
			Handlers: []string{"console-handler", "file-handler"},
		}
		zaplog.Setup(zaplog.ZapConfig{Loggers: []zaplog.LoggerConfig{rl}, Handlers: []zaplog.HandlerConfig{hc, fc}})
	} else {
		fc := zaplog.HandlerConfig{
			Name:     "file-handler",
			Type:     "file",
			Level:    "debug",
			LogDir:   ld,
			MaxAge:   2,
			FileName: "app.log",
		}
		rl := zaplog.LoggerConfig{
			Name:     "root",
			Level:    "info",
			ShowLine: true,
			Handlers: []string{"file-handler"},
		}
		zaplog.Setup(zaplog.ZapConfig{Loggers: []zaplog.LoggerConfig{rl}, Handlers: []zaplog.HandlerConfig{fc}})
	}

	return slf4g.R()
}

// This is called just after the front-end dom has been completely rendered
func (a *AppEngine) DomReady(ctx context.Context) {
	runtime.WindowShow(ctx)
}

func (a *AppEngine) Shutdown(ctx context.Context) {
	for _, v := range a.lcs {
		v.Close()
	}

	// close db
	db, err := a.gdb.DB()
	if err != nil {
		db.Close()
	}

	slf4g.Sync()
}

func loadSqliteConfig(acd string) *dao.SqliteConfig {
	dsn := filepath.Join(acd, "stock.db?cache=shared&mode=rw")
	return &dao.SqliteConfig{
		Dsn:          dsn,
		LogMode:      "info",
		LogZap:       false,
		MaxIdleConns: 1,
		MaxOpenConns: 5,
		Prefix:       "t_",
		Singular:     true,
		Type:         "sqlite3",
	}

}

func (a *AppEngine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// 拦截 /avatar/{userId}.{timestamp} 请求
	if len(path) > len("/avatar/") && path[:len("/avatar/")] == "/avatar/" {
		avatorFile := path[len("/avatar/"):]
		dotIndex := strings.Index(avatorFile, ".")
		if dotIndex != -1 {
			avatorFile = avatorFile[:dotIndex]
		}

		storage := a.GetComponent("AvatorStorage").(*storage.AvatorStorage)
		fileData, err := storage.LoadAvatorFile(avatorFile)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// 设置响应头（图片类型）
		res.Header().Set("Content-Type", "image/webp")
		res.Write(fileData)
		return
	}

	// 3. 其他路径返回 404
	http.NotFound(res, req)
}

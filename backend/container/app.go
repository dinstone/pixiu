package container

import (
	"context"
	"fmt"
	"net/http"
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
	"strconv"
	"strings"
	"time"

	"github.com/vrischmann/userdir"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gorm.io/gorm"
)

// App struct
type App struct {
	ctx context.Context

	env runtime.EnvironmentInfo

	acd string // app config directory

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

func (a *App) ConfigHome() string {
	return a.acd
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// setup logger
	logger := slf4g.Setup(ctx)

	// setup config home
	a.env = runtime.Environment(ctx)
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
	a.svs["UaacService"] = uacs
	a.svs["StockService"] = ss

	pls := storage.NewLocalStorage(a.acd, "preferences.yaml")
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

func (a *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// 拦截 /avatar/{userId}.{timestamp} 请求
	if len(path) > len("/avatar/") && path[:len("/avatar/")] == "/avatar/" {
		avatorFile := path[len("/avatar/"):]
		dotIndex := strings.Index(avatorFile, ".")
		if dotIndex != -1 {
			avatorFile = avatorFile[:dotIndex]
		}

		fileData, err := a.LoadAvatorFile(avatorFile)
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

func (a *App) LoadAvatorFile(avatorFile string) ([]byte, error) {
	localPath := filepath.Join(a.acd, "avatars", avatorFile)
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return nil, err
	}

	// 读取并返回文件内容
	return os.ReadFile(localPath)
}

func (a *App) SaveAvatorFile(filePath string, userId string) (string, error) {
	// 1. 创建本地上传目录（若不存在）
	uploadDir := filepath.Join(a.acd, "avatars")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return "", fmt.Errorf("创建上传目录失败: %v", err)
		}
	}

	// 2. 读取文件内容
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	// 3. 保存文件至本地目录（以用户ID命名）
	savePath := filepath.Join(uploadDir, userId)
	if err := os.WriteFile(savePath, fileData, 0644); err != nil {
		return "", fmt.Errorf("保存文件失败: %v", err)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	avatorPath := "/avatar/" + userId + "." + timestamp
	return avatorPath, nil
}

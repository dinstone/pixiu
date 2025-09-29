package main

import (
	"embed"
	"fmt"
	"pixiu/backend/runtime/engine"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create an instance of the appEngine structure
	appEngine := engine.NewAppEngine()

	// menu
	isMacOS := runtime.GOOS == "darwin"
	appMenu := menu.NewMenu()
	if isMacOS {
		appMenu.Append(menu.AppMenu())
		appMenu.Append(menu.EditMenu())
		appMenu.Append(menu.WindowMenu())
	}

	appInfo := appEngine.AppInfo()
	// Create application with options
	err := wails.Run(&options.App{
		Title:            appInfo.AppName,
		Width:            1024,
		Height:           768,
		MinWidth:         960,
		MinHeight:        640,
		Menu:             appMenu,
		LogLevel:         logger.INFO,
		WindowStartState: options.Maximised,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: appEngine.AvatorHandler(),
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		StartHidden:      true,
		OnStartup:        appEngine.Startup,
		OnDomReady:       appEngine.DomReady,
		OnShutdown:       appEngine.Shutdown,
		Bind:             appEngine.BindAPI(),
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			About: &mac.AboutInfo{
				Title:   fmt.Sprintf("%s %s", appInfo.AppName, appInfo.Version),
				Message: appInfo.Comments + "\n\n" + appInfo.Copyright,
				Icon:    icon,
			},
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			DisableFramelessWindowDecorations: false,
		},
		Linux: &linux.Options{
			ProgramName:         appInfo.AppName,
			Icon:                icon,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyOnDemand,
			WindowIsTranslucent: true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

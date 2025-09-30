package ipc

import (
	"os/exec"
	"pixiu/backend/adapter/container"
	"pixiu/backend/business/system"
	"pixiu/backend/pkg/slf4g"
	"runtime"
)

type SystemApi struct {
	ac container.Container
	ss *system.SystemService
}

func NewSystemApi(c container.Container) *SystemApi {
	return &SystemApi{ac: c}
}

func (sa *SystemApi) Start() {
	sa.ss = sa.ac.GetComponent("SystemService").(*system.SystemService)
}

func (sa *SystemApi) Close() {

}

func (sa *SystemApi) GetAppInfo() *Result {
	return Success(sa.ac.AppInfo())
}

func (sa *SystemApi) CheckForUpdate() *Result {
	u, err := sa.ss.GetLatestUpdate(sa.ac.AppInfo().Version)
	if err != nil {
		slf4g.R().Warn("check update failed, %s", err)
		return Failure(err)
	}
	return Success(u)
}

func (sa *SystemApi) GetPreferences() *Result {
	return Success(sa.ss.GetPreferences())
}

func (sa *SystemApi) SetPreferences(pf *system.Preferences) (ret Result) {
	if err := sa.ss.SetPreferences(pf); err != nil {
		ret.Code = 500
		ret.Mesg = err.Error()
	}
	return
}

func (sa *SystemApi) UpdatePreferences(values map[string]any) (ret Result) {
	if err := sa.ss.UpdatePreferences(values); err != nil {
		ret.Code = 500
		ret.Mesg = err.Error()
	}
	return
}

// 打开文件管理器并选中指定路径的文件夹
func (sa *SystemApi) OpenConfigFolder() error {
	var cmd *exec.Cmd
	path := sa.ac.ConfigHome()
	switch runtime.GOOS {
	case "windows":
		// Windows: explorer /select,"C:\path\to\folder"
		cmd = exec.Command("explorer", "/select,\""+path+"\"")
	case "darwin":
		// macOS: open -R /path/to/folder
		cmd = exec.Command("open", "-R", path)
	case "linux":
		// Linux: xdg-open /path/to/folder
		// 注意：xdg-open 只是打开文件夹，无法自动选中
		cmd = exec.Command("xdg-open", path)
	default:
		return nil
	}

	return cmd.Start() // 使用 Start() 而不是 Run()，避免阻塞
}

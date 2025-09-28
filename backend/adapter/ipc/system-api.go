package ipc

import (
	"pixiu/backend/business/system"
	"pixiu/backend/pkg/slf4g"
	"pixiu/backend/runtime/container"
)

type SystemApi struct {
	app *container.App
}

func NewSystemApi(app *container.App) *SystemApi {
	return &SystemApi{app}
}

func (sa *SystemApi) GetAppInfo() *Result {
	return Success(sa.app.Info)
}

func (sa *SystemApi) CheckForUpdate() *Result {
	ps := getSystemService(sa.app)
	u, err := ps.GetLatestUpdate(sa.app.Info.Version)
	if err != nil {
		slf4g.R().Warn("check update failed, %s", err)
		return Failure(err)
	}
	return Success(u)
}

func (p *SystemApi) GetPreferences() *Result {
	ps := getSystemService(p.app)
	return Success(ps.GetPreferences())
}

func (p *SystemApi) SetPreferences(pf *system.Preferences) (ret Result) {
	ps := getSystemService(p.app)
	if err := ps.SetPreferences(pf); err != nil {
		ret.Code = 500
		ret.Mesg = err.Error()
	}
	return
}

func (p *SystemApi) UpdatePreferences(values map[string]any) (ret Result) {
	ps := getSystemService(p.app)
	if err := ps.UpdatePreferences(values); err != nil {
		ret.Code = 500
		ret.Mesg = err.Error()
	}
	return
}

func getSystemService(app *container.App) *system.SystemService {
	return app.Service("SystemService").(*system.SystemService)
}

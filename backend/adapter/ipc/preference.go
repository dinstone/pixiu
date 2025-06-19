package ipc

import (
	"pixiu/backend/business/system"
	"pixiu/backend/container"
)

type PreferenceApi struct {
	app *container.App
}

func NewPreferenceApi(app *container.App) *PreferenceApi {
	return &PreferenceApi{app}
}

func (p *PreferenceApi) GetPreferences() (ret Result) {
	ps := getPreferenceService(p.app)
	ret.Data = ps.GetPreferences()
	return
}

func (p *PreferenceApi) SetPreferences(pf *system.Preferences) (ret Result) {
	ps := getPreferenceService(p.app)
	if err := ps.SetPreferences(pf); err != nil {
		ret.Code = 500
		ret.Mesg = err.Error()
	}
	return
}

func (p *PreferenceApi) UpdatePreferences(values map[string]any) (ret Result) {
	ps := getPreferenceService(p.app)
	if err := ps.UpdatePreferences(values); err != nil {
		ret.Code = 500
		ret.Mesg = err.Error()
	}
	return
}

func getPreferenceService(app *container.App) *system.PreferenceService {
	return app.Service("PreferenceService").(*system.PreferenceService)
}

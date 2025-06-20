package system

type Preferences struct {
	Theme Theme `json:"theme" yaml:"theme"`
}

type Theme struct {
	Color  string `json:"color" yaml:"color"`
	Dark   bool   `json:"dark" yaml:"dark"`
	Font   string `json:"font" yaml:"font,omitempty"`
	Layout string `json:"layout" yaml:"layout"`
}

type AppInfo struct {
	AppName   string `json:"appName" yaml:"appName"`
	AppCode   string `json:"appCode" yaml:"appCode"`
	Version   string `json:"version" yaml:"version"`
	Comments  string `json:"comments" yaml:"comments"`
	Copyright string `json:"copyright" yaml:"copyright"`
}

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

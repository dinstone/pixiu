package uaac

type Account struct {
	Id       int64  `json:"id"`
	Disabled bool   `json:"disabled"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Profile struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	NickName string `json:"nickName"`
	Avatar   string `json:"avatar"`
}

type UserDetail struct {
	Account *Account `json:"account"`
	Profile *Profile `json:"profile"`
}

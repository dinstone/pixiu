package account

type Account struct {
	Id       int64
	Username string
	Password string
	// Profile  *Profile `orm:"rel(one)"`
}

type Profile struct {
	Id      int
	Gender  string
	Age     int
	Address string
	Email   string
}

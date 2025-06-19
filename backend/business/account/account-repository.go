package account

type AccountRepository interface {
	All() (*[]Account, error)
	Get(id int64) (*Account, error)
	Save(user *Account) (int64, error)
	Update(user *Account) (int64, error)
	Delete(id int64) (int64, error)
}

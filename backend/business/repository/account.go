package repository

import "pixiu/backend/business/model"

type AccountRepository interface {
	All() (*[]model.Account, error)
	Get(id int64) (*model.Account, error)
	Save(user *model.Account) (int64, error)
	Update(user *model.Account) (int64, error)
	Delete(id int64) (int64, error)
}

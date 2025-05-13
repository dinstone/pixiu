package dao

import (
	"errors"
	"pixiu/backend/pkg/exception"

	"gorm.io/gorm"
)

func WrapGormError(err error) exception.AppError {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.WrapBusiness(404, err.Error(), err)
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		return exception.WrapBusiness(403, err.Error(), err)
	}
	return exception.WrapService(600, err.Error(), err)
}

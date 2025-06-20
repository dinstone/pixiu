package dao

import (
	"context"
	"pixiu/backend/business/uaac"
	"pixiu/backend/pkg/gormer"
)

type UaacDao struct {
	ormer gormer.GormDS
}

func NewUaacDao(ormer gormer.GormDS) *UaacDao {
	return &UaacDao{
		ormer: ormer,
	}
}

func (d *UaacDao) FindAccount(ctx context.Context, username string) (*uaac.Account, error) {
	var account uaac.Account
	err := d.ormer.GDB(ctx).Where("username = ?", username).First(&account).Error
	return &account, WrapGormError(err)
}

func (d *UaacDao) FindProfile(ctx context.Context, username string) (*uaac.Profile, error) {
	var profile uaac.Profile
	err := d.ormer.GDB(ctx).Where("username = ?", username).First(&profile).Error
	return &profile, WrapGormError(err)
}

func (d *UaacDao) UpdatePassword(ctx context.Context, username string, password string) error {
	return WrapGormError(d.ormer.GDB(ctx).Model(&uaac.Account{}).Where("username = ?", username).Update("password", password).Error)
}

func (d *UaacDao) UpdateProfile(ctx context.Context, profile *uaac.Profile) error {
	return WrapGormError(d.ormer.GDB(ctx).Model(&uaac.Profile{}).Where("username = ?", profile.Username).Updates(profile).Error)
}

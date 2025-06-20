package uaac

import (
	"pixiu/backend/pkg/exception"
	"pixiu/backend/pkg/gormer"
	"pixiu/backend/pkg/utils"
)

type UaacService struct {
	gtm gormer.GormTM
	ucr UaacRepository
}

func NewUaacService(gtm gormer.GormTM, ucr UaacRepository) *UaacService {
	return &UaacService{
		gtm: gtm,
		ucr: ucr,
	}
}

func (us *UaacService) GetUserDetail(username string) (*UserDetail, error) {
	account, err := us.ucr.FindAccount(us.gtm.Context(), username)
	if err != nil {
		return nil, err
	}
	profile, err := us.ucr.FindProfile(us.gtm.Context(), username)
	if err != nil {
		return nil, err
	}
	return &UserDetail{
		Account: account,
		Profile: profile,
	}, nil
}

func (us *UaacService) UpdateProfile(profile *Profile) error {
	return us.ucr.UpdateProfile(us.gtm.Context(), profile)
}

func (us *UaacService) UpdatePassword(username string, password string) error {
	if username == "" || password == "" {
		return exception.NewBusiness(404, "用户名和密码不能为空")
	}

	return us.ucr.UpdatePassword(us.gtm.Context(), username, utils.BcryptHash(password))
}

func (us *UaacService) AuthenAccessToken(username string, token string) error {
	if username == "" || token == "" {
		return exception.NewBusiness(404, "用户名和访问令牌不能为空")
	}
	jwt := utils.NewJWT()
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return exception.WrapBusiness(403, "访问令牌无效", err)
	}
	if claims.Username != username {
		return exception.NewBusiness(403, "访问令牌无效")
	}
	return nil
}

func (us *UaacService) AuthenPassword(username string, password string) (*Account, error) {
	if username == "" || password == "" {
		return nil, exception.NewBusiness(404, "用户名和密码不能为空")
	}

	account, err := us.ucr.FindAccount(us.gtm.Context(), username)
	if err != nil {
		return nil, err
	}

	if account.Disabled {
		return nil, exception.NewBusiness(403, "账号已被禁用")
	}

	if !utils.BcryptCheck(password, account.Password) {
		return nil, exception.NewBusiness(403, "用户名或密码错误")
	}

	return account, nil
}

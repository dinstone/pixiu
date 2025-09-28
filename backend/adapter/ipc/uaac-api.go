package ipc

import (
	"pixiu/backend/adapter/container"
	"pixiu/backend/adapter/storage"
	"pixiu/backend/business/uaac"
	"pixiu/backend/pkg/exception"
	"pixiu/backend/pkg/utils"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type UaacApi struct {
	ac container.Container
	us *uaac.UaacService
}

func NewUaacApi(c container.Container) *UaacApi {
	return &UaacApi{
		ac: c,
	}
}

func (u *UaacApi) Start() {
	u.us = u.ac.GetComponent("UaacService").(*uaac.UaacService)
}

func (u *UaacApi) Close() {

}

func (u *UaacApi) GetUserDetail(token string) *Result {
	claims, err := checkToken(token)
	if err != nil {
		return Failure(err)
	}

	ud, err := u.us.GetUserDetail(claims.Username)
	if err != nil {
		return Failure(err)
	}

	return Success(ud)
}

func (u *UaacApi) AuthenPassword(username string, password string) *Result {
	account, err := u.us.AuthenPassword(username, password)
	if err != nil {
		return Failure(err)
	}

	token, err := utils.GenerateToken(account.Id, username)
	if err != nil {
		return Failure(err)
	}

	return Success(token)
}

func (u *UaacApi) AuthenAccessToken(username string, token string) *Result {
	claims, err := checkToken(token)
	if err != nil {
		return Failure(err)
	}

	if claims.Username != username {
		return Failure(exception.NewBusiness(403, "访问令牌无效"))
	}

	newToken, err := utils.GenerateToken(claims.Identity, username)
	if err != nil {
		return Failure(err)
	}
	return Success(newToken)
}

func (u *UaacApi) UpdatePassword(token string, password string) *Result {
	claims, err := checkToken(token)
	if err != nil {
		return Failure(err)
	}

	err = u.us.UpdatePassword(claims.Username, password)
	if err != nil {
		return Failure(err)
	}

	return Success(nil)
}

func (u *UaacApi) UpdateAvator(token string) *Result {
	claims, err := checkToken(token)
	if err != nil {
		return Failure(err)
	}

	// select file
	imgFile, err := runtime.OpenFileDialog(u.ac.WailsContext(), runtime.OpenDialogOptions{
		Title: "选择图片",
		Filters: []runtime.FileFilter{{
			DisplayName: "Images (*.png;*.jpg;*.jpeg;*.gif)",
			Pattern:     "*.png;*.jpg;*.jpeg;*.gif",
		}},
	})
	if err != nil {
		return Failure(err)
	}

	// upload file
	if imgFile != "" {
		// save file
		userId := claims.Username
		as := u.ac.GetComponent("AvatorStorage").(*storage.AvatorStorage)
		aurl, err := as.SaveAvatorFile(imgFile, userId)
		if err != nil {
			return Failure(err)
		}

		// update db
		profile := &uaac.Profile{
			Username: userId,
			Avatar:   aurl,
		}
		err = u.us.UpdateProfile(profile)
		if err != nil {
			return Failure(err)
		}

		return Success("OK")
	}

	return Success("NG")
}

func (u *UaacApi) UpdateProfile(token string, profile *uaac.Profile) *Result {
	claims, err := checkToken(token)
	if err != nil {
		return Failure(err)
	}

	profile.Username = claims.Username
	err = u.us.UpdateProfile(profile)
	if err != nil {
		return Failure(err)
	}
	return Success(nil)
}

func checkToken(token string) (*utils.CustomClaims, error) {
	if token == "" {
		return nil, exception.NewBusiness(403, "访问令牌为空")
	}

	jwt := utils.NewJWT()
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return nil, exception.WrapBusiness(403, "访问令牌无效", err)
	}

	return claims, nil
}

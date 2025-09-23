package ipc

import (
	"pixiu/backend/business/uaac"
	"pixiu/backend/pkg/exception"
	"pixiu/backend/pkg/utils"
	"pixiu/backend/runtime/container"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type UaacApi struct {
	app *container.App
}

func NewUaacApi(app *container.App) *UaacApi {
	return &UaacApi{
		app: app,
	}
}

func (u *UaacApi) GetUserDetail(token string) *Result {
	claims, err := checkToken(token)
	if err != nil {
		return Failure(err)
	}

	ud, err := getUaacService(u.app).GetUserDetail(claims.Username)
	if err != nil {
		return Failure(err)
	}

	return Success(ud)
}

func (u *UaacApi) AuthenPassword(username string, password string) *Result {
	service := getUaacService(u.app)
	account, err := service.AuthenPassword(username, password)
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

	service := getUaacService(u.app)
	err = service.UpdatePassword(claims.Username, password)
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
	imgFile, err := runtime.OpenFileDialog(u.app.Context(), runtime.OpenDialogOptions{
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
		aurl, err := u.app.SaveAvatorFile(imgFile, userId)
		if err != nil {
			return Failure(err)
		}

		// update db
		profile := &uaac.Profile{
			Username: userId,
			Avatar:   aurl,
		}
		service := getUaacService(u.app)
		err = service.UpdateProfile(profile)
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

	service := getUaacService(u.app)
	profile.Username = claims.Username
	err = service.UpdateProfile(profile)
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

func getUaacService(app *container.App) *uaac.UaacService {
	return app.Service("UaacService").(*uaac.UaacService)
}

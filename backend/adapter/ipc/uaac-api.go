package ipc

import (
	"fmt"
	"os"
	"path/filepath"
	"pixiu/backend/business/uaac"
	"pixiu/backend/container"
	"pixiu/backend/pkg/constant"
	"pixiu/backend/pkg/exception"
	"pixiu/backend/pkg/utils"
	"strconv"
	"time"

	"github.com/vrischmann/userdir"
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

	if imgFile != "" {
		// upload file
		userId := claims.Username
		err := saveFile(imgFile, userId)
		if err != nil {
			return Failure(err)
		}

		// update db
		err = saveInfo(userId, u)
		if err != nil {
			return Failure(err)
		}
		return Success("OK")
	}

	return Success("NG")
}

func saveInfo(userId string, u *UaacApi) error {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	avatorPath := "/avatar/" + userId + "." + timestamp
	profile := &uaac.Profile{
		Username: userId,
		Avatar:   avatorPath,
	}
	service := getUaacService(u.app)
	return service.UpdateProfile(profile)
}

func saveFile(filePath string, userId string) error {
	// 1. 创建本地上传目录（若不存在）
	uploadDir := filepath.Join(userdir.GetConfigHome(), constant.AppCode, "avatars")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return fmt.Errorf("创建上传目录失败: %v", err)
		}
	}

	// 2. 读取文件内容
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	// 3. 保存文件至本地目录（以用户ID命名）
	savePath := filepath.Join(uploadDir, userId)
	if err := os.WriteFile(savePath, fileData, 0644); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	return nil
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

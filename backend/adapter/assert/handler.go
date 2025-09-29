package assert

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"pixiu/backend/pkg/slf4g"
	"strconv"
	"strings"
	"time"
)

type AvatorHandler struct {
	uploadDir string
}

func (ah *AvatorHandler) Startup(configHome string) error {
	// 1. 创建本地上传目录（若不存在）
	uploadDir := filepath.Join(configHome, "avatars")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return fmt.Errorf("创建上传目录失败: %v", err)
		}
	}

	slf4g.R().Info("头像上传目录：%s", uploadDir)

	ah.uploadDir = uploadDir

	return nil
}

func (ah *AvatorHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// 拦截 /avatar/{userId}.{timestamp} 请求
	if len(path) > len("/avatar/") && path[:len("/avatar/")] == "/avatar/" {
		avatorFile := path[len("/avatar/"):]
		fileData, err := ah.LoadAvatorFile(avatorFile)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// 设置响应头（图片类型）
		res.Header().Set("Content-Type", "image/webp")
		res.Write(fileData)
		return
	}

	// 3. 其他路径返回 404
	http.NotFound(res, req)
}

func (a *AvatorHandler) LoadAvatorFile(avatorFile string) ([]byte, error) {
	// 构建本地文件读取路径
	dotIndex := strings.Index(avatorFile, ".")
	if dotIndex != -1 {
		avatorFile = avatorFile[:dotIndex]
	}
	localPath := filepath.Join(a.uploadDir, avatorFile)
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return nil, err
	}

	// 读取并返回文件内容
	return os.ReadFile(localPath)
}

func (a *AvatorHandler) SaveAvatorFile(filePath string, userId string) (string, error) {
	// 读取文件内容
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	// 保存文件至本地目录（以用户ID命名）
	savePath := filepath.Join(a.uploadDir, userId)
	if err := os.WriteFile(savePath, fileData, 0644); err != nil {
		return "", fmt.Errorf("保存文件失败: %v", err)
	}

	// 构建访问路径
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	avatorPath := "/avatar/" + userId + "." + timestamp
	return avatorPath, nil
}

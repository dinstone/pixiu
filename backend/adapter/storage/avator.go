package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type AvatorStorage struct {
	acd string
}

func NewAvatorStorage(appdir string) *AvatorStorage {
	return &AvatorStorage{
		acd: appdir,
	}
}

func (a *AvatorStorage) LoadAvatorFile(avatorFile string) ([]byte, error) {
	localPath := filepath.Join(a.acd, "avatars", avatorFile)
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return nil, err
	}

	// 读取并返回文件内容
	return os.ReadFile(localPath)
}

func (a *AvatorStorage) SaveAvatorFile(filePath string, userId string) (string, error) {
	// 1. 创建本地上传目录（若不存在）
	uploadDir := filepath.Join(a.acd, "avatars")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return "", fmt.Errorf("创建上传目录失败: %v", err)
		}
	}

	// 2. 读取文件内容
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	// 3. 保存文件至本地目录（以用户ID命名）
	savePath := filepath.Join(uploadDir, userId)
	if err := os.WriteFile(savePath, fileData, 0644); err != nil {
		return "", fmt.Errorf("保存文件失败: %v", err)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	avatorPath := "/avatar/" + userId + "." + timestamp
	return avatorPath, nil
}

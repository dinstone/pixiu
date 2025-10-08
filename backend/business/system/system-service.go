package system

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pixiu/backend/pkg/slf4g"
	"reflect"
	"strings"
	"sync"
	"time"
)

type SystemService struct {
	storage PreferenceRepository
	mutex   sync.Mutex
}

func NewSystemService(pr PreferenceRepository) *SystemService {
	return &SystemService{
		storage: pr,
	}
}

func (p *SystemService) GetLatestUpdate(currentVersion string) (*UpdateInfo, error) {
	url := "https://api.github.com/repos/dinstone/pixiu/releases/latest"

	// 创建HTTP客户端并设置超时时间
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 发起GET请求
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch update info: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch update info, status code: %d", resp.StatusCode)
	}

	// 解析响应体中的JSON数据
	var release struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if currentVersion < release.TagName {
		// 构造返回结果
		updateInfo := UpdateInfo{
			HtmlUrl: release.HTMLURL,
			Latest:  release.TagName,
			Current: currentVersion,
		}

		return &updateInfo, nil
	} else {
		return nil, fmt.Errorf("no update available")
	}
}

func (p *SystemService) getPreferences() (ret Preferences) {
	b, err := p.storage.Load()
	if err != nil {
		return
	}

	json.Unmarshal(b, &ret)
	return
}

// GetPreferences Get preferences from local
func (p *SystemService) GetPreferences() (ret Preferences) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.getPreferences()
}

func (p *SystemService) setPreferences(pf *Preferences, key string, value any) error {
	parts := strings.Split(key, ".")
	if len(parts) > 0 {
		var reflectValue reflect.Value
		if reflect.TypeOf(pf).Kind() == reflect.Ptr {
			reflectValue = reflect.ValueOf(pf).Elem()
		} else {
			reflectValue = reflect.ValueOf(pf)
		}
		for i, part := range parts {
			part = strings.ToUpper(part[:1]) + part[1:]
			reflectValue = reflectValue.FieldByName(part)
			if reflectValue.IsValid() {
				if i == len(parts)-1 {
					reflectValue.Set(reflect.ValueOf(value))
					return nil
				}
			} else {
				break
			}
		}
	}

	return fmt.Errorf("invalid key path(%s)", key)
}

func (p *SystemService) savePreferences(pf *Preferences) error {
	b, err := json.Marshal(pf)
	if err != nil {
		return err
	}

	if err = p.storage.Store(b); err != nil {
		return err
	}
	return nil
}

// SetPreferences replace preferences
func (p *SystemService) SetPreferences(pf *Preferences) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.savePreferences(pf)
}

// UpdatePreferences update values by key paths, the key path use "." to indicate multiple level
func (p *SystemService) UpdatePreferences(values map[string]any) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	pf := p.getPreferences()
	for path, v := range values {
		if err := p.setPreferences(&pf, path, v); err != nil {
			return err
		}
	}
	slf4g.R().Info("after save %v", pf)

	return p.savePreferences(&pf)
}

package service

import (
	"fmt"
	"pixiu/backend/business/model"
	"pixiu/backend/business/repository"
	"pixiu/backend/pkg/slf4g"
	"reflect"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type PreferenceService struct {
	storage repository.PreferenceRepository
	mutex   sync.Mutex
}

func NewPreferenceService(pr repository.PreferenceRepository) *PreferenceService {
	// storage := NewLocalStore("preferences.yaml")
	storage := pr
	return &PreferenceService{
		storage: storage,
	}
}

func (p *PreferenceService) getPreferences() (ret model.Preferences) {
	b, err := p.storage.Load()
	if err != nil {
		return
	}

	yaml.Unmarshal(b, &ret)
	return
}

// GetPreferences Get preferences from local
func (p *PreferenceService) GetPreferences() (ret model.Preferences) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.getPreferences()
}

func (p *PreferenceService) setPreferences(pf *model.Preferences, key string, value any) error {
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

func (p *PreferenceService) savePreferences(pf *model.Preferences) error {
	b, err := yaml.Marshal(pf)
	if err != nil {
		return err
	}

	if err = p.storage.Store(b); err != nil {
		return err
	}
	return nil
}

// SetPreferences replace preferences
func (p *PreferenceService) SetPreferences(pf *model.Preferences) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.savePreferences(pf)
}

// UpdatePreferences update values by key paths, the key path use "." to indicate multiple level
func (p *PreferenceService) UpdatePreferences(values map[string]any) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	pf := p.getPreferences()
	for path, v := range values {
		if err := p.setPreferences(&pf, path, v); err != nil {
			return err
		}
	}
	slf4g.Get().Info("after save %v", pf)

	return p.savePreferences(&pf)
}

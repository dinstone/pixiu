package container

import (
	"context"
	"pixiu/backend/business/system"
)

type Container interface {
	AppInfo() *system.AppInfo
	ConfigHome() string
	WailsContext() context.Context
	GetComponent(name string) interface{}
}

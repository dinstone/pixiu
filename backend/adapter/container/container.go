package container

import (
	"context"
	"pixiu/backend/adapter/assert"
	"pixiu/backend/business/system"
)

type Container interface {
	AppInfo() *system.AppInfo
	ConfigHome() string
	WailsContext() context.Context
	AvatorHandler() *assert.AvatorHandler
	GetComponent(name string) interface{}
}

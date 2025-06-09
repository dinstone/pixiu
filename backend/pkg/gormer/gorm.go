package gormer

import (
	"context"

	"gorm.io/gorm"
)

type GormDS interface {
	GDB(ctx context.Context) *gorm.DB
}

type GormTM interface {
	Context() context.Context
	Execute(fc func(ctx context.Context) error) error
}

type GormID struct{}

type Gormer struct {
	gdb *gorm.DB
}

func NewGormer(gdb *gorm.DB) *Gormer {
	return &Gormer{gdb: gdb}
}

func (g *Gormer) GDB(ctx context.Context) *gorm.DB {
	if ctx != nil {
		gdb := ctx.Value(GormID{}).(*gorm.DB)
		if gdb != nil {
			return gdb
		}
	}
	return g.gdb
}

// GormDB 实现GormTM接口
func (g *Gormer) Context() context.Context {
	return context.Background()
}

func (g *Gormer) Execute(fc func(ctx context.Context) error) error {
	return g.gdb.Transaction(func(tx *gorm.DB) error {
		ctx := context.WithValue(g.Context(), GormID{}, tx)
		return fc(ctx)
	})
}

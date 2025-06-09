package repository

import (
	"context"
	"pixiu/backend/business/model"
)

type StockRepository interface {
	GetStock(ctx context.Context, code string) (*model.StockInfo, error)

	SaveStock(ctx context.Context, si *model.StockInfo) error
	UpdateStock(ctx context.Context, si *model.StockInfo) error
	DeleteStock(ctx context.Context, code string) error

	AliveStocks(ctx context.Context) (*[]model.StockInfo, error)

	GetHolding(ctx context.Context, code string) (*model.Investment, error)

	CreateInvestment(ctx context.Context, invest *model.Investment) error
	UpdateInvestment(ctx context.Context, invest *model.Investment) error
	GetInvestment(ctx context.Context, id int64) (*model.Investment, error)
	DeleteInvestment(ctx context.Context, id int64) error

	CreateTransaction(ctx context.Context, trans *model.Transaction) error
	UpdateTransaction(ctx context.Context, trans *model.Transaction) error
	GetTransaction(ctx context.Context, id int64) (*model.Transaction, error)
	DeleteTransaction(ctx context.Context, id int64) error
	GetTransactions(ctx context.Context, investId int64) (*[]model.Transaction, error)
}

package stock

import (
	"context"
)

type StockRepository interface {
	GetStock(ctx context.Context, code string) (*StockInfo, error)

	SaveStock(ctx context.Context, si *StockInfo) error
	UpdateStock(ctx context.Context, si *StockInfo) error
	DeleteStock(ctx context.Context, code string) error

	AliveStocks(ctx context.Context) (*[]StockInfo, error)

	GetHolding(ctx context.Context, code string) (*Investment, error)

	CreateInvestment(ctx context.Context, invest *Investment) error
	UpdateInvestment(ctx context.Context, invest *Investment) error
	GetInvestment(ctx context.Context, id int64) (*Investment, error)
	DeleteInvestment(ctx context.Context, id int64) error

	CreateTransaction(ctx context.Context, trans *Transaction) error
	UpdateTransaction(ctx context.Context, trans *Transaction) error
	GetTransaction(ctx context.Context, id int64) (*Transaction, error)
	DeleteTransaction(ctx context.Context, id int64) error
	GetTransactions(ctx context.Context, investId int64) (*[]Transaction, error)

	GetClearList(context context.Context, stime string, ftime string) (*[]ClearStats, error)
	GetClearInvest(context context.Context, stockCode string, startTime string, finishTime string) (*[]Investment, error)
}

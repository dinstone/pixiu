package repository

import "pixiu/backend/business/model"

type StockRepository interface {
	GetStock(code string) (*model.StockInfo, error)

	SaveStock(si *model.StockInfo) error
	UpdateStock(si *model.StockInfo) error
	DeleteStock(code string) error

	AliveStocks() (*[]model.StockInfo, error)

	GetHolding(code string) (*model.Investment, error)

	CreateInvestment(invest *model.Investment) error
	UpdateInvestment(invest *model.Investment) error
	GetInvestment(id int64) (*model.Investment, error)
	DeleteInvestment(id int64) error

	CreateTransaction(trans *model.Transaction) error
	UpdateTransaction(trans *model.Transaction) error
	GetTransaction(id int64) (*model.Transaction, error)
	DeleteTransaction(id int64) error
	GetTransactions(investId int64) (*[]model.Transaction, error)
}

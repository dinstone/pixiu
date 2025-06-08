package ipc

import (
	"pixiu/backend/business/model"
	"pixiu/backend/business/service"
	"pixiu/backend/container"
	"pixiu/backend/pkg/slf4g"
)

type StockApi struct {
	app *container.App
}

func NewStockApi(app *container.App) *StockApi {
	return &StockApi{
		app: app,
	}
}

func (s *StockApi) GetStockList() (result Result) {
	ss := getStockService(s.app)
	sis, err := ss.GetStockList()
	if err != nil {
		handleError(err, &result)
		return
	}
	result.Data = sis
	return
}

func (s *StockApi) GetStock(code string) (result Result) {
	ss := getStockService(s.app)
	si, err := ss.GetStock(code)
	if err != nil {
		handleError(err, &result)
		return
	}
	result.Data = si
	return
}

func (s *StockApi) AddStock(si *model.StockInfo) (result Result) {
	ss := getStockService(s.app)
	err := ss.SaveStock(si)
	if err != nil {
		handleError(err, &result)
		return
	}
	return
}

func (s *StockApi) UpdateStock(si *model.StockInfo) (result Result) {
	ss := getStockService(s.app)
	err := ss.UpdateStock(si)
	if err != nil {
		handleError(err, &result)
		return
	}
	return
}

func (s *StockApi) DeleteStock(code string) (result Result) {
	ss := getStockService(s.app)
	err := ss.DeleteStock(code)
	if err != nil {
		handleError(err, &result)
		return
	}
	return
}

func (s *StockApi) GetHolding(code string) (result Result) {
	ss := getStockService(s.app)
	invest, err := ss.GetHolding(code)
	if err != nil {
		handleError(err, &result)
		return
	}
	result.Data = invest
	return
}

func (s *StockApi) AddTransaction(tran *model.Transaction) (result Result) {
	ss := getStockService(s.app)
	err := ss.AddTransaction(tran)
	if err != nil {
		slf4g.Get().Warn("add transaction error: %v", err)
		handleError(err, &result)
		return
	}
	return
}

func (s *StockApi) UpdateTransaction(tran *model.Transaction) (result Result) {
	ss := getStockService(s.app)
	err := ss.UpdateTransaction(tran)
	if err != nil {
		handleError(err, &result)
		return
	}
	return
}

func (s *StockApi) DeleteTransaction(tranId int64) (result Result) {
	ss := getStockService(s.app)
	err := ss.DeleteTransaction(tranId)
	if err != nil {
		handleError(err, &result)
		return
	}
	return
}

func (s *StockApi) GetTransactions(investId int64) (result Result) {
	ss := getStockService(s.app)
	tranArray, err := ss.GetTransactions(investId)
	if err != nil {
		handleError(err, &result)
		return
	}
	result.Data = tranArray
	return
}

func getStockService(app *container.App) *service.StockService {
	return app.Service("StockService").(*service.StockService)
}

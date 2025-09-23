package ipc

import (
	"pixiu/backend/business/stock"
	"pixiu/backend/runtime/container"
)

type StockApi struct {
	app *container.App
}

type ClearQuery struct {
	StartTime  string `json:"startTime"`
	FinishTime string `json:"finishTime"`
}

func NewStockApi(app *container.App) *StockApi {
	return &StockApi{
		app: app,
	}
}

func (s *StockApi) GetStockList() *Result {
	ss := getStockService(s.app)
	sis, err := ss.GetStockList()
	if err != nil {
		return Failure(err)
	}
	return Success(sis)
}

func (s *StockApi) GetStock(code string) *Result {
	ss := getStockService(s.app)
	si, err := ss.GetStock(code)
	if err != nil {
		return Failure(err)
	}
	return Success(si)
}

func (s *StockApi) AddStock(si *stock.StockInfo) *Result {
	ss := getStockService(s.app)
	err := ss.SaveStock(si)
	if err != nil {
		return Failure(err)
	}
	return Success(true)
}

func (s *StockApi) UpdateStock(si *stock.StockInfo) *Result {
	ss := getStockService(s.app)
	err := ss.UpdateStock(si)
	if err != nil {
		return Failure(err)
	}
	return Success(true)
}

func (s *StockApi) DeleteStock(code string) *Result {
	ss := getStockService(s.app)
	err := ss.DeleteStock(code)
	if err != nil {
		return Failure(err)
	}
	return Success(true)
}

func (s *StockApi) GetHolding(code string) *Result {
	ss := getStockService(s.app)
	invest, err := ss.GetHolding(code)
	if err != nil {
		return Failure(err)
	}
	return Success(invest)
}

func (s *StockApi) AddTransaction(tran *stock.Transaction) *Result {
	ss := getStockService(s.app)
	err := ss.AddTransaction(tran)
	if err != nil {
		return Failure(err)
	}
	return Success(true)
}

func (s *StockApi) UpdateTransaction(tran *stock.Transaction) *Result {
	ss := getStockService(s.app)
	err := ss.UpdateTransaction(tran)
	if err != nil {
		return Failure(err)
	}
	return Success(true)
}

func (s *StockApi) DeleteTransaction(tranId int64) *Result {
	ss := getStockService(s.app)
	err := ss.DeleteTransaction(tranId)
	if err != nil {
		return Failure(err)
	}
	return Success(true)
}

func (s *StockApi) GetTransactions(investId int64) *Result {
	ss := getStockService(s.app)
	tranArray, err := ss.GetTransactions(investId)
	if err != nil {
		return Failure(err)
	}
	return Success(tranArray)
}

func (s *StockApi) GetClearList(cq ClearQuery) *Result {
	ss := getStockService(s.app)
	clearList, err := ss.GetClearList(cq.StartTime, cq.FinishTime)
	if err != nil {
		return Failure(err)
	}
	return Success(clearList)
}

func (s *StockApi) GetStockClear(stockCode string, startTime string, finishTime string) *Result {
	ss := getStockService(s.app)
	cstats, err := ss.GetStockClear(stockCode, startTime, finishTime)
	if err != nil {
		return Failure(err)
	}
	return Success(cstats)
}

func getStockService(app *container.App) *stock.StockService {
	return app.Service("StockService").(*stock.StockService)
}

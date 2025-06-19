package ipc

import (
	"pixiu/backend/business/stock"
	"pixiu/backend/container"
	"pixiu/backend/pkg/slf4g"
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

func (s *StockApi) AddStock(si *stock.StockInfo) (result Result) {
	ss := getStockService(s.app)
	err := ss.SaveStock(si)
	if err != nil {
		handleError(err, &result)
		return
	}
	return
}

func (s *StockApi) UpdateStock(si *stock.StockInfo) (result Result) {
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

func (s *StockApi) AddTransaction(tran *stock.Transaction) (result Result) {
	ss := getStockService(s.app)
	err := ss.AddTransaction(tran)
	if err != nil {
		slf4g.Get().Warn("add transaction error: %v", err)
		handleError(err, &result)
		return
	}
	return
}

func (s *StockApi) UpdateTransaction(tran *stock.Transaction) (result Result) {
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

func (s *StockApi) GetClearList(cq ClearQuery) (result Result) {
	ss := getStockService(s.app)
	clearList, err := ss.GetClearList(cq.StartTime, cq.FinishTime)
	if err != nil {
		handleError(err, &result)
		return
	}
	result.Data = clearList
	return
}

func (s *StockApi) GetStockClear(stockCode string, startTime string, finishTime string) (result Result) {
	ss := getStockService(s.app)
	cstats, err := ss.GetStockClear(stockCode, startTime, finishTime)
	if err != nil {
		handleError(err, &result)
		return
	}
	result.Data = cstats
	return
}

func getStockService(app *container.App) *stock.StockService {
	return app.Service("StockService").(*stock.StockService)
}

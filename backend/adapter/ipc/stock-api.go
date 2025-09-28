package ipc

import (
	"context"
	"pixiu/backend/adapter/container"
	"pixiu/backend/business/stock"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type StockApi struct {
	ac     container.Container
	ss     *stock.StockService
	cancel context.CancelFunc
}

type ClearQuery struct {
	StartTime  string `json:"startTime"`
	FinishTime string `json:"finishTime"`
}

func NewStockApi(ac container.Container) *StockApi {
	return &StockApi{
		ac: ac,
	}
}

func (s *StockApi) Start() {
	s.ss = s.ac.GetComponent("StockService").(*stock.StockService)

	var cctx context.Context
	cctx, s.cancel = context.WithCancel(context.Background())

	go loopWindowEvent(s.ac.WailsContext(), cctx)
}

func (s *StockApi) Close() {
	if s.cancel != nil {
		s.cancel()
	}
}

func loopWindowEvent(wctx context.Context, cctx context.Context) {
	var fullscreen, maximised, minimised, normal bool
	var width, height int
	var dirty bool
	for {
		select {
		case <-cctx.Done():
			return // 应用上下文已取消，退出循环
		default:
			time.Sleep(300 * time.Millisecond)
			if wctx == nil {
				continue
			}

			dirty = false
			if f := runtime.WindowIsFullscreen(wctx); f != fullscreen {
				// full-screen switched
				fullscreen = f
				dirty = true
			}

			if w, h := runtime.WindowGetSize(wctx); w != width || h != height {
				// window size changed
				width, height = w, h
				dirty = true
			}

			if m := runtime.WindowIsMaximised(wctx); m != maximised {
				maximised = m
				dirty = true
			}

			if m := runtime.WindowIsMinimised(wctx); m != minimised {
				minimised = m
				dirty = true
			}

			if n := runtime.WindowIsNormal(wctx); n != normal {
				normal = n
				dirty = true
			}

			if dirty {
				runtime.EventsEmit(wctx, "window_changed", map[string]any{
					"fullscreen": fullscreen,
					"width":      width,
					"height":     height,
					"maximised":  maximised,
					"minimised":  minimised,
					"normal":     normal,
				})
			}
		}
	}
}

func (s *StockApi) GetStockList() *Result {
	sis, err := s.ss.GetStockList()
	if err != nil {
		return Failure(err)
	}
	return Success(sis)
}

func (s *StockApi) GetStock(code string) *Result {
	si, err := s.ss.GetStock(code)
	if err != nil {
		return Failure(err)
	}
	return Success(si)
}

func (s *StockApi) AddStock(si *stock.StockInfo) *Result {
	err := s.ss.SaveStock(si)
	if err != nil {
		return Failure(err)
	}
	return Success(true)
}

func (s *StockApi) UpdateStock(si *stock.StockInfo) *Result {
	err := s.ss.UpdateStock(si)
	if err != nil {
		return Failure(err)
	}
	return Success(true)
}

func (s *StockApi) DeleteStock(code string) *Result {
	err := s.ss.DeleteStock(code)
	if err != nil {
		return Failure(err)
	}
	return Success(true)
}

func (s *StockApi) GetHolding(code string) *Result {
	invest, err := s.ss.GetHolding(code)
	if err != nil {
		return Failure(err)
	}
	return Success(invest)
}

func (s *StockApi) AddTransaction(tran *stock.Transaction) *Result {
	err := s.ss.AddTransaction(tran)
	if err != nil {
		return Failure(err)
	}
	return Success(true)
}

func (s *StockApi) UpdateTransaction(tran *stock.Transaction) *Result {
	err := s.ss.UpdateTransaction(tran)
	if err != nil {
		return Failure(err)
	}
	return Success(true)
}

func (s *StockApi) DeleteTransaction(tranId int64) *Result {
	err := s.ss.DeleteTransaction(tranId)
	if err != nil {
		return Failure(err)
	}
	return Success(true)
}

func (s *StockApi) GetTransactions(investId int64) *Result {
	tranArray, err := s.ss.GetTransactions(investId)
	if err != nil {
		return Failure(err)
	}
	return Success(tranArray)
}

func (s *StockApi) GetClearList(cq ClearQuery) *Result {
	clearList, err := s.ss.GetClearList(cq.StartTime, cq.FinishTime)
	if err != nil {
		return Failure(err)
	}
	return Success(clearList)
}

func (s *StockApi) GetStockClear(stockCode string, startTime string, finishTime string) *Result {
	cstats, err := s.ss.GetStockClear(stockCode, startTime, finishTime)
	if err != nil {
		return Failure(err)
	}
	return Success(cstats)
}

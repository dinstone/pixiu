package service

import (
	"pixiu/backend/business/model"
	"pixiu/backend/business/repository"
	"pixiu/backend/pkg/exception"
	"time"

	"github.com/shopspring/decimal"
)

const DateTimeLayout = "2006-01-02 15:04:05"

type StockService struct {
	sr repository.StockRepository
}

func NewStockService(sr repository.StockRepository) *StockService {
	return &StockService{sr: sr}
}

func (ss StockService) GetStockList() (*[]model.StockInfo, error) {
	return ss.sr.AliveStocks()
}

func (ss StockService) GetStock(code string) (*model.StockInfo, error) {
	if code == "" {
		return nil, exception.NewBusiness(400, "code is required")
	}
	return ss.sr.GetStock(code)
}

func (ss StockService) SaveStock(si *model.StockInfo) error {
	if si.Code == "" {
		return exception.NewBusiness(400, "code is required")
	}
	if si.Name == "" {
		return exception.NewBusiness(400, "name is required")
	}
	if si.Currency == "" {
		si.Currency = "人民币"
	}

	si.Status = 0
	si.CreatedAt = time.Now()
	si.UpdatedAt = time.Now()
	return ss.sr.SaveStock(si)
}

func (ss StockService) UpdateStock(si *model.StockInfo) error {
	if si.Code == "" {
		return exception.NewBusiness(400, "code is required")
	}
	if si.Name == "" {
		return exception.NewBusiness(400, "name is required")
	}

	osi, err := ss.sr.GetStock(si.Code)
	if err != nil {
		return err
	}

	osi.UpdatedAt = time.Now()
	osi.Currency = si.Currency
	osi.Market = si.Market
	osi.Name = si.Name
	osi.Status = 0

	return ss.sr.UpdateStock(osi)
}

func (ss StockService) DeleteStock(code string) error {
	if code == "" {
		return exception.NewBusiness(400, "code is required")
	}
	return ss.sr.DeleteStock(code)
}

func (ss StockService) GetHolding(code string) (*model.Investment, error) {
	if code == "" {
		return nil, exception.NewBusiness(400, "code is required")
	}

	invest, err := ss.sr.GetHolding(code)
	if err != nil {
		return nil, err
	}

	openTime, _ := time.Parse(DateTimeLayout, invest.OpenTime)
	var closeTime time.Time
	if invest.CloseTime == "" {
		closeTime = time.Now()
	} else {
		closeTime, _ = time.Parse(DateTimeLayout, invest.CloseTime)
	}
	invest.HoldingDays = int(closeTime.Sub(openTime).Hours() / 24)

	return invest, nil
}

func (ss StockService) DeleteTransaction(tranId int64) error {
	tran, err := ss.sr.GetTransaction(tranId)
	if err != nil {
		return exception.WrapService(500, "dao error", err)
	}
	err = ss.sr.DeleteTransaction(tranId)
	if err != nil {
		return exception.WrapService(500, "dao error", err)
	}
	return ss.computeHolding(tran.InvestID)
}

func (ss StockService) UpdateTransaction(tran *model.Transaction) error {
	if tran.ID == 0 {
		return exception.NewService(400, "transaction id is required")
	}
	otran, err := ss.sr.GetTransaction(tran.ID)
	if err != nil {
		return exception.WrapService(500, "dao error", err)
	}
	if otran == nil {
		return exception.NewBusiness(404, "transaction not found")
	}

	otran.Action = tran.Action
	otran.Price = tran.Price
	otran.Quantity = tran.Quantity
	otran.TaxFee = tran.TaxFee
	otran.Amount = floatMul(tran.Price, tran.Quantity)
	otran.FinishTime = tran.FinishTime
	otran.UpdatedAt = time.Now()
	err = ss.sr.UpdateTransaction(otran)
	if err != nil {
		return exception.WrapService(500, "dao error", err)
	}

	return ss.computeHolding(otran.InvestID)
}

func floatMul(p float64, q int) float64 {
	pd := decimal.NewFromFloat(p)
	qd := decimal.NewFromInt(int64(q))
	// 转换为 float64 存储到结构体字段
	return pd.Mul(qd).RoundBank(3).InexactFloat64()
}

func (ss StockService) AddTransaction(tran *model.Transaction) error {
	if tran.StockCode == "" {
		return exception.NewBusiness(400, "stock code is empty")
	}
	if tran.Action == 0 {
		return exception.NewBusiness(400, "action is empty")
	}

	nowTime := time.Now()

	invest, err := ss.sr.GetHolding(tran.StockCode)
	if err != nil {
		// no holding investment
		if tran.Action == -1 {
			return exception.NewBusiness(400, "action is sell but holding is closed")
		}
		// add holding investment for opening
		invest = &model.Investment{StockCode: tran.StockCode, Status: 0,
			CreatedAt: nowTime, UpdatedAt: nowTime, OpenTime: nowTime.Format(DateTimeLayout)}
		err := ss.sr.CreateInvestment(invest)
		if err != nil {
			return exception.WrapService(500, "create holding error", err)
		}
	} else {
		if tran.Action == -1 && invest.Quantity < tran.Quantity {
			return exception.NewBusiness(400, "holding quantity is less than sell quantity")
		}
	}

	tran.InvestID = invest.ID
	tran.CreatedAt = nowTime
	tran.UpdatedAt = nowTime
	tran.Amount = floatMul(tran.Price, tran.Quantity)

	err = ss.sr.CreateTransaction(tran)
	if err != nil {
		return err
	}

	// 根据持仓的交易记录计算持仓信息
	return ss.computeHolding(tran.InvestID)
}

func (ss StockService) computeHolding(investId int64) error {
	invest, err := ss.sr.GetInvestment(investId)
	if err != nil {
		return exception.WrapService(500, "dao error", err)
	}

	trans, err := ss.sr.GetTransactions(investId)
	if err != nil {
		return exception.WrapService(500, "dao error", err)
	}

	inQuantity := 0
	ouQuantity := 0

	inAmount := decimal.NewFromFloat(0)
	ouAmount := decimal.NewFromFloat(0)
	totalTaxFee := decimal.NewFromFloat(0)

	for _, t := range *trans {
		if t.Action == 1 {
			inQuantity += t.Quantity
			inAmount = inAmount.Add(decimal.NewFromFloat(t.Price).Mul(decimal.NewFromInt(int64(t.Quantity))))
		} else if t.Action == -1 {
			ouQuantity += t.Quantity
			ouAmount = ouAmount.Add(decimal.NewFromFloat(t.Price).Mul(decimal.NewFromInt(int64(t.Quantity))))
		}
		totalTaxFee = totalTaxFee.Add(decimal.NewFromFloat(t.TaxFee))
	}

	invest.Quantity = inQuantity - ouQuantity

	invest.TotalTaxFee = totalTaxFee.InexactFloat64()
	if invest.Quantity == 0 {
		invest.CostPrice = 0
	} else {
		invest.CostPrice = inAmount.Div(decimal.NewFromInt(int64(inQuantity))).RoundBank(3).InexactFloat64()
	}
	invest.Amount = floatMul(invest.CostPrice, invest.Quantity)
	invest.ProfitLoss = ouAmount.Sub(inAmount).Add(decimal.NewFromFloat(invest.Amount)).InexactFloat64()

	if invest.Quantity == 0 {
		invest.Status = 1

		invest.CloseTime = time.Now().Format(DateTimeLayout)
	}

	invest.UpdatedAt = time.Now()

	err = ss.sr.UpdateInvestment(invest)
	if err != nil {
		return exception.WrapService(500, "update invest error", err)
	}

	return nil
}

func (ss StockService) GetTransactions(investId int64) (*[]model.Transaction, error) {
	trans, err := ss.sr.GetTransactions(investId)
	if err != nil {
		return nil, exception.WrapService(500, "dao error", err)
	}
	return trans, nil
}

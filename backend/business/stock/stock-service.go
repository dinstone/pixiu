package stock

import (
	"pixiu/backend/pkg/exception"
	"pixiu/backend/pkg/gormer"
	"time"

	"github.com/shopspring/decimal"
)

const DateTimeLayout = "2006-01-02 15:04:05"

type StockService struct {
	gtm gormer.GormTM
	sr  StockRepository
}

func NewStockService(gtm gormer.GormTM, sr StockRepository) *StockService {
	return &StockService{gtm, sr}
}

func (ss *StockService) GetClearList(stime string, ftime string) (*[]ClearStats, error) {
	return ss.sr.GetClearList(ss.gtm.Context(), stime, ftime)
}

func (ss *StockService) GetStockClear(stockCode string, startTime string, finishTime string) (*ClearInvest, error) {
	if stockCode == "" {
		return nil, exception.NewBusiness(400, "stock code is required")
	}
	sinfo, err := ss.sr.GetStock(ss.gtm.Context(), stockCode)
	if err != nil {
		return nil, err
	}
	cinvests, err := ss.sr.GetClearInvest(ss.gtm.Context(), stockCode, startTime, finishTime)
	if err != nil {
		return nil, err
	}

	var totalCount int
	var profitLoss float64
	var totalAmount float64
	var invests []Investment
	for _, ci := range *cinvests {
		totalCount++
		profitLoss += ci.ProfitLoss
		totalAmount += ci.Amount
		openTime, _ := time.Parse(DateTimeLayout, ci.OpenTime)
		closeTime, _ := time.Parse(DateTimeLayout, ci.CloseTime)
		ci.HoldingDays = daysBetweenDates(openTime, closeTime)
		invests = append(invests, ci)
	}
	roi := 0.00
	if totalAmount != 0 {
		roi = round2Decimal((profitLoss / totalAmount) * 100)
	}
	return &ClearInvest{
		Stock:   sinfo,
		Stats:   &ClearStats{TotalCount: totalCount, ProfitLoss: profitLoss, Roi: roi, StartTime: startTime, FinishTime: finishTime},
		Invests: &invests}, nil
}

// daysBetweenDates 计算两个时间点之间相差的天数（只看日期，忽略时分秒）
func daysBetweenDates(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, t2.Location())
	return int(t2.Sub(t1).Hours() / 24)
}

func (ss StockService) GetStockList() (*[]StockInfo, error) {
	return ss.sr.AliveStocks(ss.gtm.Context())
}

func (ss StockService) GetStock(code string) (*StockInfo, error) {
	if code == "" {
		return nil, exception.NewBusiness(400, "code is required")
	}
	return ss.sr.GetStock(ss.gtm.Context(), code)
}

func (ss StockService) SaveStock(si *StockInfo) error {
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
	return ss.sr.SaveStock(ss.gtm.Context(), si)
}

func (ss StockService) UpdateStock(si *StockInfo) error {
	if si.Code == "" {
		return exception.NewBusiness(400, "code is required")
	}
	if si.Name == "" {
		return exception.NewBusiness(400, "name is required")
	}

	osi, err := ss.sr.GetStock(ss.gtm.Context(), si.Code)
	if err != nil {
		return err
	}

	osi.UpdatedAt = time.Now()
	osi.Currency = si.Currency
	osi.Market = si.Market
	osi.Name = si.Name
	osi.Status = 0

	return ss.sr.UpdateStock(ss.gtm.Context(), osi)
}

func (ss StockService) DeleteStock(code string) error {
	if code == "" {
		return exception.NewBusiness(400, "code is required")
	}
	return ss.sr.DeleteStock(ss.gtm.Context(), code)
}

func (ss StockService) GetHolding(code string) (*Investment, error) {
	if code == "" {
		return nil, exception.NewBusiness(400, "code is required")
	}

	invest, err := ss.sr.GetHolding(ss.gtm.Context(), code)
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
	tran, err := ss.sr.GetTransaction(ss.gtm.Context(), tranId)
	if err != nil {
		return exception.WrapService(500, "dao error", err)
	}
	err = ss.sr.DeleteTransaction(ss.gtm.Context(), tranId)
	if err != nil {
		return exception.WrapService(500, "dao error", err)
	}
	return ss.computeHolding(tran.InvestID)
}

func (ss StockService) UpdateTransaction(tran *Transaction) error {
	if tran.ID == 0 {
		return exception.NewService(400, "transaction id is required")
	}
	otran, err := ss.sr.GetTransaction(ss.gtm.Context(), tran.ID)
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
	otran.Amount = floatMulInt(tran.Price, tran.Quantity)
	otran.FinishTime = tran.FinishTime
	otran.UpdatedAt = time.Now()
	err = ss.sr.UpdateTransaction(ss.gtm.Context(), otran)
	if err != nil {
		return exception.WrapService(500, "dao error", err)
	}

	return ss.computeHolding(otran.InvestID)
}

func floatMulInt(p float64, q int) float64 {
	pd := decimal.NewFromFloat(p)
	qd := decimal.NewFromInt(int64(q))
	// 转换为 float64 存储到结构体字段
	return pd.Mul(qd).RoundBank(2).InexactFloat64()
}

func round2Decimal(value float64) float64 {
	d := decimal.NewFromFloat(value)
	return d.RoundBank(2).InexactFloat64()
}

func (ss StockService) AddTransaction(tran *Transaction) error {
	if tran.StockCode == "" {
		return exception.NewBusiness(400, "stock code is empty")
	}
	if tran.Action == 0 {
		return exception.NewBusiness(400, "action is empty")
	}

	nowTime := time.Now()

	invest, err := ss.sr.GetHolding(ss.gtm.Context(), tran.StockCode)
	if err != nil {
		// no holding investment
		if tran.Action == -1 {
			return exception.NewBusiness(400, "action is sell but holding is closed")
		}
		// add holding investment for opening
		invest = &Investment{StockCode: tran.StockCode, Status: 0,
			CreatedAt: nowTime, UpdatedAt: nowTime, OpenTime: nowTime.Format(DateTimeLayout)}
		err := ss.sr.CreateInvestment(ss.gtm.Context(), invest)
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
	tran.Amount = floatMulInt(tran.Price, tran.Quantity)
	if tran.TaxFee == 0 {
		tran.TaxFee = round2Decimal(tran.Amount * 0.00137)
	}

	err = ss.sr.CreateTransaction(ss.gtm.Context(), tran)
	if err != nil {
		return err
	}

	// 根据持仓的交易记录计算持仓信息
	return ss.computeHolding(tran.InvestID)
}

func (ss StockService) computeHolding(investId int64) error {
	invest, err := ss.sr.GetInvestment(ss.gtm.Context(), investId)
	if err != nil {
		return exception.WrapService(500, "dao error", err)
	}

	trans, err := ss.sr.GetTransactions(ss.gtm.Context(), investId)
	if err != nil {
		return exception.WrapService(500, "dao error", err)
	}

	inQuantity := 0
	ouQuantity := 0

	inAmount := decimal.NewFromFloat(0)
	ouAmount := decimal.NewFromFloat(0)
	totalTaxFee := decimal.NewFromFloat(0)

	for _, t := range *trans {
		switch t.Action {
		case 1:
			inQuantity += t.Quantity
			inAmount = inAmount.Add(decimal.NewFromFloat(t.Price).Mul(decimal.NewFromInt(int64(t.Quantity))))
		case -1:
			ouQuantity += t.Quantity
			ouAmount = ouAmount.Add(decimal.NewFromFloat(t.Price).Mul(decimal.NewFromInt(int64(t.Quantity))))
		}
		totalTaxFee = totalTaxFee.Add(decimal.NewFromFloat(t.TaxFee))
	}

	invest.TotalTaxFee = totalTaxFee.InexactFloat64()
	invest.CostPrice = inAmount.Div(decimal.NewFromInt(int64(inQuantity))).RoundBank(3).InexactFloat64()

	invest.Quantity = inQuantity - ouQuantity
	invest.Amount = inAmount.RoundBank(2).InexactFloat64()
	invest.ProfitLoss = ouAmount.Sub(inAmount).Add(decimal.NewFromFloat(invest.Amount)).InexactFloat64()

	// 第一个元素
	firstElement := (*trans)[0]
	invest.OpenTime = firstElement.FinishTime
	if invest.Quantity == 0 {
		// 清仓
		invest.Status = 1
		// 最后一个元素
		lastElement := (*trans)[len(*trans)-1]
		invest.CloseTime = lastElement.FinishTime
	}

	invest.UpdatedAt = time.Now()

	err = ss.sr.UpdateInvestment(ss.gtm.Context(), invest)
	if err != nil {
		return exception.WrapService(500, "update invest error", err)
	}

	return nil
}

func (ss StockService) GetTransactions(investId int64) (*[]Transaction, error) {
	trans, err := ss.sr.GetTransactions(ss.gtm.Context(), investId)
	if err != nil {
		return nil, exception.WrapService(500, "dao error", err)
	}
	return trans, nil
}

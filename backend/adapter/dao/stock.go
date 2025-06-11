package dao

import (
	"context"
	"pixiu/backend/business/model"
	"pixiu/backend/pkg/gormer"
)

type StockDao struct {
	ormer gormer.GormDS
}

func NewStockDao(gds gormer.GormDS) *StockDao {
	return &StockDao{ormer: gds}
}

func (s StockDao) AliveStocks(ctx context.Context) (*[]model.StockInfo, error) {
	var stocks []model.StockInfo
	err := s.ormer.GDB(ctx).Where("status = ?", 0).Order("market").Find(&stocks).Error
	return &stocks, WrapGormError(err)
}

func (s StockDao) GetStock(ctx context.Context, code string) (*model.StockInfo, error) {
	var stock model.StockInfo
	err := s.ormer.GDB(ctx).Where("code = ?", code).First(&stock).Error
	if err != nil {
		return nil, WrapGormError(err)
	}
	return &stock, nil
}

func (s StockDao) SaveStock(ctx context.Context, si *model.StockInfo) error {
	return WrapGormError(s.ormer.GDB(ctx).Create(si).Error)
}

func (s StockDao) UpdateStock(ctx context.Context, si *model.StockInfo) error {
	return WrapGormError(s.ormer.GDB(ctx).Model(si).Where("code=?", si.Code).Updates(si).Error)
}

func (s StockDao) DeleteStock(ctx context.Context, code string) error {
	return WrapGormError(s.ormer.GDB(ctx).Model(&model.StockInfo{}).Where("code = ?", code).UpdateColumn("status", -1).Error)
}

func (s StockDao) GetHolding(ctx context.Context, code string) (*model.Investment, error) {
	var investment model.Investment
	err := s.ormer.GDB(ctx).Where("stock_code = ? and status = 0", code).First(&investment).Error
	if err != nil {
		return nil, WrapGormError(err)
	}
	return &investment, nil
}

func (s StockDao) CreateInvestment(ctx context.Context, invest *model.Investment) error {
	return WrapGormError(s.ormer.GDB(ctx).Create(invest).Error)
}

func (s StockDao) UpdateInvestment(ctx context.Context, invest *model.Investment) error {
	return WrapGormError(s.ormer.GDB(ctx).Model(invest).Select("ProfitLoss", "TotalTaxFee", "CostPrice", "Quantity", "Amount", "OpenTime", "CloseTime", "Status").Updates(invest).Error)
}

func (s StockDao) GetInvestment(ctx context.Context, id int64) (*model.Investment, error) {
	var investment model.Investment
	err := s.ormer.GDB(ctx).Where("id = ?", id).First(&investment).Error
	return &investment, WrapGormError(err)
}

func (s StockDao) DeleteInvestment(ctx context.Context, id int64) error {
	return WrapGormError(s.ormer.GDB(ctx).Where("id = ?", id).Delete(&model.Investment{}).Error)
}

func (s StockDao) CreateTransaction(ctx context.Context, trans *model.Transaction) error {
	return WrapGormError(s.ormer.GDB(ctx).Create(trans).Error)
}

func (s StockDao) UpdateTransaction(ctx context.Context, trans *model.Transaction) error {
	return WrapGormError(s.ormer.GDB(ctx).Model(trans).Select("TaxFee", "Action", "Price", "Quantity", "Amount", "FinishTime", "UpdatedAt").Updates(trans).Error)
}

func (s StockDao) GetTransaction(ctx context.Context, id int64) (*model.Transaction, error) {
	var transaction model.Transaction
	err := s.ormer.GDB(ctx).Where("id = ?", id).First(&transaction).Error
	return &transaction, WrapGormError(err)
}

func (s StockDao) DeleteTransaction(ctx context.Context, id int64) error {
	return WrapGormError(s.ormer.GDB(ctx).Where("id = ?", id).Delete(model.Transaction{}).Error)
}

func (s StockDao) GetTransactions(ctx context.Context, investId int64) (*[]model.Transaction, error) {
	var transactions []model.Transaction
	err := s.ormer.GDB(ctx).Where("invest_id = ?", investId).Order("finish_time").Find(&transactions).Error
	return &transactions, WrapGormError(err)
}

func (s StockDao) GetClearList(ctx context.Context, stime string, ftime string) (*[]model.ClearStats, error) {
	subQuery := s.ormer.GDB(ctx).Model(model.Investment{}).
		Select("stock_code, Sum(profit_loss) profit_loss, COUNT(*) total_count, SUM(CASE WHEN profit_loss >= 0 THEN 1 ELSE 0 END) profit_count, SUM(CASE WHEN profit_loss < 0 THEN 1 ELSE 0 END) loss_count").
		Where("status = ?", 1)
	if stime != "" {
		subQuery = subQuery.Where("open_time >= ?", stime)
	}
	if ftime != "" {
		subQuery = subQuery.Where("open_time <= ?", ftime)
	}
	subQuery = subQuery.Group("stock_code")

	var clears []model.ClearStats
	err := s.ormer.GDB(ctx).Model(&model.StockInfo{}).
		Select("code stock_code, name stock_name, i.profit_loss, i.total_count, i.profit_count, i.loss_count").
		Joins("JOIN (?) i ON code = i.stock_code", subQuery).
		Find(&clears).Error
	return &clears, WrapGormError(err)
}

func (s *StockDao) GetClearInvest(ctx context.Context, stockCode string, startTime string, finishTime string) (*[]model.Investment, error) {
	db := s.ormer.GDB(ctx).Model(&model.Investment{}).Where("status=1 and stock_code = ?", stockCode)
	if startTime != "" {
		db = db.Where("open_time >= ?", startTime)
	}
	if finishTime != "" {
		db = db.Where("open_time <= ?", finishTime)
	}
	var invests []model.Investment
	err := db.Find(&invests).Error
	return &invests, WrapGormError(err)
}

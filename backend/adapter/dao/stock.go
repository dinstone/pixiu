package dao

import (
	"pixiu/backend/business/model"

	"gorm.io/gorm"
)

type StockDao struct {
	ormer *gorm.DB
}

func NewStockDao(gdb *gorm.DB) *StockDao {
	return &StockDao{ormer: gdb}
}

func (s StockDao) AliveStocks() (*[]model.StockInfo, error) {
	var stocks []model.StockInfo
	err := s.ormer.Where("status = ?", 0).Order("market").Find(&stocks).Error
	return &stocks, WrapGormError(err)
}

func (s StockDao) GetStock(code string) (*model.StockInfo, error) {
	var stock model.StockInfo
	err := s.ormer.Where("code = ?", code).First(&stock).Error
	if err != nil {
		return nil, WrapGormError(err)
	}
	return &stock, nil
}

func (s StockDao) SaveStock(si *model.StockInfo) error {
	return WrapGormError(s.ormer.Create(si).Error)
}

func (s StockDao) UpdateStock(si *model.StockInfo) error {
	return WrapGormError(s.ormer.Model(si).Where("code=?", si.Code).Updates(si).Error)
}

func (s StockDao) DeleteStock(code string) error {
	return WrapGormError(s.ormer.Model(&model.StockInfo{}).Where("code = ?", code).UpdateColumn("status", -1).Error)
}

func (s StockDao) GetHolding(code string) (*model.Investment, error) {
	var investment model.Investment
	err := s.ormer.Where("stock_code = ? and status = 0", code).First(&investment).Error
	if err != nil {
		return nil, WrapGormError(err)
	}
	return &investment, nil
}

func (s StockDao) CreateInvestment(invest *model.Investment) error {
	return WrapGormError(s.ormer.Create(invest).Error)
}

func (s StockDao) UpdateInvestment(invest *model.Investment) error {
	return WrapGormError(s.ormer.Model(invest).Select("ProfitLoss", "TotalTaxFee", "CostPrice", "Quantity", "Amount", "OpenTime", "CloseTime", "Status").Updates(invest).Error)
}

func (s StockDao) GetInvestment(id int64) (*model.Investment, error) {
	var investment model.Investment
	err := s.ormer.Where("id = ?", id).First(&investment).Error
	return &investment, WrapGormError(err)
}

func (s StockDao) DeleteInvestment(id int64) error {
	return WrapGormError(s.ormer.Where("id = ?", id).Delete(&model.Investment{}).Error)
}

func (s StockDao) CreateTransaction(trans *model.Transaction) error {
	return WrapGormError(s.ormer.Create(trans).Error)
}

func (s StockDao) UpdateTransaction(trans *model.Transaction) error {
	return WrapGormError(s.ormer.Model(trans).Select("TaxFee", "Action", "Price", "Quantity", "Amount", "FinishTime", "UpdatedAt").Updates(trans).Error)
}

func (s StockDao) GetTransaction(id int64) (*model.Transaction, error) {
	var transaction model.Transaction
	err := s.ormer.Where("id = ?", id).First(&transaction).Error
	return &transaction, WrapGormError(err)
}

func (s StockDao) DeleteTransaction(id int64) error {
	return WrapGormError(s.ormer.Where("id = ?", id).Delete(model.Transaction{}).Error)
}

func (s StockDao) GetTransactions(investId int64) (*[]model.Transaction, error) {
	var transactions []model.Transaction
	err := s.ormer.Where("invest_id = ?", investId).Order("finish_time").Find(&transactions).Error
	return &transactions, WrapGormError(err)
}

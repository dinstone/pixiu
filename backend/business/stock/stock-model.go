package stock

import (
	"time"
)

// 定义股票信息结构体
type StockInfo struct {
	Code      string    `gorm:"primaryKey" json:"code"` // 编码
	Name      string    `json:"name"`                   // 名称
	Market    string    `json:"market"`                 // 股市（A股、港股等）
	Currency  string    `json:"currency"`               // 币种（人民币、港币、美元等）
	Status    int       `json:"status"`                 // 状态（-1:删除、0:正常）
	CreatedAt time.Time `json:"createdAt"`              // 创建时间
	UpdatedAt time.Time `json:"updatedAt"`              // 更新时间
}

// 投资信息结构体
type Investment struct {
	ID          int64     `gorm:"primaryKey" json:"id"` // 标识（唯一标识符）
	StockCode   string    `json:"stockCode"`            // 股票编码
	ProfitLoss  float64   `json:"profitLoss"`           // 持仓盈亏金额
	TotalTaxFee float64   `json:"totalTaxFee"`          // 税费合计
	CostPrice   float64   `json:"costPrice"`            // 成本价格
	Quantity    int       `json:"quantity"`             // 持仓数量
	Amount      float64   `json:"amount"`               // 投资金额
	Status      int       `json:"status"`               // 状态（-1:删除、0:持仓、1:清仓）
	HoldingDays int       `gorm:"-" json:"holdingDays"` // 持仓天数
	OpenTime    string    `json:"openTime"`             // 建仓时间
	CloseTime   string    `json:"closeTime"`            // 清仓时间
	CreatedAt   time.Time `json:"createdAt"`            // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`            // 更新时间
}

// 交易信息结构体
type Transaction struct {
	ID         int64     `gorm:"primaryKey" json:"id"` // 标识（唯一标识符）
	InvestID   int64     `json:"investId"`             // 投资标识（关联 Investment 结构体的 ID）
	StockCode  string    `json:"stockCode"`            // 股票编码
	Action     int8      `json:"action"`               // 操作类型：买入:1、删除:0、卖出:-1
	TaxFee     float64   `json:"taxFee"`               // 税费（如交易税、手续费等）
	Price      float64   `json:"price"`                // 成交价格
	Quantity   int       `json:"quantity"`             // 成交数量
	Amount     float64   `json:"amount"`               // 交易金额
	FinishTime string    `json:"finishTime"`           // 成交时间
	CreatedAt  time.Time `json:"createdAt"`            // 创建时间
	UpdatedAt  time.Time `json:"updatedAt"`            // 更新时间
}

type ClearStats struct {
	StockCode   string  `json:"stockCode"`
	StockName   string  `json:"stockName"`
	ProfitLoss  float64 `json:"profitLoss"`
	Roi         float64 `json:"roi"`
	TotalCount  int     `json:"totalCount"`
	ProfitCount int     `json:"profitCount"`
	LossCount   int     `json:"lossCount"`
	StartTime   string  `json:"startTime"`
	FinishTime  string  `json:"finishTime"`
}

type ClearInvest struct {
	Stock   *StockInfo    `json:"stock"`
	Stats   *ClearStats   `json:"stats"`
	Invests *[]Investment `json:"invests"`
}

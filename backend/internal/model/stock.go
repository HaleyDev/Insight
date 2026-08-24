package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// BasicStockDataModel 股票基础数据表
type BasicStockDataModel struct {
	ID            uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	Date          time.Time `gorm:"column:date" json:"date"`
	TradingVolume float64   `gorm:"column:trading_volume" json:"trading_volume"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"-"`
}

// TableName 表名
func (s *BasicStockDataModel) TableName() string {
	return "stock_data"
}

// StockDataInfo 对外暴露的股票数据结构
type StockDataInfo struct {
	ID            uint64    `json:"id"`
	Date          time.Time `json:"date"`
	TradingVolume float64   `json:"trading_volume"`
}

// ToInfo 转换为对外结构
func (s *BasicStockDataModel) ToInfo() *StockDataInfo {
	if s == nil {
		return &StockDataInfo{}
	}
	return &StockDataInfo{
		ID:            s.ID,
		Date:          s.Date,
		TradingVolume: s.TradingVolume,
	}
}

// 市场常量
const (
	MarketAStock  = "A股"
	MarketHKStock = "港股"
	MarketUSStock = "美股"
)

// 币种常量
const (
	CurrencyCNY = "CNY"
	CurrencyUSD = "USD"
	CurrencyHKD = "HKD"
)

// 标的类型常量
const (
	SecurityTypeStock = "stock"
	SecurityTypeETF   = "ETF"
)

// SecuritiesModel 标的基础信息表（静态元数据）
type SecuritiesModel struct {
	ID        uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`                             // 主键
	Symbol    string    `gorm:"column:symbol;size:20;not null;uniqueIndex:idx_symbol_market" json:"symbol"` // 代码：588120, AAPL, 00700
	Name      string    `gorm:"column:name;size:50;not null" json:"name"`                                   // 名称：科创东财, 苹果, 腾讯
	Market    string    `gorm:"column:market;size:10;not null;uniqueIndex:idx_symbol_market" json:"market"` // 市场：A股 / 港股 / 美股
	Currency  string    `gorm:"column:currency;size:3;not null" json:"currency"`                            // 币种：CNY / USD / HKD
	Type      string    `gorm:"column:type;size:20;not null;default:stock" json:"type"`                     // 类型：stock / ETF
	IsActive  bool      `gorm:"column:is_active;not null;default:true" json:"is_active"`                    // 是否在持仓中
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`                                                 // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`                                                 // 更新时间
}

// TableName 表名
func (s *SecuritiesModel) TableName() string {
	return "securities"
}

// SecurityInfo 对外暴露的标的结构
type SecurityInfo struct {
	ID       uint64 `json:"id"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Market   string `json:"market"`
	Currency string `json:"currency"`
	Type     string `json:"type"`
	IsActive bool   `json:"is_active"`
}

// ToInfo 转换为对外结构
func (s *SecuritiesModel) ToInfo() *SecurityInfo {
	if s == nil {
		return &SecurityInfo{}
	}
	return &SecurityInfo{
		ID:       s.ID,
		Symbol:   s.Symbol,
		Name:     s.Name,
		Market:   s.Market,
		Currency: s.Currency,
		Type:     s.Type,
		IsActive: s.IsActive,
	}
}

// 交易类型常量
const (
	TradeTypeBuy  = "买入"
	TradeTypeSell = "卖出"
)

// 交易平台常量
const (
	PlatformEastMoney = "东方财富"
	PlatformIBKR      = "IBKR"
)

// TradesModel 交易记录表（唯一数据源，所有持仓和盈亏都从这张表计算）
type TradesModel struct {
	ID          uint64           `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`                       // 主键
	TradeDate   time.Time        `gorm:"column:trade_date;type:date;not null" json:"trade_date"`               // 交易日期
	TradeTime   *time.Time       `gorm:"column:trade_time;type:time" json:"trade_time,omitempty"`              // 交易时间（可选）
	Timezone    *string          `gorm:"column:timezone;size:32" json:"timezone,omitempty"`                    // 交易时区：Asia/Shanghai, America/New_York, Asia/Hong_Kong
	SecurityID  uint64           `gorm:"column:security_id;not null;index" json:"security_id"`                 // 关联 securities.id
	TradeType   string           `gorm:"column:trade_type;size:10;not null" json:"trade_type"`                 // 类型：买入 / 卖出
	Quantity    int32            `gorm:"column:quantity;not null" json:"quantity"`                             // 数量（正=买入，负=卖出）
	Price       decimal.Decimal  `gorm:"column:price;type:decimal(12,4);not null" json:"price"`                // 成交价格
	Amount      decimal.Decimal  `gorm:"column:amount;type:decimal(15,2);not null" json:"amount"`              // 成交金额 = |数量 × 价格|
	Fee         decimal.Decimal  `gorm:"column:fee;type:decimal(10,2);default:0" json:"fee"`                   // 手续费 / 佣金 / 印花税
	RealizedPnl *decimal.Decimal `gorm:"column:realized_pnl;type:decimal(15,2)" json:"realized_pnl,omitempty"` // 卖出时的已实现盈亏（历史事实）
	Platform    string           `gorm:"column:platform;size:20;not null" json:"platform"`                     // 平台：东方财富 / IBKR
	Notes       *string          `gorm:"column:notes;type:text" json:"notes,omitempty"`                        // 交易备注
	CreatedAt   time.Time        `gorm:"column:created_at" json:"-"`                                           // 创建时间
	UpdatedAt   time.Time        `gorm:"column:updated_at" json:"-"`                                           // 更新时间
}

// TableName 表名
func (t *TradesModel) TableName() string {
	return "trades"
}

// TradeInfo 对外暴露的交易记录结构
type TradeInfo struct {
	ID          uint64           `json:"id"`
	TradeDate   time.Time        `json:"trade_date"`
	TradeTime   *time.Time       `json:"trade_time,omitempty"`
	Timezone    *string          `json:"timezone,omitempty"`
	SecurityID  uint64           `json:"security_id"`
	TradeType   string           `json:"trade_type"`
	Quantity    int32            `json:"quantity"`
	Price       decimal.Decimal  `json:"price"`
	Amount      decimal.Decimal  `json:"amount"`
	Fee         decimal.Decimal  `json:"fee"`
	RealizedPnl *decimal.Decimal `json:"realized_pnl,omitempty"`
	Platform    string           `json:"platform"`
	Notes       *string          `json:"notes,omitempty"`
}

// ToInfo 转换为对外结构
func (t *TradesModel) ToInfo() *TradeInfo {
	if t == nil {
		return &TradeInfo{}
	}
	return &TradeInfo{
		ID:          t.ID,
		TradeDate:   t.TradeDate,
		TradeTime:   t.TradeTime,
		Timezone:    t.Timezone,
		SecurityID:  t.SecurityID,
		TradeType:   t.TradeType,
		Quantity:    t.Quantity,
		Price:       t.Price,
		Amount:      t.Amount,
		Fee:         t.Fee,
		RealizedPnl: t.RealizedPnl,
		Platform:    t.Platform,
		Notes:       t.Notes,
	}
}

// HoldingsModel 持仓表（从 trades 汇总而来的当前持仓缓存，只存稳定的派生值：数量、成本）
type HoldingsModel struct {
	ID            uint64          `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`                             // 主键
	SecurityID    uint64          `gorm:"column:security_id;not null;uniqueIndex:idx_security_id" json:"security_id"` // 关联 securities.id
	Quantity      int32           `gorm:"column:quantity;not null" json:"quantity"`                                   // 持仓数量
	CostPrice     decimal.Decimal `gorm:"column:cost_price;type:decimal(20,8);not null" json:"cost_price"`            // 加权平均成本价
	TotalCost     decimal.Decimal `gorm:"column:total_cost;type:decimal(18,2);not null" json:"total_cost"`            // 总成本 = quantity × cost_price
	RiskLevel     int32           `gorm:"column:risk_level;default:3" json:"risk_level"`                              // 风险评级 1-10
	Notes         *string         `gorm:"column:notes;type:text" json:"notes,omitempty"`                              // 备注：买入逻辑、目标价
	FirstBuyDate  *time.Time      `gorm:"column:first_buy_date;type:date" json:"first_buy_date,omitempty"`            // 首次买入日期
	LastTradeDate *time.Time      `gorm:"column:last_trade_date;type:date" json:"last_trade_date,omitempty"`          // 最后交易日期
	CreatedAt     time.Time       `gorm:"column:created_at" json:"-"`                                                 // 创建时间
	UpdatedAt     time.Time       `gorm:"column:updated_at" json:"-"`                                                 // 更新时间
}

// TableName 表名
func (h *HoldingsModel) TableName() string {
	return "holdings"
}

// HoldingInfo 对外暴露的持仓结构
type HoldingInfo struct {
	ID            uint64          `json:"id"`
	SecurityID    uint64          `json:"security_id"`
	Quantity      int32           `json:"quantity"`
	CostPrice     decimal.Decimal `json:"cost_price"`
	TotalCost     decimal.Decimal `json:"total_cost"`
	RiskLevel     int32           `json:"risk_level"`
	Notes         *string         `json:"notes,omitempty"`
	FirstBuyDate  *time.Time      `json:"first_buy_date,omitempty"`
	LastTradeDate *time.Time      `json:"last_trade_date,omitempty"`
}

// ToInfo 转换为对外结构
func (h *HoldingsModel) ToInfo() *HoldingInfo {
	if h == nil {
		return &HoldingInfo{}
	}
	return &HoldingInfo{
		ID:            h.ID,
		SecurityID:    h.SecurityID,
		Quantity:      h.Quantity,
		CostPrice:     h.CostPrice,
		TotalCost:     h.TotalCost,
		RiskLevel:     h.RiskLevel,
		Notes:         h.Notes,
		FirstBuyDate:  h.FirstBuyDate,
		LastTradeDate: h.LastTradeDate,
	}
}

// DailySnapshotsModel 每日收益快照表（业绩曲线核心，每天收盘后生成一条，用收盘价算出后冻结不变）
type DailySnapshotsModel struct {
	ID               uint64           `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`                                             // 主键
	SnapshotDate     time.Time        `gorm:"column:snapshot_date;type:date;not null;uniqueIndex:idx_snapshot_date" json:"snapshot_date"` // 快照日期
	TotalAssets      decimal.Decimal  `gorm:"column:total_assets;type:decimal(18,2);not null" json:"total_assets"`                        // 当日总资产（含现金）
	TotalMarketValue *decimal.Decimal `gorm:"column:total_market_value;type:decimal(18,2)" json:"total_market_value,omitempty"`           // 持仓总市值（用收盘价算）
	TotalCost        *decimal.Decimal `gorm:"column:total_cost;type:decimal(18,2)" json:"total_cost,omitempty"`                           // 持仓总成本
	DailyPnl         decimal.Decimal  `gorm:"column:daily_pnl;type:decimal(15,2);not null" json:"daily_pnl"`                              // 当日盈亏
	DailyPnlPct      *decimal.Decimal `gorm:"column:daily_pnl_pct;type:decimal(8,4)" json:"daily_pnl_pct,omitempty"`                      // 当日收益率
	CumulativePnl    *decimal.Decimal `gorm:"column:cumulative_pnl;type:decimal(18,2)" json:"cumulative_pnl,omitempty"`                   // 累计盈亏
	CumulativePnlPct *decimal.Decimal `gorm:"column:cumulative_pnl_pct;type:decimal(8,4)" json:"cumulative_pnl_pct,omitempty"`            // 累计收益率
	InitialCapital   decimal.Decimal  `gorm:"column:initial_capital;type:decimal(18,2);not null" json:"initial_capital"`                  // 初始本金
	CashBalance      *decimal.Decimal `gorm:"column:cash_balance;type:decimal(18,2)" json:"cash_balance,omitempty"`                       // 现金余额
	PositionRatio    *decimal.Decimal `gorm:"column:position_ratio;type:decimal(5,4)" json:"position_ratio,omitempty"`                    // 仓位比例 = 总市值 / 总资产
	HoldingCount     *int32           `gorm:"column:holding_count" json:"holding_count,omitempty"`                                        // 持仓数量
	Notes            *string          `gorm:"column:notes;type:text" json:"notes,omitempty"`                                              // 备注：重大事件、操作说明
	CreatedAt        time.Time        `gorm:"column:created_at" json:"-"`                                                                 // 创建时间
}

// TableName 表名
func (d *DailySnapshotsModel) TableName() string {
	return "daily_snapshots"
}

// DailySnapshotInfo 对外暴露的每日快照结构
type DailySnapshotInfo struct {
	ID               uint64           `json:"id"`
	SnapshotDate     time.Time        `json:"snapshot_date"`
	TotalAssets      decimal.Decimal  `json:"total_assets"`
	TotalMarketValue *decimal.Decimal `json:"total_market_value,omitempty"`
	TotalCost        *decimal.Decimal `json:"total_cost,omitempty"`
	DailyPnl         decimal.Decimal  `json:"daily_pnl"`
	DailyPnlPct      *decimal.Decimal `json:"daily_pnl_pct,omitempty"`
	CumulativePnl    *decimal.Decimal `json:"cumulative_pnl,omitempty"`
	CumulativePnlPct *decimal.Decimal `json:"cumulative_pnl_pct,omitempty"`
	InitialCapital   decimal.Decimal  `json:"initial_capital"`
	CashBalance      *decimal.Decimal `json:"cash_balance,omitempty"`
	PositionRatio    *decimal.Decimal `json:"position_ratio,omitempty"`
	HoldingCount     *int32           `json:"holding_count,omitempty"`
	Notes            *string          `json:"notes,omitempty"`
}

// ToInfo 转换为对外结构
func (d *DailySnapshotsModel) ToInfo() *DailySnapshotInfo {
	if d == nil {
		return &DailySnapshotInfo{}
	}
	return &DailySnapshotInfo{
		ID:               d.ID,
		SnapshotDate:     d.SnapshotDate,
		TotalAssets:      d.TotalAssets,
		TotalMarketValue: d.TotalMarketValue,
		TotalCost:        d.TotalCost,
		DailyPnl:         d.DailyPnl,
		DailyPnlPct:      d.DailyPnlPct,
		CumulativePnl:    d.CumulativePnl,
		CumulativePnlPct: d.CumulativePnlPct,
		InitialCapital:   d.InitialCapital,
		CashBalance:      d.CashBalance,
		PositionRatio:    d.PositionRatio,
		HoldingCount:     d.HoldingCount,
		Notes:            d.Notes,
	}
}

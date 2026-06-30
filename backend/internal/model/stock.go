package model

import "time"

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
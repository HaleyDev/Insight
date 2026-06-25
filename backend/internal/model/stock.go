package model

import "time"

type BasicStockData struct {
	ID            uint64    `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	Date          time.Time `gorm:"column:date" json:"date"`
	TradingVolume float64   `gorm:"column:trading_volume" json:"trading_volume"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"-"`
}

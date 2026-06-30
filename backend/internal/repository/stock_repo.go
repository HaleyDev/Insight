package repository

import (
	"context"
	"time"

	"github.com/insight/backend/internal/model"
	"github.com/pkg/errors"
)

// ListStockData 根据时间范围获取市场信息
func (d *repository) ListStockData(ctx context.Context, startTime, endTime time.Time) ([]*model.BasicStockDataModel, error) {
	var stocks []*model.BasicStockDataModel
	err := d.orm.WithContext(ctx).
		Where("date BETWEEN ? AND ?", startTime, endTime).
		Order("date ASC").
		Find(&stocks).Error
	if err != nil {
		return nil, errors.Wrap(err, "[repo.stock] list stock data err")
	}
	return stocks, nil
}

func (d *repository) AddStockData(ctx context.Context, stockData *model.BasicStockDataModel) (id uint64, err error) {
	err = d.orm.WithContext(ctx).Create(&stockData).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo.stock] add stock data err")
	}
	return stockData.ID, nil
}

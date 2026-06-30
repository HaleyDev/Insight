package service

import (
	"context"
	"time"

	pkgerrors "github.com/pkg/errors"

	"github.com/insight/backend/internal/model"
	"github.com/insight/backend/internal/repository"
)

// StockService 股票服务接口
type StockService interface {
	ListStockData(ctx context.Context, startTime, endTime time.Time) ([]*model.StockDataInfo, error)
	AddStockInfo(ctx context.Context, stockInfo *model.StockDataInfo) error
}

type stockService struct {
	repo repository.Repository
}

var _ StockService = (*stockService)(nil)

func newStocks(svc *service) *stockService {
	return &stockService{repo: svc.repo}
}

// ListStockData 根据时间范围获取市场数据
func (s *stockService) ListStockData(ctx context.Context, startTime, endTime time.Time) ([]*model.StockDataInfo, error) {
	stocks, err := s.repo.ListStockData(ctx, startTime, endTime)
	if err != nil {
		return nil, err
	}
	infos := make([]*model.StockDataInfo, 0, len(stocks))
	for _, stock := range stocks {
		infos = append(infos, stock.ToInfo())
	}
	return infos, nil
}

// AddStockInfo 新增市场数据
func (s *stockService) AddStockInfo(ctx context.Context, stockInfo *model.StockDataInfo) error {
	stock := &model.BasicStockDataModel{
		Date:          stockInfo.Date,
		TradingVolume: stockInfo.TradingVolume,
	}
	if _, err := s.repo.AddStockData(ctx, stock); err != nil {
		return pkgerrors.Wrap(err, "add stock data err")
	}
	return nil
}

package service

import (
	"github.com/insight/backend/internal/repository"
)

// Svc 全局 service
var Svc Service

// Service 业务服务汇总接口
type Service interface {
	Users() UserService
	Stocks() StockService
}

type service struct {
	repo repository.Repository
}

// New 创建 service
func New(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Users() UserService {
	return newUsers(s)
}

func (s *service) Stocks() StockService { return newStocks(s) }

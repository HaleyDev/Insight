package repository

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/insight/backend/internal/cache"
	"github.com/insight/backend/internal/model"
	"github.com/insight/backend/pkg/storage/sql"
)

var (
	// ErrNotFound 数据不存在
	ErrNotFound = gorm.ErrRecordNotFound
)

var _ Repository = (*repository)(nil)

// Repository 用户仓库接口
type Repository interface {
	CreateUser(ctx context.Context, user *model.UserBaseModel) (id uint64, err error)
	UpdateUser(ctx context.Context, id uint64, userMap map[string]interface{}) error
	GetUser(ctx context.Context, id uint64) (*model.UserBaseModel, error)
	GetUserByEmail(ctx context.Context, email string) (*model.UserBaseModel, error)
	UserIsExist(user *model.UserBaseModel) (bool, error)
	Close()
}

// repository 仓库实现
type repository struct {
	orm       *gorm.DB
	db        *sql.DB
	tracer    trace.Tracer
	userCache *cache.Cache
}

// New 创建仓库实例
func New(db *gorm.DB) Repository {
	return &repository{
		orm:       db,
		tracer:    otel.Tracer("repository"),
		userCache: cache.NewUserCache(),
	}
}

// Close 释放连接
func (d *repository) Close() {}

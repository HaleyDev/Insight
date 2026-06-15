package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/insight/backend/internal/model"
	"github.com/insight/backend/pkg/cache"
	"github.com/insight/backend/pkg/log"
	"github.com/insight/backend/pkg/redis"
)

var g singleflight.Group

// CreateUser 创建用户
func (d *repository) CreateUser(ctx context.Context, user *model.UserBaseModel) (id uint64, err error) {
	err = d.orm.WithContext(ctx).Create(&user).Error
	if err != nil {
		return 0, errors.Wrap(err, "[repo.user_base] create user err")
	}
	return user.ID, nil
}

// UpdateUser 更新用户信息
func (d *repository) UpdateUser(ctx context.Context, id uint64, userMap map[string]interface{}) error {
	user, err := d.GetUser(ctx, id)
	if err != nil {
		return errors.Wrap(err, "[repo.user_base] update user data err")
	}

	if err := d.userCache.DelUserBaseCache(ctx, id); err != nil {
		log.Warnf("[repo.user_base] delete user cache err: %v", err)
	}

	if err := d.orm.WithContext(ctx).Model(user).Updates(userMap).Error; err != nil {
		return errors.Wrap(err, "[repo.user_base] update user err")
	}
	return nil
}

// GetUser 根据 id 获取用户（Cache Aside Pattern）
func (d *repository) GetUser(ctx context.Context, uid uint64) (userBase *model.UserBaseModel, err error) {
	ctx, span := d.tracer.Start(ctx, "GetUser", oteltrace.WithAttributes(
		attribute.String("param.uid", cast.ToString(uid)),
	))
	defer span.End()

	var (
		data *model.UserBaseModel
		val  interface{}
	)

	userBase, err = d.userCache.GetUserBaseCache(ctx, uid)
	if errors.Is(err, cache.ErrPlaceholder) {
		span.RecordError(err)
		return nil, ErrNotFound
	} else if errors.Is(err, redis.ErrRedisNotFound) {
		key := fmt.Sprintf("get_user_base_%d", uid)
		val, err, _ = g.Do(key, func() (interface{}, error) {
			data = new(model.UserBaseModel)
			err = d.orm.WithContext(ctx).First(data, uid).Error
			if errors.Is(err, ErrNotFound) {
				if err := d.userCache.SetCacheWithNotFound(ctx, uid); err != nil {
					log.Warnf("[repo.user_base] SetCacheWithNotFound err, uid: %d", uid)
				}
				return nil, ErrNotFound
			} else if err != nil {
				span.RecordError(err)
				return nil, errors.Wrapf(err, "[repo.user_base] query db err")
			}
			if err := d.userCache.SetUserBaseCache(ctx, uid, data, cache.DefaultExpireTime); err != nil {
				return nil, errors.Wrap(err, "[repo.user_base] SetUserBaseCache err")
			}
			return data, nil
		})

		if err != nil && err != ErrNotFound {
			span.RecordError(err)
			return nil, errors.Wrap(err, "[repo.user_base] get user base err via single flight do")
		}
		if val != nil {
			data = val.(*model.UserBaseModel)
		}
	} else if err != nil {
		return nil, err
	}

	if userBase != nil {
		return userBase, nil
	}
	return data, nil
}

// GetUserByEmail 根据邮箱获取用户
func (d *repository) GetUserByEmail(ctx context.Context, email string) (*model.UserBaseModel, error) {
	userBase := &model.UserBaseModel{}
	err := d.orm.WithContext(ctx).Where("email = ?", email).First(userBase).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "[repo.user_base] get user err by email")
	}
	return userBase, nil
}

// UserIsExist 判断用户名或邮箱是否已存在
func (d *repository) UserIsExist(user *model.UserBaseModel) (bool, error) {
	err := d.orm.Where("username = ? or email = ?", user.Username, user.Email).First(&model.UserBaseModel{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

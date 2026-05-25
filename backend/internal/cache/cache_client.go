package cache

import (
	"context"

	"github.com/insight/backend/internal/model"
	"github.com/insight/backend/pkg/cache"
	"github.com/insight/backend/pkg/encoding"
	"github.com/insight/backend/pkg/redis"
)

func getCacheClient(ctx context.Context) cache.Cache {
	sonicEncoding := encoding.SonicEncoding{}
	cachePrefix := ""
	client := cache.NewRedisCache(redis.RedisClient, cachePrefix, sonicEncoding, func() interface{} {
		return &model.UserBaseModel{}
	})

	return client
}

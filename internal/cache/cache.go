package cache

import (
	"context"

	"github.com/AlexandrBurak/TaxiApp/internal/config"

	"github.com/AlexandrBurak/TaxiApp/internal/model"
)

type Cache struct {
	Cache *redis.Client
}

func NewCache() (*Cache, error) {

	appCfg, err := config.GetAppCfg()
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     appCfg.REDIS_URL,
		Password: appCfg.REDIS_PASSWORD,
	})
	_, err = rdb.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}
	return &Cache{Cache: rdb}, nil
}

func (cache *Cache) CacheUser(ctx context.Context, user model.User) error {
	if _, err := cache.Cache.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		rdb.HSet(ctx, user.Phone, "phone", user.Phone)
		rdb.HSet(ctx, user.Phone, "password", user.Password)
		rdb.HSet(ctx, user.Phone, "name", user.Name)
		rdb.HSet(ctx, user.Phone, "email", user.Email)
		rdb.HSet(ctx, user.Phone, "id", user.ID)
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (cache *Cache) ExistsInCache(ctx context.Context, phone string) (model.User, error) {
	var user model.User
	if err := cache.Cache.HGetAll(ctx, phone).Scan(&user); err != nil {
		return user, err
	}
	return user, nil
}

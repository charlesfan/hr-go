package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/charlesfan/hr-go/config"
)

var item *cache

type ICache interface {
	Ping(context.Context) (string, error)
	Set(context.Context, string, interface{}, time.Duration) error
	BindJSON(context.Context, string, interface{}) (bool, error)
}

type cache struct {
	redis ICache
}

func (c *cache) Redis() ICache {
	return item.redis
}

func Init(c config.Config) {
	ropts := &redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	}
	item = &cache{
		redis: newRedis(ropts),
	}
}

func NewRedis() ICache {
	if item != nil {
		return item.Redis()
	}
	return nil
}

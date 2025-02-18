package fixture

import (
	"context"
	"time"

	"github.com/charlesfan/hr-go/repository/cache"
)

type cacheMock struct {
	whitelist map[string]interface{}
}

func (c *cacheMock) Ping(ctx context.Context) (string, error) {
	return "success", nil
}

func (c *cacheMock) Set(ctx context.Context, key string, val interface{}, du time.Duration) error {
	c.whitelist[key] = val
	return nil
}

func (c *cacheMock) BindJSON(ctx context.Context, key string, dest interface{}) (bool, error) {
	if _, ok := c.whitelist[key]; ok {
		return true, nil
	}
	return false, nil
}

func NewCacheMock() cache.ICache {
	return &cacheMock{
		whitelist: make(map[string]interface{}),
	}
}

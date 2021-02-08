package cache

import (
	"go-sso/pkg/redis"
	"time"
)

type (
	CacheClient interface {
		GetCache(key string) (string, error)
		SetCache(key string, value string) error
		SetCacheExpired(key string, value string, expired time.Duration) error
		RemoveCache(key string) (bool, error)
	}

	cacheClient struct {
		client redis.RedisClient
	}
)

func NewCacheClient(client redis.RedisClient) CacheClient {
	return newCacheClient(client)
}

func newCacheClient(client redis.RedisClient) *cacheClient {
	return &cacheClient{client: client}
}

func (c cacheClient) GetCache(key string) (string, error) {
	return c.client.Get(key)
}

func (c cacheClient) SetCache(key string, value string) error {
	_, err := c.client.Set(key, value)
	return err
}

func (c cacheClient) SetCacheExpired(key string, value string, expired time.Duration) error {
	_, err := c.client.SetExpired(key, value, int(expired/time.Second))
	return err
}

func (c cacheClient) RemoveCache(key string) (bool, error) {
	v, err := c.client.Del(key)
	if err != nil {
		return false, err
	}
	return v == 1, err
}

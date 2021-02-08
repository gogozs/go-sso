package registry

import (
	"go-sso/conf"
	"go-sso/internal/repository/cache"
	"go-sso/pkg/redis"
)

var (
	cacheStore cache.CacheClient
)

func InitCacheClient(c *conf.Config) cache.CacheClient {
	cacheStore = cache.NewCacheClient(
		redis.NewRedisClient(
			redis.SetHost(c.Redis.Host),
			redis.SetPort(c.Redis.Port),
			redis.SetPassword(c.Redis.Password),
		))
	return cacheStore
}

func GetCacheStore() cache.CacheClient {
	return cacheStore
}

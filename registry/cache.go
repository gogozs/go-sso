package registry

import (
	"github.com/go-zs/cache"
	"go-sso/pkg/storage"
)

var (
	cacheStore = cache.NewStore()
)

// 可以自行替换redis等
func GetCacheStore() storage.CacheStore {
	return cacheStore
}

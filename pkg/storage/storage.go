package storage

import "github.com/go-zs/cache"

var (
	cacheStore = cache.NewStore()
)

// 可以自行替换redis等
func GetStore() CacheStore {
	return cacheStore
}


type CacheStore interface {
	GetCache(string) (interface{}, bool)
	SetCache(string, interface{}) bool
	RemoveCache(string) bool
}


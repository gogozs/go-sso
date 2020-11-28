package storage

type CacheStore interface {
	GetCache(string) (interface{}, bool)
	SetCache(string, interface{}) bool
	RemoveCache(string) bool
}

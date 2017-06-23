package main

type (
	storeContext struct {
		CacheStorage *CacheStorage
	}
)

var StoreContext *storeContext

func init() {
	StoreContext = &storeContext{
		CacheStorage: NewCacheStorage(),
	}
}

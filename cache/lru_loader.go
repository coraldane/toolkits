package cache

//
//import (
//	lru "github.com/hashicorp/golang-lru"
//)
//
//type LRUCacheLoader[Key comparable, Value any] struct {
//	dataFetcher Fetcher[Key, Value]
//	lruCache    *lru.Cache
//}
//
//func NewLRULoader[Key comparable, Value any](fetcher Fetcher[Key, Value], capacity int) *LRUCacheLoader {
//	loader := LRUCacheLoader{}
//	loader.dataFetcher = fetcher
//	loader.lruCache, _ = lru.New(capacity)
//	return &loader
//}

package cache

import (
	"fmt"
	"github.com/linyerun/GeeCache/lru"
	"sync"
)

type safetyCache struct {
	mu       sync.RWMutex // 大部分使用写锁，唯一使用读锁就是获取Cache的Size
	lruCache lru.ICache
}

func NewSafetyCache(maxBytes int64, onEvicted func(key string, value IByteView)) ISafetyCache {
	if onEvicted == nil {
		onEvicted = func(key string, value IByteView) {
			fmt.Printf("execute default onEvicted: the key = %v has been deleted, the value is %v", key, value.String())
		}
	}
	return &safetyCache{
		mu: sync.RWMutex{},
		lruCache: lru.NewLruCache(maxBytes, func(key string, value lru.Value) {
			onEvicted(key, value.(IByteView))
		}),
	}
}

func (s *safetyCache) Get(key string) (value IByteView, ok bool) {
	// 写锁
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.lruCache.Get(key)
	if !ok {
		return nil, false
	}

	return v.(IByteView), ok
}

func (s *safetyCache) Add(key string, value IByteView) (err error) {
	// 写锁
	s.mu.Lock()
	defer s.mu.Unlock()

	err = s.lruCache.Add(key, value)
	return
}

func (s *safetyCache) Delete(key string) (err error) {
	// 写锁
	s.mu.Lock()
	defer s.mu.Unlock()

	err = s.lruCache.Delete(key)
	return
}

func (s *safetyCache) Size() uint {
	// 读锁
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lruCache.Size()
}

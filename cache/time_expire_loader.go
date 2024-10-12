package cache

import (
	"github.com/coraldane/logger"
	"github.com/coraldane/toolkits/concurrent"
	"time"
)

type TimeExpireCacheLoader struct {
	dataFetcher    Fetcher
	expireDuration time.Duration
	clearInterval  time.Duration
	closed         chan struct{}
	container      concurrent.Map
	expireKeyMap   concurrent.Map
}

func NewTimeExpireLoader(fetcher Fetcher,
	expire time.Duration, clearInterval time.Duration) *TimeExpireCacheLoader {
	loader := &TimeExpireCacheLoader{
		dataFetcher:    fetcher,
		expireDuration: expire,
		clearInterval:  clearInterval,
		closed:         make(chan struct{}),
		container:      concurrent.Map{},
		expireKeyMap:   concurrent.Map{},
	}
	go loader.checkExpire()
	//runtime.SetFinalizer(loader, loader.Close)
	return loader
}

func (this *TimeExpireCacheLoader) Close() {
	close(this.closed)
}

func (this *TimeExpireCacheLoader) Get(k any) (any, bool) {
	item, ok := this.container.Load(k)
	if ok {
		expItem := item.(expiredItem)
		if !expItem.HasExpire() {
			return expItem.Item, true
		} else {
			this.container.Delete(k)
		}
	} else {
		//key stay in expire keys and no expire
		keyItem, exists := this.expireKeyMap.Load(k)
		if exists {
			expKey := keyItem.(expiredItem)
			if !expKey.HasExpire() {
				return nil, false
			} else {
				this.expireKeyMap.Delete(k)
			}
		}
	}

	data, err := this.dataFetcher(k)
	if nil != err {
		logger.Error("fetch data for %v in cache load fail, err: %v", k, err)
		this.expireKeyMap.Store(k, newExpiredItem(k, this.expireDuration))
		return nil, false
	}
	expItem := newExpiredItem(data, this.expireDuration)
	this.container.Store(k, expItem)

	return data, true
}

func (this *TimeExpireCacheLoader) checkExpire() {
	if this.clearInterval <= 0 {
		return
	}

	ticker := time.NewTicker(this.clearInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			this.expireValues()
			this.expireKeys()
		case <-this.closed:
			break
		}
	}
}

func (this *TimeExpireCacheLoader) expireValues() {
	this.container.Range(func(key, val any) bool {
		item := val.(expiredItem)
		if item.HasExpire() {
			this.container.Delete(key)
		}
		return true
	})
}

func (this *TimeExpireCacheLoader) expireKeys() {
	this.expireKeyMap.Range(func(key, val any) bool {
		item := val.(expiredItem)
		if item.HasExpire() {
			this.expireKeyMap.Delete(key)
		}
		return true
	})
}

type expiredItem struct {
	Item       any
	ExpireTime time.Time
}

func newExpiredItem(val any, expire time.Duration) expiredItem {
	return expiredItem{
		Item:       val,
		ExpireTime: time.Now().Add(expire),
	}
}

func (this *expiredItem) HasExpire() bool {
	return time.Now().After(this.ExpireTime)
}

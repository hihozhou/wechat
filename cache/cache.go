package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	goCache *cache.Cache
)

func init() {
	goCache = cache.New(5*time.Minute, 10*time.Minute)
}

func Set(k string, x interface{}, d time.Duration) {
	goCache.Set(k, x, d)
}

func Get(k string) (interface{}, bool) {
	return goCache.Get(k)
}

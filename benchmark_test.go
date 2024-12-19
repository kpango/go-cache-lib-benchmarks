package main

import (
	"math/rand"
	"sync"
	"testing"
	"time"
	"unsafe"

	// "git.mills.io/prologic/bitcask"
	// "github.com/VictoriaMetrics/fastcache"
	"github.com/bluele/gcache"
	// "github.com/coocood/freecache"
	gachev2 "github.com/kpango/gache/v2"
	"github.com/kpango/gache"
	bigcache "github.com/allegro/bigcache/v3"
	gocache "github.com/patrickmn/go-cache"
	mcache "github.com/OrlovEvgeny/go-mcache"
	ttlcache "github.com/jellydator/ttlcache/v3"
)

type DefaultMap struct {
	mu   sync.RWMutex
	data map[interface{}]interface{}
}

func NewDefault() *DefaultMap {
	return &DefaultMap{
		data: make(map[interface{}]interface{}),
	}
}

func (m *DefaultMap) Get(key interface{}) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}

func (m *DefaultMap) Set(key, val interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
}

const (
	NoTTL = time.Duration(-1)
)

var (
	ttl time.Duration = 50 * time.Millisecond

	parallelism = 10000

	bigData    = map[string]string{}
	bigDataLen = 2 << 10

	smallData = map[string]string{
		"string": "aaaa",
		"int":    "123",
		"float":  "99.99",
		"struct": "struct{}{}",
	}
)

func init() {
	for i := 0; i < bigDataLen; i++ {
		bigData[randStr(bigDataLen)] = randStr(bigDataLen)
	}
}

var randSrc = rand.NewSource(time.Now().UnixNano())

const (
	rs6Letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rs6LetterIdxBits = 6
	rs6LetterIdxMask = 1<<rs6LetterIdxBits - 1
	rs6LetterIdxMax  = 63 / rs6LetterIdxBits
)

func randStr(n int) string {
	b := make([]byte, n)
	cache, remain := randSrc.Int63(), rs6LetterIdxMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), rs6LetterIdxMax
		}
		idx := int(cache & rs6LetterIdxMask)
		if idx < len(rs6Letters) {
			b[i] = rs6Letters[idx]
			i--
		}
		cache >>= rs6LetterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func benchmark(b *testing.B, data map[string]string,
	t time.Duration,
	set func(string, string, time.Duration),
	get func(string)) {
	b.Helper()
	b.SetParallelism(parallelism)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for k, v := range data {
				set(k, v, t)
			}
			for k := range data {
				get(k)
			}
		}
	})

}

func BenchmarkDefaultMapSetSmallDataNoTTL(b *testing.B) {
	m := NewDefault()
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { m.Set(k, v) },
		func(k string) { m.Get(k) })
}
func BenchmarkDefaultMapSetBigDataNoTTL(b *testing.B) {
	m := NewDefault()
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { m.Set(k, v) },
		func(k string) { m.Get(k) })
}
func BenchmarkSyncMapSetSmallDataNoTTL(b *testing.B) {
	var m sync.Map
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { m.Store(k, v) },
		func(k string) { m.Load(k) })
}
func BenchmarkSyncMapSetBigDataNoTTL(b *testing.B) {
	var m sync.Map
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { m.Store(k, v) },
		func(k string) { m.Load(k) })
}
func BenchmarkGacheV2SetSmallDataNoTTL(b *testing.B) {
	g := gachev2.New[string]().SetDefaultExpire(NoTTL)
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}
func BenchmarkGacheV2SetSmallDataWithTTL(b *testing.B) {
	g := gachev2.New[string]().SetDefaultExpire(ttl)
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { g.SetWithExpire(k, v, t) },
		func(k string) { g.Get(k) })
}
func BenchmarkGacheV2SetBigDataNoTTL(b *testing.B) {
	g := gachev2.New[string]().SetDefaultExpire(NoTTL)
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}
func BenchmarkGacheV2SetBigDataWithTTL(b *testing.B) {
	g := gachev2.New[string]().SetDefaultExpire(ttl)
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { g.SetWithExpire(k, v, t) },
		func(k string) { g.Get(k) })
}
func BenchmarkGacheSetSmallDataNoTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(NoTTL)
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}
func BenchmarkGacheSetSmallDataWithTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(ttl)
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { g.SetWithExpire(k, v, t) },
		func(k string) { g.Get(k) })
}
func BenchmarkGacheSetBigDataNoTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(NoTTL)
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}
func BenchmarkGacheSetBigDataWithTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(ttl)
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { g.SetWithExpire(k, v, t) },
		func(k string) { g.Get(k) })
}
func BenchmarkTTLCacheSetSmallDataNoTTL(b *testing.B) {
	cache := ttlcache.New[string, string]()
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { cache.Set(k, v, ttlcache.NoTTL) },
		func(k string) { cache.Get(k) })
}
func BenchmarkTTLCacheSetSmallDataWithTTL(b *testing.B) {
	cache := ttlcache.New[string, string](
		ttlcache.WithTTL[string, string](ttl),
	)
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { cache.Set(k, v, t) },
		func(k string) { cache.Get(k) })
}
func BenchmarkTTLCacheSetBigDataNoTTL(b *testing.B) {
	cache := ttlcache.New[string, string]()
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { cache.Set(k, v, ttlcache.NoTTL) },
		func(k string) { cache.Get(k) })
}
func BenchmarkTTLCacheSetBigDataWithTTL(b *testing.B) {
	cache := ttlcache.New[string, string](
		ttlcache.WithTTL[string, string](ttl),
	)
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { cache.Set(k, v, t) },
		func(k string) { cache.Get(k) })
}
func BenchmarkGoCacheSetSmallDataNoTTL(b *testing.B) {
	c := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkGoCacheSetSmallDataWithTTL(b *testing.B) {
	c := gocache.New(ttl, ttl)
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkGoCacheSetBigDataNoTTL(b *testing.B) {
	c := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkGoCacheSetBigDataWithTTL(b *testing.B) {
	c := gocache.New(ttl, ttl)
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkBigCacheSetSmallDataNoTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(NoTTL)
	cfg.Verbose = false
	bc, _ := bigcache.NewBigCache(cfg)
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}
func BenchmarkBigCacheSetSmallDataWithTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(ttl)
	cfg.Verbose = false
	bc, _ := bigcache.NewBigCache(cfg)
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}
func BenchmarkBigCacheSetBigDataNoTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(NoTTL)
	cfg.Verbose = false
	bc, _ := bigcache.NewBigCache(cfg)
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}
func BenchmarkBigCacheSetBigDataWithTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(ttl)
	cfg.Verbose = false
	bc, _ := bigcache.NewBigCache(cfg)
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}

// func BenchmarkFastCacheSetSmallDataNoTTL(b *testing.B) {
// 	fc := fastcache.New(20)
// 	benchmark(b, smallData, NoTTL,
// 		func(k, v string, t time.Duration) { fc.Set([]byte(k), []byte(v)) },
// 		func(k string) {
// 			var val []byte
// 			val = fc.Get(val, []byte(k))
// 		})
// }
// func BenchmarkFastCacheSetBigDataNoTTL(b *testing.B) {
// 	fc := fastcache.New(20)
// 	benchmark(b, bigData, NoTTL,
// 		func(k, v string, t time.Duration) { fc.SetBig([]byte(k), []byte(v)) },
// 		func(k string) {
// 			var val []byte
// 			val = fc.Get(val, []byte(k))
// 		})
// }
// func BenchmarkFreeCacheSetSmallDataNoTTL(b *testing.B) {
// 	c := freecache.NewCache(100 * 1024 * 1024)
// 	benchmark(b, smallData, NoTTL,
// 		func(k, v string, t time.Duration) { c.Set([]byte(k), []byte(v), -1) },
// 		func(k string) { c.Get([]byte(k)) })
// }
// func BenchmarkFreeCacheSetSmallDataWithTTL(b *testing.B) {
// 	c := freecache.NewCache(100 * 1024 * 1024)
// 	benchmark(b, smallData, ttl,
// 		func(k, v string, t time.Duration) { c.Set([]byte(k), []byte(v), 1) },
// 		func(k string) { c.Get([]byte(k)) })
// }
// func BenchmarkFreeCacheSetBigDataNoTTL(b *testing.B) {
// 	c := freecache.NewCache(100 * 1024 * 1024)
// 	benchmark(b, bigData, NoTTL,
// 		func(k, v string, t time.Duration) { c.Set([]byte(k), []byte(v), -1) },
// 		func(k string) { c.Get([]byte(k)) })
// }
// func BenchmarkFreeCacheSetBigDataWithTTL(b *testing.B) {
// 	c := freecache.NewCache(100 * 1024 * 1024)
// 	benchmark(b, bigData, ttl,
// 		func(k, v string, t time.Duration) { c.Set([]byte(k), []byte(v), 1) },
// 		func(k string) { c.Get([]byte(k)) })
// }
func BenchmarkGCacheLRUSetSmallDataNoTTL(b *testing.B) {
	c := gcache.New(20).LRU().Build()
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkGCacheLRUSetSmallDataWithTTL(b *testing.B) {
	c := gcache.New(20).LRU().Build()
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkGCacheLRUSetBigDataNoTTL(b *testing.B) {
	c := gcache.New(20).LRU().Build()
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkGCacheLRUSetBigDataWithTTL(b *testing.B) {
	c := gcache.New(20).LRU().Build()
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkGCacheLFUSetSmallDataNoTTL(b *testing.B) {
	c := gcache.New(20).LFU().Build()
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkGCacheLFUSetSmallDataWithTTL(b *testing.B) {
	c := gcache.New(20).LFU().Build()
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkGCacheLFUSetBigDataNoTTL(b *testing.B) {
	c := gcache.New(20).LFU().Build()
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkGCacheLFUSetBigDataWithTTL(b *testing.B) {
	c := gcache.New(20).LFU().Build()
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkGCacheARCSetSmallDataNoTTL(b *testing.B) {
	c := gcache.New(20).ARC().Build()
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkGCacheARCSetSmallDataWithTTL(b *testing.B) {
	c := gcache.New(20).ARC().Build()
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkGCacheARCSetBigDataNoTTL(b *testing.B) {
	c := gcache.New(20).ARC().Build()
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkGCacheARCSetBigDataWithTTL(b *testing.B) {
	c := gcache.New(20).ARC().Build()
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}
func BenchmarkMCacheSetSmallDataNoTTL(b *testing.B) {
	c := mcache.StartInstance()
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.SetPointer(k, v, t) },
		func(k string) { c.GetPointer(k) })
}
func BenchmarkMCacheSetSmallDataWithTTL(b *testing.B) {
	c := mcache.StartInstance()
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { c.SetPointer(k, v, t) },
		func(k string) { c.GetPointer(k) })
}
func BenchmarkMCacheSetBigDataNoTTL(b *testing.B) {
	c := mcache.StartInstance()
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.SetPointer(k, v, t) },
		func(k string) { c.GetPointer(k) })
}
func BenchmarkMCacheSetBigDataWithTTL(b *testing.B) {
	c := mcache.StartInstance()
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { c.SetPointer(k, v, t) },
		func(k string) { c.GetPointer(k) })
}
// func BenchmarkBitcaskSetSmallDataNoTTL(b *testing.B) {
// 	bc, _ := bitcask.Open("/tmp/db")
// 	benchmark(b, smallData, NoTTL,
// 		func(k, v string, t time.Duration) { bc.Put([]byte(k), []byte(v)) },
// 		func(k string) { bc.Get([]byte(k)) })
// }
// func BenchmarkBitcaskSetSmallDataWithTTL(b *testing.B) {
// 	bc, _ := bitcask.Open("/tmp/db")
// 	benchmark(b, smallData, ttl,
// 		func(k, v string, t time.Duration) { bc.Put([]byte(k), []byte(v)) },
// 		func(k string) { bc.Get([]byte(k)) })
// }
// func BenchmarkBitcaskSetBigDataNoTTL(b *testing.B) {
// 	bc, _ := bitcask.Open("/tmp/db")
// 	benchmark(b, bigData, NoTTL,
// 		func(k, v string, t time.Duration) { bc.Put([]byte(k), []byte(v)) },
// 		func(k string) { bc.Get([]byte(k)) })
// }
// func BenchmarkBitcaskSetBigDataWithTTL(b *testing.B) {
// 	bc, _ := bitcask.Open("/tmp/db")
// 	benchmark(b, bigData, ttl,
// 		func(k, v string, t time.Duration) { bc.Put([]byte(k), []byte(v)) },
// 		func(k string) { bc.Get([]byte(k)) })
// }

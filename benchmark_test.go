package main

import (
	// "git.mills.io/prologic/bitcask"
	// "github.com/VictoriaMetrics/fastcache"
	// "github.com/coocood/freecache"
	// mcache "github.com/OrlovEvgeny/go-mcache"
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	bigcache "github.com/allegro/bigcache/v3"
	"github.com/bluele/gcache"
	ttlcache "github.com/jellydator/ttlcache/v3"
	"github.com/kpango/gache"
	gachev2 "github.com/kpango/gache/v2"
	gocache "github.com/patrickmn/go-cache"
)

type DefaultMap struct {
	mu   sync.RWMutex
	data map[any]any
}

func NewDefault() *DefaultMap {
	return &DefaultMap{
		data: make(map[any]any),
	}
}

func (m *DefaultMap) Get(key any) (any, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}

func (m *DefaultMap) Set(key, val any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val
}

const NoTTL = time.Duration(-1)

// bigCacheNoTTL uses a large positive duration for BigCache since it does not
// handle negative TTL values correctly.
const bigCacheNoTTL = 24 * time.Hour

var benchParallelismFlag string

var parallelismValues []int

type keyValue struct {
	key   string
	value string
}

var (
	ttl time.Duration = 50 * time.Millisecond

	smallData []keyValue
	bigData   []keyValue
)

func init() {
	flag.StringVar(&benchParallelismFlag, "benchparallelism", "", "comma-separated list of parallelism values for benchmarks (default: 100,1000,5000,10000)")

	var (
		bigDataLen     = 2 << 10
		bigDataCount   = 2 << 16
		smallDataLen   = 2 << 5
		smallDataCount = 2 << 3
	)
	bigData = make([]keyValue, 0, bigDataCount)
	for range bigDataCount {
		bigData = append(bigData, keyValue{
			key:   randStr(bigDataLen),
			value: randStr(bigDataLen),
		})
	}
	slices.SortFunc(bigData, func(a, b keyValue) int {
		return strings.Compare(a.key, b.key)
	})
	smallData = make([]keyValue, 0, smallDataCount)
	for range smallDataCount {
		smallData = append(smallData, keyValue{
			key:   randStr(smallDataLen),
			value: randStr(smallDataLen),
		})
	}
	slices.SortFunc(smallData, func(a, b keyValue) int {
		return strings.Compare(a.key, b.key)
	})

	if benchParallelismFlag != "" {
		for s := range strings.SplitSeq(benchParallelismFlag, ",") {
			v, err := strconv.Atoi(strings.TrimSpace(s))
			if err == nil && v > 0 {
				parallelismValues = append(parallelismValues, v)
			}
		}
	}
	if len(parallelismValues) == 0 {
		parallelismValues = []int{100, 1000, 10000}
	}
}

var randSrc = rand.NewSource(42)

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
	return string(b)
}

func benchmark(b *testing.B, data []keyValue,
	t time.Duration,
	set func(string, string, time.Duration),
	get func(string),
) {
	b.Helper()
	nprocs := runtime.GOMAXPROCS(0)
	for _, p := range parallelismValues {
		b.Run(fmt.Sprintf("P%d", p), func(b *testing.B) {
			mp := max(p/nprocs, 1)
			b.SetParallelism(mp)
			runBenchParallel(b, func(_ *testing.PB, _ int) {
				for _, kv := range data {
					set(kv.key, kv.value, t)
				}
				for _, kv := range data {
					get(kv.key)
				}
			})
		})
	}
}

func runBenchParallel(b *testing.B, f func(pb *testing.PB, i int)) {
	b.Helper()
	runtime.GC()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var i int
		for pb.Next() {
			f(pb, i)
			i++
		}
	})
}

// benchmarkSetOnly runs a write-only workload benchmark for each configured
// parallelism value, emitting sub-benchmarks named "P<n>".
func benchmarkSetOnly(b *testing.B, data []keyValue,
	t time.Duration,
	set func(string, string, time.Duration),
) {
	b.Helper()
	for _, p := range parallelismValues {
		b.Run(fmt.Sprintf("P%d", p), func(b *testing.B) {
			b.SetParallelism(p)
			b.ReportAllocs()
			runtime.GC()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					for _, kv := range data {
						set(kv.key, kv.value, t)
					}
				}
			})
		})
	}
}

// benchmarkGetOnly runs a read-only workload benchmark. The cache is
// pre-populated before timing begins to isolate read performance.
// Each configured parallelism value produces a sub-benchmark named "P<n>".
func benchmarkGetOnly(b *testing.B, data []keyValue,
	t time.Duration,
	setup func(string, string, time.Duration),
	get func(string),
) {
	b.Helper()
	for _, kv := range data {
		setup(kv.key, kv.value, t)
	}
	for _, p := range parallelismValues {
		b.Run(fmt.Sprintf("P%d", p), func(b *testing.B) {
			b.SetParallelism(p)
			b.ReportAllocs()
			runtime.GC()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					for _, kv := range data {
						get(kv.key)
					}
				}
			})
		})
	}
}

func BenchmarkDefaultMapSetGetSmallDataNoTTL(b *testing.B) {
	m := NewDefault()
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { m.Set(k, v) },
		func(k string) { m.Get(k) })
}

func BenchmarkDefaultMapSetGetBigDataNoTTL(b *testing.B) {
	m := NewDefault()
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { m.Set(k, v) },
		func(k string) { m.Get(k) })
}

func BenchmarkSyncMapSetGetSmallDataNoTTL(b *testing.B) {
	var m sync.Map
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { m.Store(k, v) },
		func(k string) { m.Load(k) })
}

func BenchmarkSyncMapSetGetBigDataNoTTL(b *testing.B) {
	var m sync.Map
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { m.Store(k, v) },
		func(k string) { m.Load(k) })
}

func BenchmarkGacheV2SetGetSmallDataNoTTL(b *testing.B) {
	g := gachev2.New(gachev2.WithDefaultExpiration[string](gachev2.NoTTL))
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheV2SetGetSmallDataWithTTL(b *testing.B) {
	g := gachev2.New(gachev2.WithDefaultExpiration[string](ttl))
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { g.SetWithExpire(k, v, t) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheV2SetGetBigDataNoTTL(b *testing.B) {
	g := gachev2.New(
		gachev2.WithDefaultExpiration[string](gachev2.NoTTL),
		gachev2.WithMaxKeyLength[string](0))
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheV2SetGetBigDataWithTTL(b *testing.B) {
	g := gachev2.New(
		gachev2.WithDefaultExpiration[string](ttl),
		gachev2.WithMaxKeyLength[string](0))
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { g.SetWithExpire(k, v, t) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheSetGetSmallDataNoTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(NoTTL)
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheSetGetSmallDataWithTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(ttl)
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { g.SetWithExpire(k, v, t) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheSetGetBigDataNoTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(NoTTL)
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheSetGetBigDataWithTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(ttl)
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { g.SetWithExpire(k, v, t) },
		func(k string) { g.Get(k) })
}

func BenchmarkTTLCacheSetGetSmallDataNoTTL(b *testing.B) {
	cache := ttlcache.New[string, string]()
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { cache.Set(k, v, ttlcache.NoTTL) },
		func(k string) { cache.Get(k) })
}

func BenchmarkTTLCacheSetGetSmallDataWithTTL(b *testing.B) {
	cache := ttlcache.New(
		ttlcache.WithTTL[string, string](ttl),
	)
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { cache.Set(k, v, t) },
		func(k string) { cache.Get(k) })
}

func BenchmarkTTLCacheSetGetBigDataNoTTL(b *testing.B) {
	cache := ttlcache.New[string, string]()
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { cache.Set(k, v, ttlcache.NoTTL) },
		func(k string) { cache.Get(k) })
}

func BenchmarkTTLCacheSetGetBigDataWithTTL(b *testing.B) {
	cache := ttlcache.New(
		ttlcache.WithTTL[string, string](ttl),
	)
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { cache.Set(k, v, t) },
		func(k string) { cache.Get(k) })
}

func BenchmarkGoCacheSetGetSmallDataNoTTL(b *testing.B) {
	c := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGoCacheSetGetSmallDataWithTTL(b *testing.B) {
	c := gocache.New(ttl, ttl)
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGoCacheSetGetBigDataNoTTL(b *testing.B) {
	c := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGoCacheSetGetBigDataWithTTL(b *testing.B) {
	c := gocache.New(ttl, ttl)
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkBigCacheSetGetSmallDataNoTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(bigCacheNoTTL)
	cfg.Verbose = false
	bc, _ := bigcache.New(b.Context(), cfg)
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}

func BenchmarkBigCacheSetGetSmallDataWithTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(ttl)
	cfg.Verbose = false
	bc, _ := bigcache.New(b.Context(), cfg)
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}

func BenchmarkBigCacheSetGetBigDataNoTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(bigCacheNoTTL)
	cfg.Verbose = false
	bc, _ := bigcache.New(b.Context(), cfg)
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}

func BenchmarkBigCacheSetGetBigDataWithTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(ttl)
	cfg.Verbose = false
	bc, _ := bigcache.New(b.Context(), cfg)
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}

//	func BenchmarkFastCacheSetGetSmallDataNoTTL(b *testing.B) {
//		fc := fastcache.New(20)
//		benchmark(b, smallData, NoTTL,
//			func(k, v string, t time.Duration) { fc.Set([]byte(k), []byte(v)) },
//			func(k string) {
//				var val []byte
//				val = fc.Get(val, []byte(k))
//			})
//	}
//
//	func BenchmarkFastCacheSetGetBigDataNoTTL(b *testing.B) {
//		fc := fastcache.New(20)
//		benchmark(b, bigData, NoTTL,
//			func(k, v string, t time.Duration) { fc.SetGetBig([]byte(k), []byte(v)) },
//			func(k string) {
//				var val []byte
//				val = fc.Get(val, []byte(k))
//			})
//	}
//
//	func BenchmarkFreeCacheSetGetSmallDataNoTTL(b *testing.B) {
//		c := freecache.NewCache(100 * 1024 * 1024)
//		benchmark(b, smallData, NoTTL,
//			func(k, v string, t time.Duration) { c.Set([]byte(k), []byte(v), -1) },
//			func(k string) { c.Get([]byte(k)) })
//	}
//
//	func BenchmarkFreeCacheSetGetSmallDataWithTTL(b *testing.B) {
//		c := freecache.NewCache(100 * 1024 * 1024)
//		benchmark(b, smallData, ttl,
//			func(k, v string, t time.Duration) { c.Set([]byte(k), []byte(v), 1) },
//			func(k string) { c.Get([]byte(k)) })
//	}
//
//	func BenchmarkFreeCacheSetGetBigDataNoTTL(b *testing.B) {
//		c := freecache.NewCache(100 * 1024 * 1024)
//		benchmark(b, bigData, NoTTL,
//			func(k, v string, t time.Duration) { c.Set([]byte(k), []byte(v), -1) },
//			func(k string) { c.Get([]byte(k)) })
//	}
//
//	func BenchmarkFreeCacheSetGetBigDataWithTTL(b *testing.B) {
//		c := freecache.NewCache(100 * 1024 * 1024)
//		benchmark(b, bigData, ttl,
//			func(k, v string, t time.Duration) { c.Set([]byte(k), []byte(v), 1) },
//			func(k string) { c.Get([]byte(k)) })
//	}
func BenchmarkGCacheLRUSetGetSmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LRU().Build()
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLRUSetGetSmallDataWithTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LRU().Build()
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLRUSetGetBigDataNoTTL(b *testing.B) {
	c := gcache.New(len(bigData)).LRU().Build()
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLRUSetGetBigDataWithTTL(b *testing.B) {
	c := gcache.New(len(bigData)).LRU().Build()
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLFUSetGetSmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LFU().Build()
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLFUSetGetSmallDataWithTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LFU().Build()
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLFUSetGetBigDataNoTTL(b *testing.B) {
	c := gcache.New(len(bigData)).LFU().Build()
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLFUSetGetBigDataWithTTL(b *testing.B) {
	c := gcache.New(len(bigData)).LFU().Build()
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheARCSetGetSmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).ARC().Build()
	benchmark(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheARCSetGetSmallDataWithTTL(b *testing.B) {
	c := gcache.New(len(smallData)).ARC().Build()
	benchmark(b, smallData, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheARCSetGetBigDataNoTTL(b *testing.B) {
	c := gcache.New(len(bigData)).ARC().Build()
	benchmark(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheARCSetGetBigDataWithTTL(b *testing.B) {
	c := gcache.New(len(bigData)).ARC().Build()
	benchmark(b, bigData, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

// ----- Set-only benchmarks -----

func BenchmarkDefaultMapSetOnlySmallDataNoTTL(b *testing.B) {
	m := NewDefault()
	benchmarkSetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { m.Set(k, v) })
}

func BenchmarkDefaultMapSetOnlyBigDataNoTTL(b *testing.B) {
	m := NewDefault()
	benchmarkSetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { m.Set(k, v) })
}

func BenchmarkSyncMapSetOnlySmallDataNoTTL(b *testing.B) {
	var m sync.Map
	benchmarkSetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { m.Store(k, v) })
}

func BenchmarkSyncMapSetOnlyBigDataNoTTL(b *testing.B) {
	var m sync.Map
	benchmarkSetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { m.Store(k, v) })
}

func BenchmarkGacheV2SetOnlySmallDataNoTTL(b *testing.B) {
	g := gachev2.New(gachev2.WithDefaultExpiration[string](gachev2.NoTTL))
	benchmarkSetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) })
}

func BenchmarkGacheV2SetOnlyBigDataNoTTL(b *testing.B) {
	g := gachev2.New(
		gachev2.WithDefaultExpiration[string](gachev2.NoTTL),
		gachev2.WithMaxKeyLength[string](0))
	benchmarkSetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) })
}

func BenchmarkGacheSetOnlySmallDataNoTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(NoTTL)
	benchmarkSetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) })
}

func BenchmarkGacheSetOnlyBigDataNoTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(NoTTL)
	benchmarkSetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) })
}

func BenchmarkTTLCacheSetOnlySmallDataNoTTL(b *testing.B) {
	cache := ttlcache.New[string, string]()
	benchmarkSetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { cache.Set(k, v, ttlcache.NoTTL) })
}

func BenchmarkTTLCacheSetOnlyBigDataNoTTL(b *testing.B) {
	cache := ttlcache.New[string, string]()
	benchmarkSetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { cache.Set(k, v, ttlcache.NoTTL) })
}

func BenchmarkGoCacheSetOnlySmallDataNoTTL(b *testing.B) {
	c := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	benchmarkSetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.Set(k, v, t) })
}

func BenchmarkGoCacheSetOnlyBigDataNoTTL(b *testing.B) {
	c := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	benchmarkSetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.Set(k, v, t) })
}

func BenchmarkBigCacheSetOnlySmallDataNoTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(bigCacheNoTTL)
	cfg.Verbose = false
	bc, _ := bigcache.New(b.Context(), cfg)
	benchmarkSetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) })
}

func BenchmarkBigCacheSetOnlyBigDataNoTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(bigCacheNoTTL)
	cfg.Verbose = false
	bc, _ := bigcache.New(b.Context(), cfg)
	benchmarkSetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) })
}

func BenchmarkGCacheLRUSetOnlySmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LRU().Build()
	benchmarkSetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) })
}

func BenchmarkGCacheLRUSetOnlyBigDataNoTTL(b *testing.B) {
	c := gcache.New(len(bigData)).LRU().Build()
	benchmarkSetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) })
}

func BenchmarkGCacheLFUSetOnlySmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LFU().Build()
	benchmarkSetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) })
}

func BenchmarkGCacheLFUSetOnlyBigDataNoTTL(b *testing.B) {
	c := gcache.New(len(bigData)).LRU().Build()
	benchmarkSetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) })
}

func BenchmarkGCacheARCSetOnlySmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).ARC().Build()
	benchmarkSetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) })
}

func BenchmarkGCacheARCSetOnlyBigDataNoTTL(b *testing.B) {
	c := gcache.New(len(bigData)).ARC().Build()
	benchmarkSetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) })
}

// ----- Get-only benchmarks (cache pre-populated before timing) -----

func BenchmarkDefaultMapGetSmallDataNoTTL(b *testing.B) {
	m := NewDefault()
	benchmarkGetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { m.Set(k, v) },
		func(k string) { m.Get(k) })
}

func BenchmarkDefaultMapGetBigDataNoTTL(b *testing.B) {
	m := NewDefault()
	benchmarkGetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { m.Set(k, v) },
		func(k string) { m.Get(k) })
}

func BenchmarkSyncMapGetSmallDataNoTTL(b *testing.B) {
	var m sync.Map
	benchmarkGetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { m.Store(k, v) },
		func(k string) { m.Load(k) })
}

func BenchmarkSyncMapGetBigDataNoTTL(b *testing.B) {
	var m sync.Map
	benchmarkGetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { m.Store(k, v) },
		func(k string) { m.Load(k) })
}

func BenchmarkGacheV2GetSmallDataNoTTL(b *testing.B) {
	g := gachev2.New(gachev2.WithDefaultExpiration[string](gachev2.NoTTL))
	benchmarkGetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheV2GetBigDataNoTTL(b *testing.B) {
	g := gachev2.New(
		gachev2.WithDefaultExpiration[string](gachev2.NoTTL),
		gachev2.WithMaxKeyLength[string](0))
	benchmarkGetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheGetSmallDataNoTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(NoTTL)
	benchmarkGetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheGetBigDataNoTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(NoTTL)
	benchmarkGetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkTTLCacheGetSmallDataNoTTL(b *testing.B) {
	cache := ttlcache.New[string, string]()
	benchmarkGetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { cache.Set(k, v, ttlcache.NoTTL) },
		func(k string) { cache.Get(k) })
}

func BenchmarkTTLCacheGetBigDataNoTTL(b *testing.B) {
	cache := ttlcache.New[string, string]()
	benchmarkGetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { cache.Set(k, v, ttlcache.NoTTL) },
		func(k string) { cache.Get(k) })
}

func BenchmarkGoCacheGetSmallDataNoTTL(b *testing.B) {
	c := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	benchmarkGetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGoCacheGetBigDataNoTTL(b *testing.B) {
	c := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	benchmarkGetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkBigCacheGetSmallDataNoTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(bigCacheNoTTL)
	cfg.Verbose = false
	bc, _ := bigcache.New(b.Context(), cfg)
	benchmarkGetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}

func BenchmarkBigCacheGetBigDataNoTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(bigCacheNoTTL)
	cfg.Verbose = false
	bc, _ := bigcache.New(b.Context(), cfg)
	benchmarkGetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}

func BenchmarkGCacheLRUGetSmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LRU().Build()
	benchmarkGetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLRUGetBigDataNoTTL(b *testing.B) {
	c := gcache.New(len(bigData)).LRU().Build()
	benchmarkGetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLFUGetSmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LFU().Build()
	benchmarkGetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLFUGetBigDataNoTTL(b *testing.B) {
	c := gcache.New(len(bigData)).LFU().Build()
	benchmarkGetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheARCGetSmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).ARC().Build()
	benchmarkGetOnly(b, smallData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheARCGetBigDataNoTTL(b *testing.B) {
	c := gcache.New(len(bigData)).ARC().Build()
	benchmarkGetOnly(b, bigData, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

// func BenchmarkMCacheSetGetSmallDataNoTTL(b *testing.B) {
// 	c := mcache.StartInstance()
// 	benchmark(b, smallData, NoTTL,
// 		func(k, v string, t time.Duration) { c.SetPointer(k, v, t) },
// 		func(k string) { c.GetPointer(k) })
// }

// func BenchmarkMCacheSetGetSmallDataWithTTL(b *testing.B) {
// 	c := mcache.StartInstance()
// 	benchmark(b, smallData, ttl,
// 		func(k, v string, t time.Duration) { c.SetPointer(k, v, t) },
// 		func(k string) { c.GetPointer(k) })
// }

// func BenchmarkMCacheSetGetBigDataNoTTL(b *testing.B) {
// 	c := mcache.StartInstance()
// 	benchmark(b, bigData, NoTTL,
// 		func(k, v string, t time.Duration) { c.SetPointer(k, v, t) },
// 		func(k string) { c.GetPointer(k) })
// }

// func BenchmarkMCacheSetGetBigDataWithTTL(b *testing.B) {
// 	c := mcache.StartInstance()
// 	benchmark(b, bigData, ttl,
// 		func(k, v string, t time.Duration) { c.SetPointer(k, v, t) },
// 		func(k string) { c.GetPointer(k) })
// }

// func BenchmarkBitcaskSetGetSmallDataNoTTL(b *testing.B) {
// 	bc, _ := bitcask.Open("/tmp/db")
// 	benchmark(b, smallData, NoTTL,
// 		func(k, v string, t time.Duration) { bc.Put([]byte(k), []byte(v)) },
// 		func(k string) { bc.Get([]byte(k)) })
// }
// func BenchmarkBitcaskSetGetSmallDataWithTTL(b *testing.B) {
// 	bc, _ := bitcask.Open("/tmp/db")
// 	benchmark(b, smallData, ttl,
// 		func(k, v string, t time.Duration) { bc.Put([]byte(k), []byte(v)) },
// 		func(k string) { bc.Get([]byte(k)) })
// }
// func BenchmarkBitcaskSetGetBigDataNoTTL(b *testing.B) {
// 	bc, _ := bitcask.Open("/tmp/db")
// 	benchmark(b, bigData, NoTTL,
// 		func(k, v string, t time.Duration) { bc.Put([]byte(k), []byte(v)) },
// 		func(k string) { bc.Get([]byte(k)) })
// }
// func BenchmarkBitcaskSetGetBigDataWithTTL(b *testing.B) {
// 	bc, _ := bitcask.Open("/tmp/db")
// 	benchmark(b, bigData, ttl,
// 		func(k, v string, t time.Duration) { bc.Put([]byte(k), []byte(v)) },
// 		func(k string) { bc.Get([]byte(k)) })
// }

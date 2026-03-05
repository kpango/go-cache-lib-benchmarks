package main

import (
	// "git.mills.io/prologic/bitcask"
	// "github.com/VictoriaMetrics/fastcache"
	// "github.com/coocood/freecache"
	// mcache "github.com/OrlovEvgeny/go-mcache"
	"flag"
	"fmt"
	"math/rand"
	"os"
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

// keyValue holds a pre-computed key-value pair for deterministic benchmark iteration.
type keyValue struct {
	key   string
	value string
}

// mapToSlice converts a map to a sorted slice of keyValue pairs.
// Sorting ensures deterministic iteration order across benchmark runs.
func mapToSlice(m map[string]string) []keyValue {
	kvs := make([]keyValue, 0, len(m))
	for k, v := range m {
		kvs = append(kvs, keyValue{key: k, value: v})
	}
	slices.SortFunc(kvs, func(a, b keyValue) int {
		return strings.Compare(a.key, b.key)
	})
	return kvs
}

const NoTTL = time.Duration(-1)

// bigCacheNoTTL uses a large positive duration for BigCache since it does not
// handle negative TTL values correctly.
const bigCacheNoTTL = 24 * time.Hour

// benchParallelismFlag holds the raw flag value for -benchparallelism.
var benchParallelismFlag string

// parallelismValues is the set of parallelism levels used by all benchmarks.
// It is populated from -benchparallelism (comma-separated integers) in
// TestMain; the default is []int{100, 1000, 10000}.
var parallelismValues []int

var (
	ttl time.Duration = 50 * time.Millisecond

	bigDataLen   = 2 << 10
	bigDataCount = 2 << 16
	bigData      = make(map[string]string, bigDataCount)

	smallData = map[string]string{
		"string": "aaaa",
		"int":    "123",
		"float":  "99.99",
		"struct": "struct{}{}",
	}

	// Pre-computed slices for deterministic iteration order in benchmarks.
	smallDataSlice []keyValue
	bigDataSlice   []keyValue
)

func init() {
	flag.StringVar(&benchParallelismFlag, "benchparallelism", "", "comma-separated list of parallelism values for benchmarks (default: 100,1000,10000)")
	for range bigDataCount {
		bigData[randStr(bigDataLen)] = randStr(bigDataLen)
	}
	smallDataSlice = mapToSlice(smallData)
	bigDataSlice = mapToSlice(bigData)
}

func TestMain(m *testing.M) {
	flag.Parse()
	if benchParallelismFlag != "" {
		for _, s := range strings.Split(benchParallelismFlag, ",") {
			v, err := strconv.Atoi(strings.TrimSpace(s))
			if err == nil && v > 0 {
				parallelismValues = append(parallelismValues, v)
			}
		}
	}
	if len(parallelismValues) == 0 {
		parallelismValues = []int{100, 1000, 10000}
	}
	os.Exit(m.Run())
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
	return string(b)
}

// benchmark runs a mixed set-and-get workload benchmark for each configured
// parallelism value, emitting sub-benchmarks named "P<n>".
func benchmark(b *testing.B, data []keyValue,
	t time.Duration,
	set func(string, string, time.Duration),
	get func(string),
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
					for _, kv := range data {
						get(kv.key)
					}
				}
			})
		})
	}
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

func BenchmarkDefaultMapSetSmallDataNoTTL(b *testing.B) {
	m := NewDefault()
	benchmark(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { m.Set(k, v) },
		func(k string) { m.Get(k) })
}

func BenchmarkDefaultMapSetBigDataNoTTL(b *testing.B) {
	m := NewDefault()
	benchmark(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { m.Set(k, v) },
		func(k string) { m.Get(k) })
}

func BenchmarkSyncMapSetSmallDataNoTTL(b *testing.B) {
	var m sync.Map
	benchmark(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { m.Store(k, v) },
		func(k string) { m.Load(k) })
}

func BenchmarkSyncMapSetBigDataNoTTL(b *testing.B) {
	var m sync.Map
	benchmark(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { m.Store(k, v) },
		func(k string) { m.Load(k) })
}

func BenchmarkGacheV2SetSmallDataNoTTL(b *testing.B) {
	g := gachev2.New(gachev2.WithDefaultExpiration[string](gachev2.NoTTL))
	benchmark(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheV2SetSmallDataWithTTL(b *testing.B) {
	g := gachev2.New(gachev2.WithDefaultExpiration[string](ttl))
	benchmark(b, smallDataSlice, ttl,
		func(k, v string, t time.Duration) { g.SetWithExpire(k, v, t) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheV2SetBigDataNoTTL(b *testing.B) {
	g := gachev2.New(
		gachev2.WithDefaultExpiration[string](gachev2.NoTTL),
		gachev2.WithMaxKeyLength[string](0))
	benchmark(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheV2SetBigDataWithTTL(b *testing.B) {
	g := gachev2.New(
		gachev2.WithDefaultExpiration[string](ttl),
		gachev2.WithMaxKeyLength[string](0))
	benchmark(b, bigDataSlice, ttl,
		func(k, v string, t time.Duration) { g.SetWithExpire(k, v, t) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheSetSmallDataNoTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(NoTTL)
	benchmark(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheSetSmallDataWithTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(ttl)
	benchmark(b, smallDataSlice, ttl,
		func(k, v string, t time.Duration) { g.SetWithExpire(k, v, t) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheSetBigDataNoTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(NoTTL)
	benchmark(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheSetBigDataWithTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(ttl)
	benchmark(b, bigDataSlice, ttl,
		func(k, v string, t time.Duration) { g.SetWithExpire(k, v, t) },
		func(k string) { g.Get(k) })
}

func BenchmarkTTLCacheSetSmallDataNoTTL(b *testing.B) {
	cache := ttlcache.New[string, string]()
	benchmark(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { cache.Set(k, v, ttlcache.NoTTL) },
		func(k string) { cache.Get(k) })
}

func BenchmarkTTLCacheSetSmallDataWithTTL(b *testing.B) {
	cache := ttlcache.New[string, string](
		ttlcache.WithTTL[string, string](ttl),
	)
	benchmark(b, smallDataSlice, ttl,
		func(k, v string, t time.Duration) { cache.Set(k, v, t) },
		func(k string) { cache.Get(k) })
}

func BenchmarkTTLCacheSetBigDataNoTTL(b *testing.B) {
	cache := ttlcache.New[string, string]()
	benchmark(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { cache.Set(k, v, ttlcache.NoTTL) },
		func(k string) { cache.Get(k) })
}

func BenchmarkTTLCacheSetBigDataWithTTL(b *testing.B) {
	cache := ttlcache.New[string, string](
		ttlcache.WithTTL[string, string](ttl),
	)
	benchmark(b, bigDataSlice, ttl,
		func(k, v string, t time.Duration) { cache.Set(k, v, t) },
		func(k string) { cache.Get(k) })
}

func BenchmarkGoCacheSetSmallDataNoTTL(b *testing.B) {
	c := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	benchmark(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGoCacheSetSmallDataWithTTL(b *testing.B) {
	c := gocache.New(ttl, ttl)
	benchmark(b, smallDataSlice, ttl,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGoCacheSetBigDataNoTTL(b *testing.B) {
	c := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	benchmark(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGoCacheSetBigDataWithTTL(b *testing.B) {
	c := gocache.New(ttl, ttl)
	benchmark(b, bigDataSlice, ttl,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkBigCacheSetSmallDataNoTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(bigCacheNoTTL)
	cfg.Verbose = false
	bc, _ := bigcache.NewBigCache(cfg)
	benchmark(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}

func BenchmarkBigCacheSetSmallDataWithTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(ttl)
	cfg.Verbose = false
	bc, _ := bigcache.NewBigCache(cfg)
	benchmark(b, smallDataSlice, ttl,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}

func BenchmarkBigCacheSetBigDataNoTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(bigCacheNoTTL)
	cfg.Verbose = false
	bc, _ := bigcache.NewBigCache(cfg)
	benchmark(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}

func BenchmarkBigCacheSetBigDataWithTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(ttl)
	cfg.Verbose = false
	bc, _ := bigcache.NewBigCache(cfg)
	benchmark(b, bigDataSlice, ttl,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}

//	func BenchmarkFastCacheSetSmallDataNoTTL(b *testing.B) {
//		fc := fastcache.New(20)
//		benchmark(b, smallData, NoTTL,
//			func(k, v string, t time.Duration) { fc.Set([]byte(k), []byte(v)) },
//			func(k string) {
//				var val []byte
//				val = fc.Get(val, []byte(k))
//			})
//	}
//
//	func BenchmarkFastCacheSetBigDataNoTTL(b *testing.B) {
//		fc := fastcache.New(20)
//		benchmark(b, bigData, NoTTL,
//			func(k, v string, t time.Duration) { fc.SetBig([]byte(k), []byte(v)) },
//			func(k string) {
//				var val []byte
//				val = fc.Get(val, []byte(k))
//			})
//	}
//
//	func BenchmarkFreeCacheSetSmallDataNoTTL(b *testing.B) {
//		c := freecache.NewCache(100 * 1024 * 1024)
//		benchmark(b, smallData, NoTTL,
//			func(k, v string, t time.Duration) { c.Set([]byte(k), []byte(v), -1) },
//			func(k string) { c.Get([]byte(k)) })
//	}
//
//	func BenchmarkFreeCacheSetSmallDataWithTTL(b *testing.B) {
//		c := freecache.NewCache(100 * 1024 * 1024)
//		benchmark(b, smallData, ttl,
//			func(k, v string, t time.Duration) { c.Set([]byte(k), []byte(v), 1) },
//			func(k string) { c.Get([]byte(k)) })
//	}
//
//	func BenchmarkFreeCacheSetBigDataNoTTL(b *testing.B) {
//		c := freecache.NewCache(100 * 1024 * 1024)
//		benchmark(b, bigData, NoTTL,
//			func(k, v string, t time.Duration) { c.Set([]byte(k), []byte(v), -1) },
//			func(k string) { c.Get([]byte(k)) })
//	}
//
//	func BenchmarkFreeCacheSetBigDataWithTTL(b *testing.B) {
//		c := freecache.NewCache(100 * 1024 * 1024)
//		benchmark(b, bigData, ttl,
//			func(k, v string, t time.Duration) { c.Set([]byte(k), []byte(v), 1) },
//			func(k string) { c.Get([]byte(k)) })
//	}
func BenchmarkGCacheLRUSetSmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LRU().Build()
	benchmark(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLRUSetSmallDataWithTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LRU().Build()
	benchmark(b, smallDataSlice, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLRUSetBigDataNoTTL(b *testing.B) {
	c := gcache.New(bigDataCount).LRU().Build()
	benchmark(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLRUSetBigDataWithTTL(b *testing.B) {
	c := gcache.New(bigDataCount).LRU().Build()
	benchmark(b, bigDataSlice, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLFUSetSmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LFU().Build()
	benchmark(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLFUSetSmallDataWithTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LFU().Build()
	benchmark(b, smallDataSlice, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLFUSetBigDataNoTTL(b *testing.B) {
	c := gcache.New(bigDataCount).LFU().Build()
	benchmark(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLFUSetBigDataWithTTL(b *testing.B) {
	c := gcache.New(bigDataCount).LFU().Build()
	benchmark(b, bigDataSlice, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheARCSetSmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).ARC().Build()
	benchmark(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheARCSetSmallDataWithTTL(b *testing.B) {
	c := gcache.New(len(smallData)).ARC().Build()
	benchmark(b, smallDataSlice, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheARCSetBigDataNoTTL(b *testing.B) {
	c := gcache.New(bigDataCount).ARC().Build()
	benchmark(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheARCSetBigDataWithTTL(b *testing.B) {
	c := gcache.New(bigDataCount).ARC().Build()
	benchmark(b, bigDataSlice, ttl,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

// ----- Set-only benchmarks -----

func BenchmarkDefaultMapSetOnlySmallDataNoTTL(b *testing.B) {
	m := NewDefault()
	benchmarkSetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { m.Set(k, v) })
}

func BenchmarkDefaultMapSetOnlyBigDataNoTTL(b *testing.B) {
	m := NewDefault()
	benchmarkSetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { m.Set(k, v) })
}

func BenchmarkSyncMapSetOnlySmallDataNoTTL(b *testing.B) {
	var m sync.Map
	benchmarkSetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { m.Store(k, v) })
}

func BenchmarkSyncMapSetOnlyBigDataNoTTL(b *testing.B) {
	var m sync.Map
	benchmarkSetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { m.Store(k, v) })
}

func BenchmarkGacheV2SetOnlySmallDataNoTTL(b *testing.B) {
	g := gachev2.New(gachev2.WithDefaultExpiration[string](gachev2.NoTTL))
	benchmarkSetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) })
}

func BenchmarkGacheV2SetOnlyBigDataNoTTL(b *testing.B) {
	g := gachev2.New(
		gachev2.WithDefaultExpiration[string](gachev2.NoTTL),
		gachev2.WithMaxKeyLength[string](0))
	benchmarkSetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) })
}

func BenchmarkGacheSetOnlySmallDataNoTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(NoTTL)
	benchmarkSetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) })
}

func BenchmarkGacheSetOnlyBigDataNoTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(NoTTL)
	benchmarkSetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) })
}

func BenchmarkTTLCacheSetOnlySmallDataNoTTL(b *testing.B) {
	cache := ttlcache.New[string, string]()
	benchmarkSetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { cache.Set(k, v, ttlcache.NoTTL) })
}

func BenchmarkTTLCacheSetOnlyBigDataNoTTL(b *testing.B) {
	cache := ttlcache.New[string, string]()
	benchmarkSetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { cache.Set(k, v, ttlcache.NoTTL) })
}

func BenchmarkGoCacheSetOnlySmallDataNoTTL(b *testing.B) {
	c := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	benchmarkSetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.Set(k, v, t) })
}

func BenchmarkGoCacheSetOnlyBigDataNoTTL(b *testing.B) {
	c := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	benchmarkSetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.Set(k, v, t) })
}

func BenchmarkBigCacheSetOnlySmallDataNoTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(bigCacheNoTTL)
	cfg.Verbose = false
	bc, _ := bigcache.NewBigCache(cfg)
	benchmarkSetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) })
}

func BenchmarkBigCacheSetOnlyBigDataNoTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(bigCacheNoTTL)
	cfg.Verbose = false
	bc, _ := bigcache.NewBigCache(cfg)
	benchmarkSetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) })
}

func BenchmarkGCacheLRUSetOnlySmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LRU().Build()
	benchmarkSetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) })
}

func BenchmarkGCacheLRUSetOnlyBigDataNoTTL(b *testing.B) {
	c := gcache.New(bigDataCount).LRU().Build()
	benchmarkSetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) })
}

func BenchmarkGCacheLFUSetOnlySmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LFU().Build()
	benchmarkSetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) })
}

func BenchmarkGCacheLFUSetOnlyBigDataNoTTL(b *testing.B) {
	c := gcache.New(bigDataCount).LFU().Build()
	benchmarkSetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) })
}

func BenchmarkGCacheARCSetOnlySmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).ARC().Build()
	benchmarkSetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) })
}

func BenchmarkGCacheARCSetOnlyBigDataNoTTL(b *testing.B) {
	c := gcache.New(bigDataCount).ARC().Build()
	benchmarkSetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) })
}

// ----- Get-only benchmarks (cache pre-populated before timing) -----

func BenchmarkDefaultMapGetSmallDataNoTTL(b *testing.B) {
	m := NewDefault()
	benchmarkGetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { m.Set(k, v) },
		func(k string) { m.Get(k) })
}

func BenchmarkDefaultMapGetBigDataNoTTL(b *testing.B) {
	m := NewDefault()
	benchmarkGetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { m.Set(k, v) },
		func(k string) { m.Get(k) })
}

func BenchmarkSyncMapGetSmallDataNoTTL(b *testing.B) {
	var m sync.Map
	benchmarkGetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { m.Store(k, v) },
		func(k string) { m.Load(k) })
}

func BenchmarkSyncMapGetBigDataNoTTL(b *testing.B) {
	var m sync.Map
	benchmarkGetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { m.Store(k, v) },
		func(k string) { m.Load(k) })
}

func BenchmarkGacheV2GetSmallDataNoTTL(b *testing.B) {
	g := gachev2.New(gachev2.WithDefaultExpiration[string](gachev2.NoTTL))
	benchmarkGetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheV2GetBigDataNoTTL(b *testing.B) {
	g := gachev2.New(
		gachev2.WithDefaultExpiration[string](gachev2.NoTTL),
		gachev2.WithMaxKeyLength[string](0))
	benchmarkGetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheGetSmallDataNoTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(NoTTL)
	benchmarkGetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkGacheGetBigDataNoTTL(b *testing.B) {
	g := gache.New().SetDefaultExpire(NoTTL)
	benchmarkGetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { g.Set(k, v) },
		func(k string) { g.Get(k) })
}

func BenchmarkTTLCacheGetSmallDataNoTTL(b *testing.B) {
	cache := ttlcache.New[string, string]()
	benchmarkGetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { cache.Set(k, v, ttlcache.NoTTL) },
		func(k string) { cache.Get(k) })
}

func BenchmarkTTLCacheGetBigDataNoTTL(b *testing.B) {
	cache := ttlcache.New[string, string]()
	benchmarkGetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { cache.Set(k, v, ttlcache.NoTTL) },
		func(k string) { cache.Get(k) })
}

func BenchmarkGoCacheGetSmallDataNoTTL(b *testing.B) {
	c := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	benchmarkGetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGoCacheGetBigDataNoTTL(b *testing.B) {
	c := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	benchmarkGetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.Set(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkBigCacheGetSmallDataNoTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(bigCacheNoTTL)
	cfg.Verbose = false
	bc, _ := bigcache.NewBigCache(cfg)
	benchmarkGetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}

func BenchmarkBigCacheGetBigDataNoTTL(b *testing.B) {
	cfg := bigcache.DefaultConfig(bigCacheNoTTL)
	cfg.Verbose = false
	bc, _ := bigcache.NewBigCache(cfg)
	benchmarkGetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { bc.Set(k, []byte(v)) },
		func(k string) { bc.Get(k) })
}

func BenchmarkGCacheLRUGetSmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LRU().Build()
	benchmarkGetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLRUGetBigDataNoTTL(b *testing.B) {
	c := gcache.New(bigDataCount).LRU().Build()
	benchmarkGetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLFUGetSmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).LFU().Build()
	benchmarkGetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheLFUGetBigDataNoTTL(b *testing.B) {
	c := gcache.New(bigDataCount).LFU().Build()
	benchmarkGetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheARCGetSmallDataNoTTL(b *testing.B) {
	c := gcache.New(len(smallData)).ARC().Build()
	benchmarkGetOnly(b, smallDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

func BenchmarkGCacheARCGetBigDataNoTTL(b *testing.B) {
	c := gcache.New(bigDataCount).ARC().Build()
	benchmarkGetOnly(b, bigDataSlice, NoTTL,
		func(k, v string, t time.Duration) { c.SetWithExpire(k, v, t) },
		func(k string) { c.Get(k) })
}

// func BenchmarkMCacheSetSmallDataNoTTL(b *testing.B) {
// 	c := mcache.StartInstance()
// 	benchmark(b, smallData, NoTTL,
// 		func(k, v string, t time.Duration) { c.SetPointer(k, v, t) },
// 		func(k string) { c.GetPointer(k) })
// }

// func BenchmarkMCacheSetSmallDataWithTTL(b *testing.B) {
// 	c := mcache.StartInstance()
// 	benchmark(b, smallData, ttl,
// 		func(k, v string, t time.Duration) { c.SetPointer(k, v, t) },
// 		func(k string) { c.GetPointer(k) })
// }

// func BenchmarkMCacheSetBigDataNoTTL(b *testing.B) {
// 	c := mcache.StartInstance()
// 	benchmark(b, bigData, NoTTL,
// 		func(k, v string, t time.Duration) { c.SetPointer(k, v, t) },
// 		func(k string) { c.GetPointer(k) })
// }

// func BenchmarkMCacheSetBigDataWithTTL(b *testing.B) {
// 	c := mcache.StartInstance()
// 	benchmark(b, bigData, ttl,
// 		func(k, v string, t time.Duration) { c.SetPointer(k, v, t) },
// 		func(k string) { c.GetPointer(k) })
// }

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

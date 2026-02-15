package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gujiaweiguo/goreport/internal/config"
)

type fakeProvider struct {
	store            map[string][]byte
	lastSetTTL       time.Duration
	failGet          error
	failSet          error
	failDeletePrefix error
}

func (f *fakeProvider) Get(ctx context.Context, key string) ([]byte, error) {
	if f.failGet != nil {
		return nil, f.failGet
	}
	value, ok := f.store[key]
	if !ok {
		return nil, nil
	}
	return value, nil
}

func (f *fakeProvider) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if f.failSet != nil {
		return f.failSet
	}
	f.lastSetTTL = ttl
	f.store[key] = value
	return nil
}

func (f *fakeProvider) Delete(ctx context.Context, key string) error {
	delete(f.store, key)
	return nil
}

func (f *fakeProvider) DeleteByPrefix(ctx context.Context, prefix string) error {
	if f.failDeletePrefix != nil {
		return f.failDeletePrefix
	}
	for key := range f.store {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			delete(f.store, key)
		}
	}
	return nil
}

func (f *fakeProvider) Close() error { return nil }

func TestNewCache(t *testing.T) {
	cache, err := New(config.CacheConfig{})
	assert.NoError(t, err)
	assert.NotNil(t, cache)
}

func TestCacheProvider(t *testing.T) {
	cache, err := New(config.CacheConfig{})
	assert.NoError(t, err)
	assert.NotNil(t, cache)
}

func TestCache_HitMiss(t *testing.T) {
	provider := &fakeProvider{store: map[string][]byte{}}
	cache := &Cache{provider: provider, cfg: config.CacheConfig{DefaultTTL: 60}, metrics: &Metrics{}}
	ctx := context.Background()

	require.NoError(t, cache.Set(ctx, "test-tenant", "datasource:tables", "ds-1", nil, []byte("ok")))

	value, hit, err := cache.Get(ctx, "test-tenant", "datasource:tables", "ds-1", nil)
	require.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, []byte("ok"), value)

	value, hit, err = cache.Get(ctx, "test-tenant", "datasource:tables", "ds-not-found", nil)
	require.NoError(t, err)
	assert.False(t, hit)
	assert.Nil(t, value)

	metrics := cache.GetMetrics()
	assert.Equal(t, int64(1), metrics.Hits)
	assert.Equal(t, int64(1), metrics.Misses)
}

func TestCache_TTL_DefaultAndOverride(t *testing.T) {
	provider := &fakeProvider{store: map[string][]byte{}}
	cache := &Cache{provider: provider, cfg: config.CacheConfig{DefaultTTL: 120}, metrics: &Metrics{}}
	ctx := context.Background()

	require.NoError(t, cache.Set(ctx, "test-tenant", "datasource:tables", "ds-1", nil, []byte("v1")))
	assert.Equal(t, 120*time.Second, provider.lastSetTTL)

	require.NoError(t, cache.Set(ctx, "test-tenant", "datasource:tables", "ds-2", nil, []byte("v2"), 10*time.Second))
	assert.Equal(t, 10*time.Second, provider.lastSetTTL)
}

func TestCache_DegradedMode(t *testing.T) {
	provider := &fakeProvider{
		store:            map[string][]byte{},
		failGet:          errors.New("redis down"),
		failSet:          errors.New("redis down"),
		failDeletePrefix: errors.New("redis down"),
	}
	cache := &Cache{provider: provider, cfg: config.CacheConfig{DefaultTTL: 60}, metrics: &Metrics{}, degraded: true}
	ctx := context.Background()

	value, hit, err := cache.Get(ctx, "test-tenant", "datasource:tables", "ds-1", nil)
	assert.NoError(t, err)
	assert.False(t, hit)
	assert.Nil(t, value)

	err = cache.Set(ctx, "test-tenant", "datasource:tables", "ds-1", nil, []byte("v"))
	assert.NoError(t, err)

	err = cache.Invalidate(ctx, "test-tenant", "datasource:tables")
	assert.NoError(t, err)

	metrics := cache.GetMetrics()
	assert.Equal(t, int64(3), metrics.Failures)
	assert.True(t, cache.IsDegraded())
}

func TestCacheKeyBuilding(t *testing.T) {
	tenantID := "tenant-1"
	domain := "datasource:tables"
	identity := "ds-1"

	key := BuildKey(tenantID, domain, identity, "")
	assert.Contains(t, key, "tenant-1")
	assert.Contains(t, key, "datasource:tables")
	assert.Contains(t, key, "ds-1")
}

func TestHashParams(t *testing.T) {
	t.Run("empty params", func(t *testing.T) {
		result := HashParams(nil)
		assert.Equal(t, "none", result)
	})

	t.Run("empty map", func(t *testing.T) {
		result := HashParams(map[string]interface{}{})
		assert.Equal(t, "none", result)
	})

	t.Run("with params", func(t *testing.T) {
		params := map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		}
		result := HashParams(params)
		assert.Len(t, result, 8)
	})
}

func TestCache_GetHitRate(t *testing.T) {
	t.Run("no requests", func(t *testing.T) {
		cache := &Cache{metrics: &Metrics{}}
		assert.Equal(t, 0.0, cache.GetHitRate())
	})

	t.Run("with hits and misses", func(t *testing.T) {
		cache := &Cache{metrics: &Metrics{Hits: 3, Misses: 1}}
		assert.Equal(t, 0.75, cache.GetHitRate())
	})

	t.Run("all hits", func(t *testing.T) {
		cache := &Cache{metrics: &Metrics{Hits: 5, Misses: 0}}
		assert.Equal(t, 1.0, cache.GetHitRate())
	})
}

func TestCache_Close(t *testing.T) {
	t.Run("nil provider", func(t *testing.T) {
		cache := &Cache{provider: nil}
		err := cache.Close()
		assert.NoError(t, err)
	})

	t.Run("with provider", func(t *testing.T) {
		provider := &fakeProvider{store: map[string][]byte{}}
		cache := &Cache{provider: provider}
		err := cache.Close()
		assert.NoError(t, err)
	})
}

func TestCache_ExportMetrics(t *testing.T) {
	cache := &Cache{
		metrics:  &Metrics{Hits: 10, Misses: 5, Failures: 2, Errors: 1},
		degraded: false,
	}

	metrics := cache.ExportMetrics()
	assert.Equal(t, "10", metrics["cache_hits"])
	assert.Equal(t, "5", metrics["cache_misses"])
	assert.Equal(t, "2", metrics["cache_failures"])
	assert.Equal(t, "1", metrics["cache_errors"])
	assert.Contains(t, metrics["cache_hit_rate"], "0.6667")
	assert.Equal(t, "false", metrics["cache_degraded"])
}

func TestBuildPrefix(t *testing.T) {
	prefix := BuildPrefix("tenant-1", "datasource")
	assert.Contains(t, prefix, "tenant-1")
	assert.Contains(t, prefix, "datasource")
	assert.True(t, len(prefix) > 0 && prefix[len(prefix)-1] == ':')
}

func TestCache_Invalidate(t *testing.T) {
	provider := &fakeProvider{store: map[string][]byte{}}
	cache := &Cache{provider: provider, cfg: config.CacheConfig{DefaultTTL: 60}, metrics: &Metrics{}}
	ctx := context.Background()

	key := BuildKey("tenant-1", "domain", "id", "none")
	provider.store[key] = []byte("value")

	err := cache.Invalidate(ctx, "tenant-1", "domain")
	assert.NoError(t, err)
}

func TestCache_NonDegradedFailure(t *testing.T) {
	provider := &fakeProvider{
		store:   map[string][]byte{},
		failGet: errors.New("redis error"),
	}
	cache := &Cache{provider: provider, cfg: config.CacheConfig{DefaultTTL: 60}, metrics: &Metrics{}, degraded: false}
	ctx := context.Background()

	_, _, err := cache.Get(ctx, "tenant-1", "domain", "id", nil)
	assert.Error(t, err)
}

func TestCache_Set_NonDegradedFailure(t *testing.T) {
	provider := &fakeProvider{
		store:   map[string][]byte{},
		failSet: errors.New("redis error"),
	}
	cache := &Cache{provider: provider, cfg: config.CacheConfig{DefaultTTL: 60}, metrics: &Metrics{}, degraded: false}
	ctx := context.Background()

	err := cache.Set(ctx, "tenant-1", "domain", "id", nil, []byte("value"))
	assert.Error(t, err)
}

func TestCache_Invalidate_NonDegradedFailure(t *testing.T) {
	provider := &fakeProvider{
		store:            map[string][]byte{},
		failDeletePrefix: errors.New("redis error"),
	}
	cache := &Cache{provider: provider, cfg: config.CacheConfig{DefaultTTL: 60}, metrics: &Metrics{}, degraded: false}
	ctx := context.Background()

	err := cache.Invalidate(ctx, "tenant-1", "domain")
	assert.Error(t, err)
}

func TestNew_EnabledButRedisUnavailable(t *testing.T) {
	cache, err := New(config.CacheConfig{
		Enabled:    true,
		Addr:       "invalid-host:6379",
		DefaultTTL: 60,
	})

	assert.NoError(t, err)
	assert.NotNil(t, cache)
	assert.True(t, cache.IsDegraded())
}

func TestNew_Disabled(t *testing.T) {
	cache, err := New(config.CacheConfig{
		Enabled:    false,
		DefaultTTL: 60,
	})

	assert.NoError(t, err)
	assert.NotNil(t, cache)
	assert.True(t, cache.IsDegraded())
}

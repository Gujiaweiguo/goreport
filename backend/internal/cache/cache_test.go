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

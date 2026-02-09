package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gujiaweiguo/goreport/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func testCacheDB() int {
	v := getEnvOrDefault("CACHE_DB", "0")
	db, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return db
}

func newRedisBackedCache(t *testing.T) *Cache {
	t.Helper()

	addr := getEnvOrDefault("TEST_REDIS_ADDR", "")
	if addr == "" {
		addr = getEnvOrDefault("REDIS_ADDR", "")
	}
	if addr == "" {
		addr = getEnvOrDefault("CACHE_ADDR", "")
	}
	if addr == "" {
		t.Skip("redis integration test skipped: TEST_REDIS_ADDR/REDIS_ADDR/CACHE_ADDR not set")
	}

	password := getEnvOrDefault("CACHE_PASSWORD", "")
	if password == "" {
		password = getEnvOrDefault("REDIS_PASSWORD", "")
	}

	cfg := config.CacheConfig{
		Enabled:    true,
		Addr:       addr,
		Password:   password,
		DB:         testCacheDB(),
		DefaultTTL: 60,
	}

	provider, err := NewRedisProvider(cfg)
	if err != nil {
		t.Skipf("redis unavailable for integration test: %v", err)
	}

	cache := &Cache{provider: provider, cfg: cfg, metrics: &Metrics{}}
	t.Cleanup(func() {
		_ = cache.Close()
	})

	return cache
}

func uniqueTenant(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func TestRedisIntegration_TTLExpiration(t *testing.T) {
	cache := newRedisBackedCache(t)
	ctx := context.Background()

	tenantID := uniqueTenant("tenant-ttl")
	domain := "datasource:tables"
	identity := "ds-ttl"

	t.Cleanup(func() {
		_ = cache.Invalidate(ctx, tenantID, domain)
	})

	require.NoError(t, cache.Set(ctx, tenantID, domain, identity, nil, []byte("value"), 300*time.Millisecond))

	value, hit, err := cache.Get(ctx, tenantID, domain, identity, nil)
	require.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, []byte("value"), value)

	deadline := time.Now().Add(3 * time.Second)
	for {
		_, hit, err = cache.Get(ctx, tenantID, domain, identity, nil)
		require.NoError(t, err)
		if !hit {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("expected key to expire, but it still exists")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func TestRedisIntegration_DeleteByPrefix(t *testing.T) {
	cache := newRedisBackedCache(t)
	ctx := context.Background()

	tenantID := uniqueTenant("tenant-prefix")
	domain := "datasource:tables"
	otherDomain := "dashboard:query"

	t.Cleanup(func() {
		_ = cache.Invalidate(ctx, tenantID, domain)
		_ = cache.Invalidate(ctx, tenantID, otherDomain)
	})

	require.NoError(t, cache.Set(ctx, tenantID, domain, "a", nil, []byte("A"), 5*time.Second))
	require.NoError(t, cache.Set(ctx, tenantID, domain, "b", nil, []byte("B"), 5*time.Second))
	require.NoError(t, cache.Set(ctx, tenantID, otherDomain, "c", nil, []byte("C"), 5*time.Second))

	_, hitA, err := cache.Get(ctx, tenantID, domain, "a", nil)
	require.NoError(t, err)
	assert.True(t, hitA)

	_, hitB, err := cache.Get(ctx, tenantID, domain, "b", nil)
	require.NoError(t, err)
	assert.True(t, hitB)

	valueC, hitC, err := cache.Get(ctx, tenantID, otherDomain, "c", nil)
	require.NoError(t, err)
	assert.True(t, hitC)
	assert.Equal(t, []byte("C"), valueC)

	require.NoError(t, cache.Invalidate(ctx, tenantID, domain))

	_, hitA, err = cache.Get(ctx, tenantID, domain, "a", nil)
	require.NoError(t, err)
	assert.False(t, hitA)

	_, hitB, err = cache.Get(ctx, tenantID, domain, "b", nil)
	require.NoError(t, err)
	assert.False(t, hitB)

	valueC, hitC, err = cache.Get(ctx, tenantID, otherDomain, "c", nil)
	require.NoError(t, err)
	assert.True(t, hitC)
	assert.Equal(t, []byte("C"), valueC)
}

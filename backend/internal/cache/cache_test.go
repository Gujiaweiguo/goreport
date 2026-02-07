package cache

import (
	"context"
	"testing"
	"time"

	"github.com/jeecg/jimureport-go/internal/config"
)

func TestCacheProvider(t *testing.T) {
	cfg := config.CacheConfig{
		Enabled:    false,
		Addr:       "localhost:6379",
		Password:   "",
		DB:         0,
		DefaultTTL: 3600,
	}

	cache, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}
	defer cache.Close()

	if !cache.IsDegraded() {
		t.Error("Expected cache to be degraded when disabled")
	}
}

func TestCacheSetAndGet(t *testing.T) {
	cfg := config.CacheConfig{
		Enabled:    false,
		Addr:       "localhost:6379",
		Password:   "",
		DB:         0,
		DefaultTTL: 3600,
	}

	cache, _ := New(cfg)
	defer cache.Close()

	ctx := context.Background()
	tenantID := "tenant-1"
	domain := "test"
	identity := "key1"

	value := []byte("test-value")

	err := cache.Set(ctx, tenantID, domain, identity, nil, value)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	retrieved, hit, err := cache.Get(ctx, tenantID, domain, identity, nil)
	if err != nil {
		t.Fatalf("Failed to get cache: %v", err)
	}
	if !hit {
		t.Error("Expected cache hit")
	}
	if string(retrieved) != string(value) {
		t.Errorf("Expected %s, got %s", string(value), string(retrieved))
	}

	metrics := cache.GetMetrics()
	if metrics.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", metrics.Hits)
	}
}

func TestCacheInvalidation(t *testing.T) {
	cfg := config.CacheConfig{
		Enabled:    false,
		Addr:       "localhost:6379",
		Password:   "",
		DB:         0,
		DefaultTTL: 3600,
	}

	cache, _ := New(cfg)
	defer cache.Close()

	ctx := context.Background()
	tenantID := "tenant-1"
	domain := "test"
	identity := "key1"

	value := []byte("test-value")
	_ = cache.Set(ctx, tenantID, domain, identity, nil, value)

	err := cache.Invalidate(ctx, tenantID, domain)
	if err != nil {
		t.Fatalf("Failed to invalidate cache: %v", err)
	}

	retrieved, hit, _ := cache.Get(ctx, tenantID, domain, identity, nil)
	if hit {
		t.Error("Expected cache miss after invalidation")
	}
	if retrieved != nil {
		t.Error("Expected nil after invalidation")
	}
}

func TestCacheKeyBuilding(t *testing.T) {
	tenantID := "tenant-1"
	domain := "test-domain"
	identity := "key1"

	params := map[string]interface{}{
		"param1": "value1",
		"param2": 42,
	}

	paramsHash := HashParams(params)
	if paramsHash == "" {
		t.Error("Expected non-empty params hash")
	}

	key := BuildKey(tenantID, domain, identity, paramsHash)
	expectedPrefix := "jr:tenant-1:test-domain:key1:"
	if len(key) <= len(expectedPrefix) || key[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Expected key to start with %s, got %s", expectedPrefix, key)
	}
}

func TestCacheMetricsExport(t *testing.T) {
	cfg := config.CacheConfig{
		Enabled:    false,
		Addr:       "localhost:6379",
		Password:   "",
		DB:         0,
		DefaultTTL: 3600,
	}

	cache, _ := New(cfg)
	defer cache.Close()

	ctx := context.Background()

	for i := 0; i < 10; i++ {
		_ = cache.Set(ctx, "tenant-1", "test", "key", nil, []byte("value"))
	}

	for i := 0; i < 5; i++ {
		_, _, _ = cache.Get(ctx, "tenant-1", "test", "key", nil)
	}

	metrics := cache.ExportMetrics()
	if metrics["cache_hits"] != "5" {
		t.Errorf("Expected 5 hits, got %s", metrics["cache_hits"])
	}
	if metrics["cache_misses"] != "5" {
		t.Errorf("Expected 5 misses, got %s", metrics["cache_misses"])
	}

	hitRate := cache.GetHitRate()
	if hitRate != 0.5 {
		t.Errorf("Expected hit rate 0.5, got %f", hitRate)
	}
}

func TestCacheDegradedMode(t *testing.T) {
	cfg := config.CacheConfig{
		Enabled:    false,
		Addr:       "localhost:6379",
		Password:   "",
		DB:         0,
		DefaultTTL: 3600,
	}

	cache, _ := New(cfg)
	defer cache.Close()

	ctx := context.Background()
	tenantID := "tenant-1"
	domain := "test"
	identity := "key1"

	value := []byte("test-value")

	err := cache.Set(ctx, tenantID, domain, identity, nil, value)
	if err != nil {
		t.Fatalf("Failed to set cache in degraded mode: %v", err)
	}

	retrieved, hit, err := cache.Get(ctx, tenantID, domain, identity, nil)
	if err != nil {
		t.Fatalf("Failed to get cache in degraded mode: %v", err)
	}
	if hit {
		t.Error("Expected cache miss in degraded mode")
	}
	if retrieved != nil {
		t.Error("Expected nil in degraded mode")
	}

	if !cache.IsDegraded() {
		t.Error("Expected cache to be in degraded mode")
	}
}

func TestCacheTTLOverride(t *testing.T) {
	cfg := config.CacheConfig{
		Enabled:    false,
		Addr:       "localhost:6379",
		Password:   "",
		DB:         0,
		DefaultTTL: 3600,
	}

	cache, _ := New(cfg)
	defer cache.Close()

	ctx := context.Background()
	tenantID := "tenant-1"
	domain := "test"
	identity := "key1"

	value := []byte("test-value")

	_ = cache.Set(ctx, tenantID, domain, identity, nil, value, 10*time.Minute)

	time.Sleep(11 * time.Millisecond)

	retrieved, hit, _ := cache.Get(ctx, tenantID, domain, identity, nil)
	if !hit {
		t.Error("Expected cache hit after set")
	}
	if string(retrieved) != string(value) {
		t.Errorf("Expected %s, got %s", string(value), string(retrieved))
	}
}

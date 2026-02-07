package cache

import (
	"context"
	"testing"
	"time"

	"github.com/gujiaweiguo/goreport/internal/config"
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

	metrics := cache.GetMetrics()
	if metrics.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", metrics.Misses)
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

func TestCacheDegradedModeSet(t *testing.T) {
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

	if !cache.IsDegraded() {
		t.Error("Expected cache to be in degraded mode")
	}
}

func TestCacheDegradedModeTTL(t *testing.T) {
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

	err := cache.Set(ctx, tenantID, domain, identity, nil, value, 10*time.Minute)
	if err != nil {
		t.Fatalf("Failed to set cache with TTL in degraded mode: %v", err)
	}
}

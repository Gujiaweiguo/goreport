package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jeecg/jimureport-go/internal/config"
)

type Provider interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	DeleteByPrefix(ctx context.Context, prefix string) error
	Close() error
}

const (
	keyPrefix = "jr"
	sep       = ":"
)

func BuildKey(tenantID, domain, identity string, paramsHash string) string {
	key := strings.Join([]string{
		keyPrefix,
		tenantID,
		domain,
		identity,
		paramsHash,
	}, sep)
	return key
}

func BuildPrefix(tenantID, domain string) string {
	return strings.Join([]string{
		keyPrefix,
		tenantID,
		domain,
	}, sep) + sep
}

func HashParams(params map[string]interface{}) string {
	if len(params) == 0 {
		return "none"
	}

	key := strings.Builder{}
	for k, v := range params {
		key.WriteString(k)
		key.WriteString("=")
		key.WriteString(fmt.Sprintf("%v", v))
		key.WriteString("&")
	}

	hash := sha256.Sum256([]byte(key.String()))
	return hex.EncodeToString(hash[:])[:8]
}

type Cache struct {
	provider Provider
	cfg      config.CacheConfig
	metrics  *Metrics
	degraded bool
}

type Metrics struct {
	Hits     int64
	Misses   int64
	Failures int64
	Errors   int64
}

func New(cfg config.CacheConfig) (*Cache, error) {
	c := &Cache{
		cfg:     cfg,
		metrics: &Metrics{},
	}

	if !cfg.Enabled {
		c.provider = &NoopProvider{}
		c.degraded = true
		return c, nil
	}

	redisProvider, err := NewRedisProvider(cfg)
	if err != nil {
		c.provider = &NoopProvider{}
		c.degraded = true
		return c, nil
	}

	c.provider = redisProvider
	return c, nil
}

func (c *Cache) Get(ctx context.Context, tenantID, domain, identity string, params map[string]interface{}) ([]byte, bool, error) {
	paramsHash := HashParams(params)
	key := BuildKey(tenantID, domain, identity, paramsHash)

	value, err := c.provider.Get(ctx, key)
	if err != nil {
		c.metrics.Failures++
		if c.degraded {
			return nil, false, nil
		}
		return nil, false, err
	}

	if value == nil {
		c.metrics.Misses++
		return nil, false, nil
	}

	c.metrics.Hits++
	return value, true, nil
}

func (c *Cache) Set(ctx context.Context, tenantID, domain, identity string, params map[string]interface{}, value []byte, ttlOverride ...time.Duration) error {
	paramsHash := HashParams(params)
	key := BuildKey(tenantID, domain, identity, paramsHash)

	ttl := time.Duration(c.cfg.DefaultTTL) * time.Second
	if len(ttlOverride) > 0 {
		ttl = ttlOverride[0]
	}

	err := c.provider.Set(ctx, key, value, ttl)
	if err != nil {
		c.metrics.Failures++
		if c.degraded {
			return nil
		}
		return err
	}

	return nil
}

func (c *Cache) Invalidate(ctx context.Context, tenantID, domain string) error {
	prefix := BuildPrefix(tenantID, domain)
	err := c.provider.DeleteByPrefix(ctx, prefix)
	if err != nil {
		c.metrics.Failures++
		if c.degraded {
			return nil
		}
		return err
	}

	return nil
}

func (c *Cache) GetMetrics() *Metrics {
	return &Metrics{
		Hits:     c.metrics.Hits,
		Misses:   c.metrics.Misses,
		Failures: c.metrics.Failures,
		Errors:   c.metrics.Errors,
	}
}

func (c *Cache) GetHitRate() float64 {
	total := c.metrics.Hits + c.metrics.Misses
	if total == 0 {
		return 0
	}
	return float64(c.metrics.Hits) / float64(total)
}

func (c *Cache) IsDegraded() bool {
	return c.degraded
}

func (c *Cache) Close() error {
	if c.provider != nil {
		return c.provider.Close()
	}
	return nil
}

func (c *Cache) ExportMetrics() map[string]string {
	return map[string]string{
		"cache_hits":     strconv.FormatInt(c.metrics.Hits, 10),
		"cache_misses":   strconv.FormatInt(c.metrics.Misses, 10),
		"cache_failures": strconv.FormatInt(c.metrics.Failures, 10),
		"cache_errors":   strconv.FormatInt(c.metrics.Errors, 10),
		"cache_hit_rate": fmt.Sprintf("%.4f", c.GetHitRate()),
		"cache_degraded": strconv.FormatBool(c.degraded),
	}
}

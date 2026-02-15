package auth

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gujiaweiguo/goreport/internal/cache"
	"github.com/gujiaweiguo/goreport/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlacklist_WithNilStore(t *testing.T) {
	blacklistStore = nil

	err := RevokeToken(context.Background(), "token-1", time.Now().Add(time.Minute))
	require.NoError(t, err)
	assert.False(t, IsTokenRevoked(context.Background(), "token-1"))
}

func TestBlacklist_RevokeAndCheck_WithRedis(t *testing.T) {
	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		addr = os.Getenv("REDIS_ADDR")
	}
	if addr == "" {
		t.Skip("TEST_REDIS_ADDR or REDIS_ADDR not set")
	}

	c, err := cache.New(config.CacheConfig{
		Enabled:    true,
		Addr:       addr,
		Password:   os.Getenv("REDIS_PASSWORD"),
		DefaultTTL: 60,
	})
	require.NoError(t, err)
	defer c.Close()

	if c.IsDegraded() {
		t.Skip("redis unavailable, cache degraded to noop")
	}

	InitBlacklist(c)

	// Use unique token to avoid conflicts with other tests
	token := fmt.Sprintf("token-test-revoke-%d", time.Now().UnixNano())
	assert.False(t, IsTokenRevoked(context.Background(), token))

	err = RevokeToken(context.Background(), token, time.Now().Add(2*time.Minute))
	require.NoError(t, err)
	assert.True(t, IsTokenRevoked(context.Background(), token))

	// 过期时间已过时应忽略写入。
	expiredToken := fmt.Sprintf("expired-token-%d", time.Now().UnixNano())
	err = RevokeToken(context.Background(), expiredToken, time.Now().Add(-time.Minute))
	require.NoError(t, err)
	assert.False(t, IsTokenRevoked(context.Background(), expiredToken))
}

func TestBlacklist_EmptyToken(t *testing.T) {
	blacklistStore = nil
	err := RevokeToken(context.Background(), "", time.Now().Add(time.Minute))
	require.NoError(t, err)
	assert.False(t, IsTokenRevoked(context.Background(), ""))
}

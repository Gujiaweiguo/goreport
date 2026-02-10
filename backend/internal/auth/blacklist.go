package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/gujiaweiguo/goreport/internal/cache"
)

const (
	blacklistKeyPrefix = "jrt:" // JSON Token Revoked
)

type BlacklistStore struct {
	cache *cache.Cache
}

var blacklistStore *BlacklistStore

func InitBlacklist(cache *cache.Cache) {
	blacklistStore = &BlacklistStore{
		cache: cache,
	}
}

func RevokeToken(ctx context.Context, token string, expiresAt time.Time) error {
	if blacklistStore == nil || blacklistStore.cache == nil {
		return nil
	}

	if token == "" {
		return nil
	}

	ttl := time.Until(expiresAt)
	if ttl < 0 {
		return nil
	}

	key := fmt.Sprintf("%s%s", blacklistKeyPrefix, token)

	return blacklistStore.cache.Set(ctx, "default", "auth_blacklist", key, nil, []byte("1"), ttl)
}

func IsTokenRevoked(ctx context.Context, token string) bool {
	if blacklistStore == nil || blacklistStore.cache == nil {
		return false
	}

	if token == "" {
		return false
	}

	key := fmt.Sprintf("%s%s", blacklistKeyPrefix, token)

	val, found, err := blacklistStore.cache.Get(ctx, "default", "auth_blacklist", key, nil)
	if err != nil || !found || val == nil {
		return false
	}

	return string(val) == "1"
}

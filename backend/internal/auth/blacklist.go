package auth

import (
	"sync"
	"time"
)

type tokenBlacklist struct {
	mu     sync.Mutex
	tokens map[string]time.Time
}

var blacklist = &tokenBlacklist{tokens: make(map[string]time.Time)}

// RevokeToken marks the token as revoked until its expiration time.
func RevokeToken(token string, expiresAt time.Time) {
	if token == "" {
		return
	}

	blacklist.mu.Lock()
	defer blacklist.mu.Unlock()

	blacklist.tokens[token] = expiresAt
}

// IsTokenRevoked returns true if the token is revoked and not expired.
func IsTokenRevoked(token string) bool {
	if token == "" {
		return false
	}

	now := time.Now()

	blacklist.mu.Lock()
	defer blacklist.mu.Unlock()

	for value, expiresAt := range blacklist.tokens {
		if expiresAt.Before(now) {
			delete(blacklist.tokens, value)
		}
	}

	expiresAt, ok := blacklist.tokens[token]
	if !ok {
		return false
	}
	if expiresAt.Before(now) {
		delete(blacklist.tokens, token)
		return false
	}
	return true
}

package cache

import (
	"context"
	"time"
)

type NoopProvider struct{}

func NewNoopProvider() *NoopProvider {
	return &NoopProvider{}
}

func (n *NoopProvider) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

func (n *NoopProvider) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return nil
}

func (n *NoopProvider) Delete(ctx context.Context, key string) error {
	return nil
}

func (n *NoopProvider) DeleteByPrefix(ctx context.Context, prefix string) error {
	return nil
}

func (n *NoopProvider) Close() error {
	return nil
}

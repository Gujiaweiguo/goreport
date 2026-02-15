package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewNoopProvider(t *testing.T) {
	provider := NewNoopProvider()
	assert.NotNil(t, provider)
}

func TestNoopProvider_Get(t *testing.T) {
	provider := NewNoopProvider()
	ctx := context.Background()

	val, err := provider.Get(ctx, "any-key")

	assert.Nil(t, val)
	assert.NoError(t, err)
}

func TestNoopProvider_Set(t *testing.T) {
	provider := NewNoopProvider()
	ctx := context.Background()

	err := provider.Set(ctx, "any-key", []byte("any-value"), time.Minute)

	assert.NoError(t, err)
}

func TestNoopProvider_Delete(t *testing.T) {
	provider := NewNoopProvider()
	ctx := context.Background()

	err := provider.Delete(ctx, "any-key")

	assert.NoError(t, err)
}

func TestNoopProvider_DeleteByPrefix(t *testing.T) {
	provider := NewNoopProvider()
	ctx := context.Background()

	err := provider.DeleteByPrefix(ctx, "prefix-")

	assert.NoError(t, err)
}

func TestNoopProvider_Close(t *testing.T) {
	provider := NewNoopProvider()

	err := provider.Close()

	assert.NoError(t, err)
}

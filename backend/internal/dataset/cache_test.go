package dataset

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpressionCache_Get_Missing(t *testing.T) {
	cache := NewExpressionCache()
	val, ok := cache.Get("nonexistent")
	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestExpressionCache_SetAndGet(t *testing.T) {
	cache := NewExpressionCache()
	cache.Set("key1", "value1", time.Hour)

	val, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", val)
}

func TestExpressionCache_Get_Expired(t *testing.T) {
	cache := NewExpressionCache()
	cache.Set("key1", "value1", 10*time.Millisecond)

	time.Sleep(20 * time.Millisecond)

	val, ok := cache.Get("key1")
	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestExpressionCache_Invalidate(t *testing.T) {
	cache := NewExpressionCache()
	cache.Set("key1", "value1", time.Hour)

	cache.Invalidate("key1")

	val, ok := cache.Get("key1")
	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestExpressionCache_Invalidate_NonExistent(t *testing.T) {
	cache := NewExpressionCache()
	cache.Invalidate("nonexistent")
}

func TestExpressionCache_Clear(t *testing.T) {
	cache := NewExpressionCache()
	cache.Set("key1", "value1", time.Hour)
	cache.Set("key2", "value2", time.Hour)

	cache.Clear()

	val1, ok1 := cache.Get("key1")
	val2, ok2 := cache.Get("key2")

	assert.False(t, ok1)
	assert.Nil(t, val1)
	assert.False(t, ok2)
	assert.Nil(t, val2)
}

func TestExpressionCache_ConcurrentAccess(t *testing.T) {
	cache := NewExpressionCache()
	done := make(chan bool)

	// Concurrent writes
	for i := 0; i < 10; i++ {
		go func(n int) {
			cache.Set("key", n, time.Hour)
			done <- true
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			cache.Get("key")
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}
}

func TestComputedFieldCache_New(t *testing.T) {
	cache := NewComputedFieldCache()
	assert.NotNil(t, cache)
}

func TestComputedFieldCache_Expression(t *testing.T) {
	cache := NewComputedFieldCache()

	// Get missing
	val, ok := cache.GetExpression("field1")
	assert.False(t, ok)
	assert.Equal(t, "", val)

	// Set and get
	cache.SetExpression("field1", "[amount] * 2", time.Hour)
	val, ok = cache.GetExpression("field1")
	assert.True(t, ok)
	assert.Equal(t, "[amount] * 2", val)
}

func TestComputedFieldCache_SQL(t *testing.T) {
	cache := NewComputedFieldCache()

	val, ok := cache.GetSQL("field1")
	assert.False(t, ok)
	assert.Equal(t, "", val)

	cache.SetSQL("field1", "amount * 2", time.Hour)
	val, ok = cache.GetSQL("field1")
	assert.True(t, ok)
	assert.Equal(t, "amount * 2", val)
}

func TestComputedFieldCache_InvalidateField(t *testing.T) {
	cache := NewComputedFieldCache()

	cache.SetExpression("field1", "[amount] * 2", time.Hour)

	cache.InvalidateField("field1")

	val, ok := cache.GetExpression("field1")
	assert.False(t, ok)
	assert.Equal(t, "", val)
}

func TestComputedFieldCache_Clear(t *testing.T) {
	cache := NewComputedFieldCache()

	cache.SetExpression("field1", "[amount] * 2", time.Hour)
	cache.SetExpression("field2", "[price] * [qty]", time.Hour)

	cache.Clear()

	val1, ok1 := cache.GetExpression("field1")
	val2, ok2 := cache.GetExpression("field2")

	assert.False(t, ok1)
	assert.Equal(t, "", val1)
	assert.False(t, ok2)
	assert.Equal(t, "", val2)
}

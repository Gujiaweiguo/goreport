package dataset

import (
	"sync"
	"time"
)

type ExpressionCache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Invalidate(key string)
	Clear()
}

type expressionCache struct {
	items map[string]cacheItem
	mu    sync.RWMutex
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

func NewExpressionCache() ExpressionCache {
	return &expressionCache{
		items: make(map[string]cacheItem),
	}
}

func (c *expressionCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

func (c *expressionCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

func (c *expressionCache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *expressionCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]cacheItem)
}

type ComputedFieldCache struct {
	expressions ExpressionCache
	sqlCache    ExpressionCache
	jsCache     ExpressionCache
}

func NewComputedFieldCache() *ComputedFieldCache {
	return &ComputedFieldCache{
		expressions: NewExpressionCache(),
		sqlCache:    NewExpressionCache(),
		jsCache:     NewExpressionCache(),
	}
}

func (c *ComputedFieldCache) GetExpression(fieldID string) (string, bool) {
	if val, ok := c.expressions.Get(fieldID); ok {
		if expr, ok := val.(string); ok {
			return expr, true
		}
	}
	return "", false
}

func (c *ComputedFieldCache) SetExpression(fieldID, expression string, ttl time.Duration) {
	c.expressions.Set(fieldID, expression, ttl)
}

func (c *ComputedFieldCache) GetSQL(fieldID string) (string, bool) {
	if val, ok := c.sqlCache.Get(fieldID); ok {
		if sql, ok := val.(string); ok {
			return sql, true
		}
	}
	return "", false
}

func (c *ComputedFieldCache) SetSQL(fieldID, sql string, ttl time.Duration) {
	c.sqlCache.Set(fieldID, sql, ttl)
}

func (c *ComputedFieldCache) InvalidateField(fieldID string) {
	c.expressions.Invalidate(fieldID)
	c.sqlCache.Invalidate(fieldID)
	c.jsCache.Invalidate(fieldID)
}

func (c *ComputedFieldCache) Clear() {
	c.expressions.Clear()
	c.sqlCache.Clear()
	c.jsCache.Clear()
}

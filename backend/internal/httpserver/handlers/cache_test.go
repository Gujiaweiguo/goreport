package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/cache"
	"github.com/gujiaweiguo/goreport/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewCacheHandler(t *testing.T) {
	handler := NewCacheHandler(nil)
	assert.NotNil(t, handler)
}

func TestCacheHandler_GetMetrics_NilCache(t *testing.T) {
	handler := NewCacheHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/cache/metrics", nil)

	handler.GetMetrics(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "cache is disabled")
}

func TestCacheHandler_GetMetrics_WithCache(t *testing.T) {
	cfg := config.CacheConfig{
		Enabled:    false,
		DefaultTTL: 60,
	}
	testCache, _ := cache.New(cfg)
	handler := NewCacheHandler(testCache)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/cache/metrics", nil)

	handler.GetMetrics(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

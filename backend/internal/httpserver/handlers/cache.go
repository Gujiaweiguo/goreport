package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeecg/jimureport-go/internal/cache"
)

type CacheHandler struct {
	cache *cache.Cache
}

func NewCacheHandler(cache *cache.Cache) *CacheHandler {
	return &CacheHandler{cache: cache}
}

func (h *CacheHandler) GetMetrics(c *gin.Context) {
	if h.cache == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "cache is disabled",
		})
		return
	}

	metrics := h.cache.ExportMetrics()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  metrics,
		"message": "success",
	})
}

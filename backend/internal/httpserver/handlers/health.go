package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Check 处理健康检查请求
func (h *HealthHandler) Check(c *gin.Context) {
	status := gin.H{
		"status": "ok",
	}

	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err == nil {
			if err := sqlDB.Ping(); err == nil {
				status["database"] = "connected"
			} else {
				status["database"] = "disconnected"
			}
		}
	}

	c.JSON(200, gin.H{
		"success": true,
		"result":  status,
		"message": "success",
	})
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/auth"
	"github.com/gujiaweiguo/goreport/internal/repository"
)

type TenantHandler struct {
	repo repository.TenantRepository
}

func NewTenantHandler(repo repository.TenantRepository) *TenantHandler {
	return &TenantHandler{repo: repo}
}

func (h *TenantHandler) GetCurrent(c *gin.Context) {
	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "tenant not found in context",
		})
		return
	}

	tenant, err := h.repo.GetByID(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "tenant not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  tenant,
		"message": "success",
	})
}

func (h *TenantHandler) List(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user not found in context",
		})
		return
	}

	tenants, err := h.repo.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to list tenants",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  tenants,
		"message": "success",
	})
}

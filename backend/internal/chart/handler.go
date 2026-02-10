package chart

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/auth"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	req.TenantID = auth.GetTenantID(c)
	if req.TenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	chart, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to create chart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": chart, "message": "chart created"})
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	req.TenantID = auth.GetTenantID(c)
	if req.TenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	chart, err := h.service.Update(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "chart not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": chart, "message": "chart updated"})
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
		return
	}

	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id, tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to delete chart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "chart deleted"})
}

func (h *Handler) Get(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
		return
	}

	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	chart, err := h.service.Get(c.Request.Context(), id, tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "chart not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": chart, "message": "success"})
}

func (h *Handler) List(c *gin.Context) {
	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	charts, err := h.service.List(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to list charts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": charts, "message": "success"})
}

func (h *Handler) Render(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
		return
	}

	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	config, err := h.service.Render(c.Request.Context(), id, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to render chart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": map[string]interface{}{"config": config}, "message": "success"})
}

package dashboard

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

func (h *Handler) List(c *gin.Context) {
	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	dashboards, err := h.service.List(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to list dashboards"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": dashboards, "message": "success"})
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

	dashboard, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to create dashboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": dashboard, "message": "dashboard created"})
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	req.ID = id
	req.TenantID = auth.GetTenantID(c)
	if req.TenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	dashboard, err := h.service.Update(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "dashboard not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": dashboard, "message": "dashboard updated"})
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
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
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to delete dashboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "dashboard deleted"})
}

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
		return
	}

	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	dashboard, err := h.service.Get(c.Request.Context(), id, tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "dashboard not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": dashboard, "message": "success"})
}

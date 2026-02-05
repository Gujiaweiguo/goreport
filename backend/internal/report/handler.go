package report

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeecg/jimureport-go/internal/auth"
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

	report, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to create report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": report, "message": "report created"})
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

	report, err := h.service.Update(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "report not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": report, "message": "report updated"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to delete report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "report deleted"})
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

	report, err := h.service.Get(c.Request.Context(), id, tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "report not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": report, "message": "success"})
}

func (h *Handler) List(c *gin.Context) {
	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	reports, err := h.service.List(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to list reports"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": reports, "message": "success"})
}

func (h *Handler) Preview(c *gin.Context) {
	var req PreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	req.TenantID = auth.GetTenantID(c)
	if req.TenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	resp, err := h.service.Preview(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to preview report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": resp, "message": "success"})
}

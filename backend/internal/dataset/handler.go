package dataset

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/auth"
)

type Handler struct {
	service       Service
	queryExecutor QueryExecutor
}

func NewHandler(service Service, queryExecutor QueryExecutor) *Handler {
	return &Handler{
		service:       service,
		queryExecutor: queryExecutor,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	req.TenantID = auth.GetTenantID(c)
	req.CreatedBy = auth.GetUserID(c)
	if req.TenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	dataset, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "result": dataset, "message": "dataset created"})
}

func (h *Handler) List(c *gin.Context) {
	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	datasets, total, err := h.service.List(c.Request.Context(), tenantID, pageInt, pageSizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to list datasets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"result":   datasets,
		"total":    total,
		"page":     pageInt,
		"pageSize": pageSizeInt,
	})
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

	dataset, err := h.service.GetWithFields(c.Request.Context(), id, tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "dataset not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": dataset, "message": "success"})
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

	dataset, err := h.service.Update(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": dataset, "message": "dataset updated"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to delete dataset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "dataset deleted"})
}

func (h *Handler) Preview(c *gin.Context) {
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

	data, err := h.service.Preview(c.Request.Context(), id, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to preview dataset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": data, "message": "success"})
}

func (h *Handler) QueryData(c *gin.Context) {
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

	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	req.DatasetID = id

	result, err := h.queryExecutor.Query(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to query dataset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": result, "message": "success"})
}

func (h *Handler) GetDimensions(c *gin.Context) {
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

	dimensions, err := h.service.ListDimensions(c.Request.Context(), id, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to get dimensions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": dimensions, "message": "success"})
}

func (h *Handler) GetMeasures(c *gin.Context) {
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

	measures, err := h.service.ListMeasures(c.Request.Context(), id, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to get measures"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": measures, "message": "success"})
}

func (h *Handler) GetSchema(c *gin.Context) {
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

	schema, err := h.service.GetSchema(c.Request.Context(), id, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to get schema"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": schema, "message": "success"})
}

func (h *Handler) CreateComputedField(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
		return
	}

	var req CreateFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	req.DatasetID = id
	req.TenantID = auth.GetTenantID(c)
	if req.TenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	field, err := h.service.CreateComputedField(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "result": field, "message": "computed field created"})
}

func (h *Handler) UpdateField(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
		return
	}

	fieldID := c.Param("fieldId")
	if fieldID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "fieldId is required"})
		return
	}

	var req UpdateFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	req.FieldID = fieldID
	req.TenantID = auth.GetTenantID(c)
	if req.TenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	field, err := h.service.UpdateField(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": field, "message": "field updated"})
}

func (h *Handler) DeleteField(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
		return
	}

	fieldID := c.Param("fieldId")
	if fieldID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "fieldId is required"})
		return
	}

	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	if err := h.service.DeleteField(c.Request.Context(), fieldID, tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to delete field"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "field deleted"})
}

package datasource

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/auth"
	"github.com/gujiaweiguo/goreport/internal/models"
)

type Handler struct {
	service  Service
	metadata *CachedMetadataService
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func NewHandlerWithMetadata(service Service, metadata *CachedMetadataService) *Handler {
	return &Handler{service: service, metadata: metadata}
}

var connectionBuilder = NewConnectionBuilder()

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
		return
	}

	ds, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	tenantID := auth.GetTenantID(c)
	if tenantID == "" || ds.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "permission denied"})
		return
	}

	ds.Password = ""
	c.JSON(http.StatusOK, gin.H{"success": true, "result": ds, "message": "success"})
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	req.TenantID = auth.GetTenantID(c)
	if req.TenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}
	if req.CreatedBy == "" {
		req.CreatedBy = auth.GetUserID(c)
		if req.CreatedBy == "" {
			req.CreatedBy = "system"
		}
	}

	ds, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "result": ds, "message": "datasource created"})
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

	datasources, total, err := h.service.List(c.Request.Context(), tenantID, pageInt, pageSizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result": gin.H{
			"datasources": datasources,
			"total":       total,
			"page":        pageInt,
			"pageSize":    pageSizeInt,
		},
		"total":    total,
		"page":     pageInt,
		"pageSize": pageSizeInt,
		"message":  "success",
	})
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

	ds, err := h.service.Update(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": ds, "message": "datasource updated"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "datasource deleted"})
}

func (h *Handler) Search(c *gin.Context) {
	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	keyword := c.Query("keyword")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	datasources, total, err := h.service.Search(c.Request.Context(), tenantID, keyword, pageInt, pageSizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result": gin.H{
			"datasources": datasources,
			"total":       total,
			"page":        pageInt,
			"pageSize":    pageSizeInt,
		},
		"total":    total,
		"page":     pageInt,
		"pageSize": pageSizeInt,
		"message":  "success",
	})
}

func (h *Handler) Copy(c *gin.Context) {
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

	ds, err := h.service.Copy(c.Request.Context(), id, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "result": ds, "message": "datasource copied"})
}

func (h *Handler) Move(c *gin.Context) {
	var req struct {
		ID     string `json:"id" binding:"required"`
		Target string `json:"target" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	if err := h.service.Move(c.Request.Context(), req.ID, tenantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "datasource moved"})
}

func (h *Handler) Rename(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required,max=255"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	ds, err := h.service.Rename(c.Request.Context(), id, tenantID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": ds, "message": "datasource renamed"})
}

func (h *Handler) GetTables(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id is required"})
		return
	}

	ds, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if ds.TenantID != auth.GetTenantID(c) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "permission denied"})
		return
	}

	if h.metadata == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"success": false, "message": "metadata service not available"})
		return
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		ds.Username,
		ds.Password,
		ds.Host,
		ds.Port,
		ds.Database,
	)

	tables, err := h.metadata.GetTables(c.Request.Context(), ds.TenantID, id, dsn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": tables, "message": "success"})
}

func (h *Handler) GetFields(c *gin.Context) {
	id := c.Param("id")
	tableName := c.Param("table")

	if id == "" || tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id and table are required"})
		return
	}

	ds, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if ds.TenantID != auth.GetTenantID(c) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "permission denied"})
		return
	}

	if h.metadata == nil {
		c.JSON(http.StatusNotImplemented, gin.H{"success": false, "message": "metadata service not available"})
		return
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		ds.Username,
		ds.Password,
		ds.Host,
		ds.Port,
		ds.Database,
	)

	fields, err := h.metadata.GetFields(c.Request.Context(), ds.TenantID, id, dsn, tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": fields, "message": "success"})
}

func (h *Handler) TestConnection(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
		return
	}

	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "tenant not found"})
		return
	}

	advanced := &AdvancedConfig{}
	if req.Advanced != nil {
		advanced = req.Advanced
	}

	testDS := &models.DataSource{
		Type:                req.Type,
		Host:                req.Host,
		Port:                req.Port,
		Database:            req.Database,
		Username:            req.Username,
		Password:            req.Password,
		SSHHost:             advanced.SSHHost,
		SSHPort:             advanced.SSHPort,
		SSHUsername:         advanced.SSHUsername,
		SSHPassword:         advanced.SSHPassword,
		SSHKey:              advanced.SSHKey,
		SSHKeyPhrase:        advanced.SSHKeyPhrase,
		MaxConnections:      advanced.MaxConnections,
		QueryTimeoutSeconds: advanced.QueryTimeoutSeconds,
	}

	if advanced.SSHHost == "" && advanced.SSHPort == 0 {
		testDS.SSHHost = ""
		testDS.SSHPort = 0
	}

	if err := connectionBuilder.TestConnection(c.Request.Context(), testDS); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": fmt.Sprintf("connection test failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "connection successful"})
}

func (h *Handler) TestSavedConnection(c *gin.Context) {
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

	ds, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if ds.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "permission denied"})
		return
	}

	if err := connectionBuilder.TestConnection(c.Request.Context(), ds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": fmt.Sprintf("connection test failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "connection successful"})
}

func (h *Handler) ListProfiles(c *gin.Context) {
	validator := NewProfileValidator()
	profiles := validator.ListProfiles()

	c.JSON(http.StatusOK, gin.H{"success": true, "result": profiles, "message": "success"})
}

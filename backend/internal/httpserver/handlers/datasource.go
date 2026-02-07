package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/auth"
	"github.com/gujiaweiguo/goreport/internal/cache"
	"github.com/gujiaweiguo/goreport/internal/datasource"
	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/gujiaweiguo/goreport/internal/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataSourceHandler struct {
	repo     repository.DataSourceRepository
	metadata *datasource.CachedMetadataService
	cache    *cache.Cache
}

func NewDataSourceHandler(db *gorm.DB, cache *cache.Cache) *DataSourceHandler {
	return &DataSourceHandler{
		repo:     repository.NewDataSourceRepository(db),
		metadata: datasource.NewCachedMetadataService(cache),
		cache:    cache,
	}
}

type CreateDataSourceRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Database string `json:"database" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TestDataSourceRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Database string `json:"database" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
}

type UpdateDataSourceRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *DataSourceHandler) CreateDatasource(c *gin.Context) {
	var req CreateDataSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
		})
		return
	}

	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "tenant not found in context",
		})
		return
	}

	ds := &models.DataSource{
		ID:           fmt.Sprintf("ds-%d", time.Now().UnixNano()),
		Name:         req.Name,
		Type:         req.Type,
		Host:         req.Host,
		Port:         req.Port,
		DatabaseName: req.Database,
		Username:     req.Username,
		Password:     req.Password,
		TenantID:     tenantID,
	}

	if err := h.repo.Create(c.Request.Context(), ds); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create datasource",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  ds,
		"message": "datasource created successfully",
	})
}

func (h *DataSourceHandler) ListDatasources(c *gin.Context) {
	tenantID := auth.GetTenantID(c)
	if tenantID == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "tenant not found in context",
		})
		return
	}

	datasources, err := h.repo.List(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to list datasources",
		})
		return
	}

	for _, ds := range datasources {
		ds.Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  datasources,
		"message": "success",
	})
}

func (h *DataSourceHandler) UpdateDatasource(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "datasource id is required",
		})
		return
	}

	existingDS, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "datasource not found",
		})
		return
	}

	tenantID := auth.GetTenantID(c)
	if tenantID == "" || existingDS.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "permission denied",
		})
		return
	}

	var req UpdateDataSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
		})
		return
	}

	if req.Name != "" {
		existingDS.Name = req.Name
	}
	if req.Type != "" {
		existingDS.Type = req.Type
	}
	if req.Host != "" {
		existingDS.Host = req.Host
	}
	if req.Port != 0 {
		existingDS.Port = req.Port
	}
	if req.Database != "" {
		existingDS.DatabaseName = req.Database
	}
	if req.Username != "" {
		existingDS.Username = req.Username
	}
	if req.Password != "" {
		existingDS.Password = req.Password
	}

	if err := h.repo.Update(c.Request.Context(), existingDS); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to update datasource",
		})
		return
	}

	if h.cache != nil {
		_ = h.cache.Invalidate(c.Request.Context(), tenantID, "datasource:tables")
	}

	existingDS.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  existingDS,
		"message": "datasource updated successfully",
	})
}

func (h *DataSourceHandler) DeleteDatasource(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "datasource id is required",
		})
		return
	}

	existingDS, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "datasource not found",
		})
		return
	}

	tenantID := auth.GetTenantID(c)
	if tenantID == "" || existingDS.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "permission denied",
		})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to delete datasource",
		})
		return
	}

	if h.cache != nil {
		_ = h.cache.Invalidate(c.Request.Context(), tenantID, "datasource:tables")
		_ = h.cache.Invalidate(c.Request.Context(), tenantID, "datasource:fields")
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "datasource deleted successfully",
	})
}

func (h *DataSourceHandler) TestDatasource(c *gin.Context) {
	var req TestDataSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request",
		})
		return
	}

	// 如果密码为空，从数据库中获取现有数据源的密码
	password := req.Password
	if password == "" {
		datasources, err := h.repo.List(c.Request.Context(), auth.GetTenantID(c))
		if err == nil {
			for _, ds := range datasources {
				if ds.Name == req.Name {
					password = ds.Password
					break
				}
			}
		}
		if password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "password is required",
			})
			return
		}
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		req.Username,
		password,
		req.Host,
		req.Port,
		req.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("connection failed: %v", err),
		})
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get database connection",
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("connection failed: %v", err),
		})
		return
	}

	sqlDB.Close()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "connection successful",
	})
}

func (h *DataSourceHandler) GetTables(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "datasource id is required",
		})
		return
	}

	ds, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "datasource not found",
		})
		return
	}

	tenantID := auth.GetTenantID(c)
	if tenantID == "" || ds.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "permission denied",
		})
		return
	}

	tables, err := h.metadata.GetTables(c.Request.Context(), tenantID, id, buildDSN(ds))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get tables",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  tables,
		"message": "success",
	})
}

func (h *DataSourceHandler) GetFields(c *gin.Context) {
	id := c.Param("id")
	tableName := c.Param("table")

	if id == "" || tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "datasource id and table name are required",
		})
		return
	}

	ds, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "datasource not found",
		})
		return
	}

	tenantID := auth.GetTenantID(c)
	if tenantID == "" || ds.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "permission denied",
		})
		return
	}

	fields, err := h.metadata.GetFields(c.Request.Context(), tenantID, id, buildDSN(ds), tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get fields",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  fields,
		"message": "success",
	})
}

func buildDSN(ds *models.DataSource) string {
	port := strconv.Itoa(ds.Port)
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		ds.Username,
		ds.Password,
		ds.Host,
		port,
		ds.DatabaseName,
	)
}

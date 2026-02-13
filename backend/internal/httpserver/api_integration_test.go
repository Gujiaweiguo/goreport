package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/config"
	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/gujiaweiguo/goreport/internal/report"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// setupTestServer 创建测试服务器
func setupTestServer(t *testing.T) (*Server, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	// 获取数据库连接
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	if dsn == "" {
		t.Skip("TEST_DB_DSN or DB_DSN not set")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移
	db.AutoMigrate(&models.User{}, &models.Tenant{}, &models.DataSource{}, &models.Dataset{}, &models.Dashboard{}, &report.Report{})

	cfg := &config.Config{
		Server: config.ServerConfig{
			Addr: ":8080",
		},
		JWT: config.JWTConfig{
			Secret:   "test-secret",
			Issuer:   "goreport",
			Audience: "goreport",
		},
		Cache: config.CacheConfig{
			Enabled:    true,
			DefaultTTL: 3600,
		},
	}

	server, err := NewServer(cfg, db)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	return server, db
}

// getTestToken 获取测试用的 JWT token
func getTestToken(t *testing.T, server *Server, username, password string) string {
	body := map[string]string{
		"username": username,
		"password": password,
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	server.Engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Login failed with status: %d, body: %s", w.Code, w.Body.String())
		return ""
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if result, ok := response["result"].(map[string]interface{}); ok {
		if token, ok := result["token"].(string); ok {
			return token
		}
	}

	return ""
}

// TestHealthEndpoint 测试健康检查端点
func TestAPI_HealthEndpoint(t *testing.T) {
	server, db := setupTestServer(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	server.Engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

// TestAuthEndpoints 测试认证端点
func TestAPI_AuthEndpoints(t *testing.T) {
	server, db := setupTestServer(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tests := []struct {
		name       string
		method     string
		path       string
		body       map[string]string
		expectCode int
	}{
		{
			name:       "login with empty credentials",
			method:     http.MethodPost,
			path:       "/api/v1/auth/login",
			body:       map[string]string{},
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "login with invalid credentials",
			method:     http.MethodPost,
			path:       "/api/v1/auth/login",
			body:       map[string]string{"username": "invalid", "password": "invalid"},
			expectCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			server.Engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectCode, w.Code)
		})
	}
}

// TestUserEndpoints 测试用户端点
func TestAPI_UserEndpoints(t *testing.T) {
	server, db := setupTestServer(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tests := []struct {
		name       string
		method     string
		path       string
		needAuth   bool
		expectCode int
	}{
		{
			name:       "get me without auth",
			method:     http.MethodGet,
			path:       "/api/v1/users/me",
			needAuth:   false,
			expectCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			server.Engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectCode, w.Code)
		})
	}
}

// TestTenantEndpoints 测试租户端点
func TestAPI_TenantEndpoints(t *testing.T) {
	server, db := setupTestServer(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tests := []struct {
		name       string
		method     string
		path       string
		expectCode int
	}{
		{
			name:       "list tenants without auth",
			method:     http.MethodGet,
			path:       "/api/v1/tenants",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "get current tenant without auth",
			method:     http.MethodGet,
			path:       "/api/v1/tenants/current",
			expectCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			server.Engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectCode, w.Code)
		})
	}
}

// TestDatasourceEndpoints 测试数据源端点
func TestAPI_DatasourceEndpoints(t *testing.T) {
	server, db := setupTestServer(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tests := []struct {
		name       string
		method     string
		path       string
		expectCode int
	}{
		{
			name:       "list datasources without auth",
			method:     http.MethodGet,
			path:       "/api/v1/datasources",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "create datasource without auth",
			method:     http.MethodPost,
			path:       "/api/v1/datasources",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "get datasource without auth",
			method:     http.MethodGet,
			path:       "/api/v1/datasources/123",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "test connection without auth",
			method:     http.MethodPost,
			path:       "/api/v1/datasources/test",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "list profiles without auth",
			method:     http.MethodGet,
			path:       "/api/v1/datasources/profiles",
			expectCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			server.Engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectCode, w.Code)
		})
	}
}

// TestDashboardEndpoints 测试仪表盘端点
func TestAPI_DashboardEndpoints(t *testing.T) {
	server, db := setupTestServer(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tests := []struct {
		name       string
		method     string
		path       string
		expectCode int
	}{
		{
			name:       "list dashboards without auth",
			method:     http.MethodGet,
			path:       "/api/v1/dashboard/list",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "create dashboard without auth",
			method:     http.MethodPost,
			path:       "/api/v1/dashboard/create",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "get dashboard without auth",
			method:     http.MethodGet,
			path:       "/api/v1/dashboard/123",
			expectCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			server.Engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectCode, w.Code)
		})
	}
}

// TestDatasetEndpoints 测试数据集端点
func TestAPI_DatasetEndpoints(t *testing.T) {
	server, db := setupTestServer(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tests := []struct {
		name       string
		method     string
		path       string
		expectCode int
	}{
		{
			name:       "list datasets without auth",
			method:     http.MethodGet,
			path:       "/api/v1/datasets",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "create dataset without auth",
			method:     http.MethodPost,
			path:       "/api/v1/datasets",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "get dataset without auth",
			method:     http.MethodGet,
			path:       "/api/v1/datasets/123",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "query dataset without auth",
			method:     http.MethodPost,
			path:       "/api/v1/datasets/123/data",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "get dimensions without auth",
			method:     http.MethodGet,
			path:       "/api/v1/datasets/123/dimensions",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "get measures without auth",
			method:     http.MethodGet,
			path:       "/api/v1/datasets/123/measures",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "get schema without auth",
			method:     http.MethodGet,
			path:       "/api/v1/datasets/123/schema",
			expectCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			server.Engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectCode, w.Code)
		})
	}
}

// TestReportEndpoints 测试报表端点
func TestAPI_ReportEndpoints(t *testing.T) {
	server, db := setupTestServer(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tests := []struct {
		name       string
		method     string
		path       string
		expectCode int
	}{
		{
			name:       "list reports without auth",
			method:     http.MethodGet,
			path:       "/api/v1/jmreport/list",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "create report without auth",
			method:     http.MethodPost,
			path:       "/api/v1/jmreport/create",
			expectCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			server.Engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectCode, w.Code)
		})
	}
}

// TestCacheEndpoints 测试缓存端点
func TestAPI_CacheEndpoints(t *testing.T) {
	server, db := setupTestServer(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tests := []struct {
		name       string
		method     string
		path       string
		expectCode int
	}{
		{
			name:       "get cache metrics without auth",
			method:     http.MethodGet,
			path:       "/api/v1/cache/metrics",
			expectCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			server.Engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectCode, w.Code)
		})
	}
}

// TestAPIEndpointsSummary API 端点汇总测试
func TestAPI_EndpointsSummary(t *testing.T) {
	server, db := setupTestServer(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 测试所有端点是否可访问（无认证）
	endpoints := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/health"},
		{http.MethodPost, "/api/v1/auth/login"},
		{http.MethodPost, "/api/v1/auth/logout"},
		{http.MethodGet, "/api/v1/users/me"},
		{http.MethodGet, "/api/v1/tenants"},
		{http.MethodGet, "/api/v1/tenants/current"},
		{http.MethodGet, "/api/v1/datasources"},
		{http.MethodPost, "/api/v1/datasources"},
		{http.MethodGet, "/api/v1/datasources/123"},
		{http.MethodPut, "/api/v1/datasources/123"},
		{http.MethodDelete, "/api/v1/datasources/123"},
		{http.MethodGet, "/api/v1/datasources/123/tables"},
		{http.MethodGet, "/api/v1/datasources/123/tables/users/fields"},
		{http.MethodPost, "/api/v1/datasources/copy/123"},
		{http.MethodPost, "/api/v1/datasources/move"},
		{http.MethodPut, "/api/v1/datasources/123/rename"},
		{http.MethodGet, "/api/v1/datasources/search"},
		{http.MethodPost, "/api/v1/datasources/test"},
		{http.MethodPost, "/api/v1/datasources/123/test"},
		{http.MethodGet, "/api/v1/datasources/profiles"},
		{http.MethodGet, "/api/v1/cache/metrics"},
		{http.MethodGet, "/api/v1/jmreport/list"},
		{http.MethodGet, "/api/v1/jmreport/get"},
		{http.MethodPost, "/api/v1/jmreport/create"},
		{http.MethodPost, "/api/v1/jmreport/update"},
		{http.MethodDelete, "/api/v1/jmreport/delete"},
		{http.MethodPost, "/api/v1/jmreport/preview"},
		{http.MethodGet, "/api/v1/dashboard/list"},
		{http.MethodPost, "/api/v1/dashboard/create"},
		{http.MethodGet, "/api/v1/dashboard/123"},
		{http.MethodPut, "/api/v1/dashboard/123"},
		{http.MethodDelete, "/api/v1/dashboard/123"},
		{http.MethodGet, "/api/v1/datasets"},
		{http.MethodPost, "/api/v1/datasets"},
		{http.MethodGet, "/api/v1/datasets/123"},
		{http.MethodPut, "/api/v1/datasets/123"},
		{http.MethodDelete, "/api/v1/datasets/123"},
		{http.MethodGet, "/api/v1/datasets/123/preview"},
		{http.MethodPost, "/api/v1/datasets/123/data"},
		{http.MethodGet, "/api/v1/datasets/123/dimensions"},
		{http.MethodGet, "/api/v1/datasets/123/measures"},
		{http.MethodGet, "/api/v1/datasets/123/schema"},
		{http.MethodPost, "/api/v1/datasets/123/fields"},
		{http.MethodPatch, "/api/v1/datasets/123/fields"},
		{http.MethodPut, "/api/v1/datasets/123/fields/456"},
		{http.MethodDelete, "/api/v1/datasets/123/fields/456"},
	}

	var passed, failed int
	for _, ep := range endpoints {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(ep.method, ep.path, nil)
		server.Engine.ServeHTTP(w, req)

		// 只要返回状态码不为 404，就认为端点存在
		if w.Code != http.StatusNotFound {
			passed++
			t.Logf("✅ %s %s - Status: %d", ep.method, ep.path, w.Code)
		} else {
			failed++
			t.Errorf("❌ %s %s - Not Found (404)", ep.method, ep.path)
		}
	}

	t.Logf("\n=== API 端点测试汇总 ===")
	t.Logf("总计: %d 个端点", len(endpoints))
	t.Logf("通过: %d", passed)
	t.Logf("失败: %d", failed)
}

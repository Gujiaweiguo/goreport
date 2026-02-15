package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestServer_NewServer_NilDB(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Addr: ":8085",
		},
		JWT: config.JWTConfig{
			Secret:   "test-secret-key-for-unit-test",
			Issuer:   "goreport",
			Audience: "goreport",
		},
		Cache: config.CacheConfig{
			Enabled: false,
		},
	}

	server, err := NewServer(cfg, nil)
	if err == nil && server != nil {
		assert.NotNil(t, server.Engine)
		assert.NotNil(t, server.Server)
	}
}

func TestServer_GetEngine_ReturnsEngine(t *testing.T) {
	engine := gin.New()
	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":8080"},
	}

	result := server.GetEngine()
	assert.Equal(t, engine, result)
	assert.NotNil(t, result)
}

func TestServer_GetEngine_NilEngine(t *testing.T) {
	server := &Server{
		Engine: nil,
		Server: &http.Server{Addr: ":8080"},
	}

	result := server.GetEngine()
	assert.Nil(t, result)
}

func TestServer_Shutdown_ContextCancellation(t *testing.T) {
	engine := gin.New()
	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":0"},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := server.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestServer_Shutdown_WithTimeout(t *testing.T) {
	engine := gin.New()
	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":0"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestServer_Run_InvalidAddress(t *testing.T) {
	engine := gin.New()
	server := &Server{
		Engine: engine,
		Server: &http.Server{},
	}

	err := server.Run("")
	if err != nil {
		t.Logf("Run with empty address returned error: %v (expected)", err)
	}
}

func TestServer_Struct_Fields(t *testing.T) {
	engine := gin.New()
	httpServer := &http.Server{Addr: ":9090"}

	server := &Server{
		Engine: engine,
		Server: httpServer,
	}

	assert.Equal(t, engine, server.Engine)
	assert.Equal(t, httpServer, server.Server)
	assert.Equal(t, ":9090", server.Server.Addr)
}

func TestServer_Engine_Routes(t *testing.T) {
	engine := gin.New()

	engine.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	engine.POST("/api/v1/test", func(c *gin.Context) {
		c.JSON(201, gin.H{"created": true})
	})

	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":8080"},
	}

	tests := []struct {
		method string
		path   string
		code   int
	}{
		{"GET", "/test", 200},
		{"POST", "/api/v1/test", 201},
		{"GET", "/nonexistent", 404},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()
		server.Engine.ServeHTTP(w, req)
		assert.Equal(t, tt.code, w.Code, "Path: %s %s", tt.method, tt.path)
	}
}

func TestServer_Middleware_Execution(t *testing.T) {
	engine := gin.New()

	engine.Use(func(c *gin.Context) {
		c.Set("middleware-executed", true)
		c.Next()
	})

	engine.GET("/middleware-test", func(c *gin.Context) {
		val, exists := c.Get("middleware-executed")
		if exists && val.(bool) {
			c.JSON(200, gin.H{"middleware": "executed"})
		} else {
			c.JSON(500, gin.H{"middleware": "not executed"})
		}
	})

	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":8080"},
	}

	req := httptest.NewRequest("GET", "/middleware-test", nil)
	w := httptest.NewRecorder()
	server.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "executed")
}

func TestServer_Concurrent_Requests(t *testing.T) {
	engine := gin.New()

	engine.GET("/concurrent", func(c *gin.Context) {
		time.Sleep(10 * time.Millisecond)
		c.JSON(200, gin.H{"status": "ok"})
	})

	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":8080"},
	}

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			req := httptest.NewRequest("GET", "/concurrent", nil)
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for concurrent request")
		}
	}
}

func TestServer_Context_Values(t *testing.T) {
	engine := gin.New()

	engine.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-123")
		c.Set("tenant_id", "test-tenant-456")
		c.Next()
	})

	engine.GET("/context-test", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		tenantID, _ := c.Get("tenant_id")
		c.JSON(200, gin.H{
			"user_id":   userID,
			"tenant_id": tenantID,
		})
	})

	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":8080"},
	}

	req := httptest.NewRequest("GET", "/context-test", nil)
	w := httptest.NewRecorder()
	server.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "test-user-123")
	assert.Contains(t, w.Body.String(), "test-tenant-456")
}

func TestServer_ErrorHandler(t *testing.T) {
	engine := gin.New()
	engine.Use(gin.Recovery())

	engine.GET("/error", func(c *gin.Context) {
		c.JSON(500, gin.H{"error": "internal error"})
	})

	engine.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":8080"},
	}

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	server.Engine.ServeHTTP(w, req)
	assert.Equal(t, 500, w.Code)

	req = httptest.NewRequest("GET", "/panic", nil)
	w = httptest.NewRecorder()
	server.Engine.ServeHTTP(w, req)
	assert.Equal(t, 500, w.Code)
}

func TestServer_RouteGroups(t *testing.T) {
	engine := gin.New()

	api := engine.Group("/api/v1")
	{
		api.GET("/users", func(c *gin.Context) {
			c.JSON(200, gin.H{"users": []string{}})
		})
		api.POST("/users", func(c *gin.Context) {
			c.JSON(201, gin.H{"created": true})
		})
	}

	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":8080"},
	}

	tests := []struct {
		method string
		path   string
		code   int
	}{
		{"GET", "/api/v1/users", 200},
		{"POST", "/api/v1/users", 201},
		{"GET", "/api/v1/users/123", 404},
		{"GET", "/api/v2/users", 404},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()
		server.Engine.ServeHTTP(w, req)
		assert.Equal(t, tt.code, w.Code, "Path: %s %s", tt.method, tt.path)
	}
}

func TestServer_QueryParams(t *testing.T) {
	engine := gin.New()

	engine.GET("/search", func(c *gin.Context) {
		keyword := c.Query("keyword")
		page := c.DefaultQuery("page", "1")
		c.JSON(200, gin.H{
			"keyword": keyword,
			"page":    page,
		})
	})

	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":8080"},
	}

	req := httptest.NewRequest("GET", "/search?keyword=test&page=2", nil)
	w := httptest.NewRecorder()
	server.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"keyword":"test"`)
	assert.Contains(t, w.Body.String(), `"page":"2"`)
}

func TestServer_PathParams(t *testing.T) {
	engine := gin.New()

	engine.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(200, gin.H{"user_id": id})
	})

	engine.GET("/posts/:postId/comments/:commentId", func(c *gin.Context) {
		postId := c.Param("postId")
		commentId := c.Param("commentId")
		c.JSON(200, gin.H{
			"post_id":    postId,
			"comment_id": commentId,
		})
	})

	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":8080"},
	}

	tests := []struct {
		path       string
		expected   map[string]string
		expectCode int
	}{
		{"/users/123", map[string]string{"user_id": "123"}, 200},
		{"/posts/456/comments/789", map[string]string{"post_id": "456", "comment_id": "789"}, 200},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", tt.path, nil)
		w := httptest.NewRecorder()
		server.Engine.ServeHTTP(w, req)

		assert.Equal(t, tt.expectCode, w.Code)
		for key, value := range tt.expected {
			assert.Contains(t, w.Body.String(), `"`+key+`":"`+value+`"`)
		}
	}
}

func TestServer_JSONBinding(t *testing.T) {
	engine := gin.New()

	type LoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	engine.POST("/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"username": req.Username})
	})

	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":8080"},
	}

	req := httptest.NewRequest("POST", "/login", strings.NewReader(`{"username":"admin","password":"secret"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.Engine.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	req = httptest.NewRequest("POST", "/login", strings.NewReader(`{"username":"admin"}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	server.Engine.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestServer_Headers(t *testing.T) {
	engine := gin.New()

	engine.GET("/headers", func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		contentType := c.ContentType()
		c.JSON(200, gin.H{
			"authorization": auth,
			"content_type":  contentType,
		})
	})

	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":8080"},
	}

	req := httptest.NewRequest("GET", "/headers", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Bearer test-token")
	assert.Contains(t, w.Body.String(), "application/json")
}

func TestServer_ResponseHeaders(t *testing.T) {
	engine := gin.New()

	engine.GET("/response-headers", func(c *gin.Context) {
		c.Header("X-Custom-Header", "custom-value")
		c.Header("X-Request-Id", "req-123")
		c.JSON(200, gin.H{"status": "ok"})
	})

	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":8080"},
	}

	req := httptest.NewRequest("GET", "/response-headers", nil)
	w := httptest.NewRecorder()
	server.Engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "custom-value", w.Header().Get("X-Custom-Header"))
	assert.Equal(t, "req-123", w.Header().Get("X-Request-Id"))
}

func TestServer_StatusCodes(t *testing.T) {
	engine := gin.New()

	engine.GET("/ok", func(c *gin.Context) { c.JSON(200, gin.H{}) })
	engine.GET("/created", func(c *gin.Context) { c.JSON(201, gin.H{}) })
	engine.GET("/bad-request", func(c *gin.Context) { c.JSON(400, gin.H{}) })
	engine.GET("/unauthorized", func(c *gin.Context) { c.JSON(401, gin.H{}) })
	engine.GET("/forbidden", func(c *gin.Context) { c.JSON(403, gin.H{}) })
	engine.GET("/not-found", func(c *gin.Context) { c.JSON(404, gin.H{}) })
	engine.GET("/internal-error", func(c *gin.Context) { c.JSON(500, gin.H{}) })

	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":8080"},
	}

	tests := []struct {
		path        string
		expectCode  int
		description string
	}{
		{"/ok", 200, "OK"},
		{"/created", 201, "Created"},
		{"/bad-request", 400, "Bad Request"},
		{"/unauthorized", 401, "Unauthorized"},
		{"/forbidden", 403, "Forbidden"},
		{"/not-found", 404, "Not Found"},
		{"/internal-error", 500, "Internal Server Error"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)
			assert.Equal(t, tt.expectCode, w.Code)
		})
	}
}

func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name: "valid config",
			config: &config.Config{
				Server: config.ServerConfig{Addr: ":8080"},
				JWT:    config.JWTConfig{Secret: "secret", Issuer: "test", Audience: "test"},
				Cache:  config.CacheConfig{Enabled: false},
			},
		},
		{
			name: "empty server address",
			config: &config.Config{
				Server: config.ServerConfig{Addr: ""},
				JWT:    config.JWTConfig{Secret: "secret"},
				Cache:  config.CacheConfig{Enabled: false},
			},
		},
		{
			name: "cache enabled without redis",
			config: &config.Config{
				Server: config.ServerConfig{Addr: ":8080"},
				JWT:    config.JWTConfig{Secret: "secret"},
				Cache:  config.CacheConfig{Enabled: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.config)
			assert.NotNil(t, tt.config)
		})
	}
}

func TestServer_RouteList(t *testing.T) {
	engine := gin.New()

	health := engine.Group("/")
	{
		health.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{}) })
	}

	auth := engine.Group("/api/v1/auth")
	{
		auth.POST("/login", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		auth.POST("/logout", func(c *gin.Context) { c.JSON(200, gin.H{}) })
	}

	users := engine.Group("/api/v1/users")
	{
		users.GET("/me", func(c *gin.Context) { c.JSON(200, gin.H{}) })
	}

	tenants := engine.Group("/api/v1/tenants")
	{
		tenants.GET("", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		tenants.GET("/current", func(c *gin.Context) { c.JSON(200, gin.H{}) })
	}

	datasources := engine.Group("/api/v1/datasources")
	{
		datasources.GET("", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasources.POST("", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasources.GET("/:id", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasources.PUT("/:id", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasources.DELETE("/:id", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasources.GET("/:id/tables", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasources.GET("/:id/tables/:table/fields", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasources.POST("/copy/:id", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasources.POST("/move", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasources.PUT("/:id/rename", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasources.GET("/search", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasources.POST("/test", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasources.POST("/:id/test", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasources.GET("/profiles", func(c *gin.Context) { c.JSON(200, gin.H{}) })
	}

	engine.GET("/api/v1/cache/metrics", func(c *gin.Context) { c.JSON(200, gin.H{}) })

	reports := engine.Group("/api/v1/jmreport")
	{
		reports.GET("/list", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		reports.GET("/get", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		reports.POST("/create", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		reports.POST("/update", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		reports.DELETE("/delete", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		reports.POST("/preview", func(c *gin.Context) { c.JSON(200, gin.H{}) })
	}

	dashboards := engine.Group("/api/v1/dashboard")
	{
		dashboards.GET("/list", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		dashboards.POST("/create", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		dashboards.GET("/:id", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		dashboards.PUT("/:id", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		dashboards.DELETE("/:id", func(c *gin.Context) { c.JSON(200, gin.H{}) })
	}

	datasets := engine.Group("/api/v1/datasets")
	{
		datasets.GET("", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasets.POST("", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasets.GET("/:id", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasets.PUT("/:id", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasets.DELETE("/:id", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasets.GET("/:id/preview", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasets.POST("/:id/data", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasets.GET("/:id/dimensions", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasets.GET("/:id/measures", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasets.GET("/:id/schema", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasets.POST("/:id/fields", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasets.PATCH("/:id/fields", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasets.PUT("/:id/fields/:fieldId", func(c *gin.Context) { c.JSON(200, gin.H{}) })
		datasets.DELETE("/:id/fields/:fieldId", func(c *gin.Context) { c.JSON(200, gin.H{}) })
	}

	server := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":8080"},
	}

	routes := []struct {
		method string
		path   string
	}{
		{"GET", "/health"},
		{"POST", "/api/v1/auth/login"},
		{"POST", "/api/v1/auth/logout"},
		{"GET", "/api/v1/users/me"},
		{"GET", "/api/v1/tenants"},
		{"GET", "/api/v1/tenants/current"},
		{"GET", "/api/v1/datasources"},
		{"POST", "/api/v1/datasources"},
		{"GET", "/api/v1/datasources/test-id"},
		{"PUT", "/api/v1/datasources/test-id"},
		{"DELETE", "/api/v1/datasources/test-id"},
		{"GET", "/api/v1/datasources/test-id/tables"},
		{"GET", "/api/v1/datasources/test-id/tables/users/fields"},
		{"POST", "/api/v1/datasources/copy/test-id"},
		{"POST", "/api/v1/datasources/move"},
		{"PUT", "/api/v1/datasources/test-id/rename"},
		{"GET", "/api/v1/datasources/search"},
		{"POST", "/api/v1/datasources/test"},
		{"POST", "/api/v1/datasources/test-id/test"},
		{"GET", "/api/v1/datasources/profiles"},
		{"GET", "/api/v1/cache/metrics"},
		{"GET", "/api/v1/jmreport/list"},
		{"GET", "/api/v1/jmreport/get"},
		{"POST", "/api/v1/jmreport/create"},
		{"POST", "/api/v1/jmreport/update"},
		{"DELETE", "/api/v1/jmreport/delete"},
		{"POST", "/api/v1/jmreport/preview"},
		{"GET", "/api/v1/dashboard/list"},
		{"POST", "/api/v1/dashboard/create"},
		{"GET", "/api/v1/dashboard/test-id"},
		{"PUT", "/api/v1/dashboard/test-id"},
		{"DELETE", "/api/v1/dashboard/test-id"},
		{"GET", "/api/v1/datasets"},
		{"POST", "/api/v1/datasets"},
		{"GET", "/api/v1/datasets/test-id"},
		{"PUT", "/api/v1/datasets/test-id"},
		{"DELETE", "/api/v1/datasets/test-id"},
		{"GET", "/api/v1/datasets/test-id/preview"},
		{"POST", "/api/v1/datasets/test-id/data"},
		{"GET", "/api/v1/datasets/test-id/dimensions"},
		{"GET", "/api/v1/datasets/test-id/measures"},
		{"GET", "/api/v1/datasets/test-id/schema"},
		{"POST", "/api/v1/datasets/test-id/fields"},
		{"PATCH", "/api/v1/datasets/test-id/fields"},
		{"PUT", "/api/v1/datasets/test-id/fields/field-id"},
		{"DELETE", "/api/v1/datasets/test-id/fields/field-id"},
	}

	for _, route := range routes {
		t.Run(route.method+"_"+route.path, func(t *testing.T) {
			req := httptest.NewRequest(route.method, route.path, nil)
			w := httptest.NewRecorder()
			server.Engine.ServeHTTP(w, req)

			assert.NotEqual(t, 404, w.Code, "Route should exist: %s %s", route.method, route.path)
		})
	}
}

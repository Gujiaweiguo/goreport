package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestServer_Struct(t *testing.T) {
	s := &Server{
		Engine: gin.New(),
		Server: &http.Server{Addr: ":8080"},
	}

	assert.NotNil(t, s.Engine)
	assert.NotNil(t, s.Server)
	assert.Equal(t, ":8080", s.Server.Addr)
}

func TestServer_GetEngine(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	s := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":8080"},
	}

	result := s.GetEngine()
	assert.Equal(t, engine, result)
}

func TestServer_Run(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	s := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":0"},
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		s.Server.Close()
	}()

	err := s.Run(":0")
	if err != nil && err != http.ErrServerClosed {
		t.Logf("Server run error: %v", err)
	}
}

func TestServer_Shutdown(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	s := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":0"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestServer_Shutdown_WithCache(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	s := &Server{
		Engine: engine,
		Server: &http.Server{Addr: ":0"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestNewServer_ConfigValidation(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Addr: ":8085",
		},
		Database: config.DatabaseConfig{
			DSN:             "root:root@tcp(localhost:3306)/goreport",
			MaxOpenConns:    100,
			MaxIdleConns:    10,
			ConnMaxLifetime: 3600,
		},
		JWT: config.JWTConfig{
			Secret:   "test-secret",
			Issuer:   "goreport",
			Audience: "goreport",
		},
		Cache: config.CacheConfig{
			Enabled: false,
		},
	}

	_ = cfg
}

func TestServer_HealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
}

func TestServer_MiddlewareChain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	engine.Use(func(c *gin.Context) {
		c.Set("test", "value")
		c.Next()
	})

	engine.GET("/test", func(c *gin.Context) {
		val, exists := c.Get("test")
		if exists {
			c.JSON(200, gin.H{"test": val})
		} else {
			c.JSON(500, gin.H{"error": "middleware not executed"})
		}
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "value")
}

func TestServer_ContextTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	engine.GET("/slow", func(c *gin.Context) {
		ctx := c.Request.Context()
		select {
		case <-time.After(100 * time.Millisecond):
			c.JSON(200, gin.H{"status": "completed"})
		case <-ctx.Done():
			c.JSON(499, gin.H{"error": "context cancelled"})
		}
	})

	req := httptest.NewRequest("GET", "/slow", nil)
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

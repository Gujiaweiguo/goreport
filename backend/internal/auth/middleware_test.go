package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/config"
	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	return r
}

func TestAuthMiddleware_AllowsPublicPath(t *testing.T) {
	r := setupAuthTestRouter()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_RejectsWithoutToken(t *testing.T) {
	r := setupAuthTestRouter()
	r.GET("/api/v1/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_RejectsInvalidToken(t *testing.T) {
	InitJWT(&config.JWTConfig{Secret: "test-secret", Issuer: "test", Audience: "test"})

	r := setupAuthTestRouter()
	r.GET("/api/v1/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_SetsClaimsInContext(t *testing.T) {
	InitJWT(&config.JWTConfig{Secret: "test-secret", Issuer: "test", Audience: "test"})

	user := &models.User{ID: "u-1", Username: "alice", Role: "admin", TenantID: "tenant-1"}
	token, err := GenerateToken(user)
	require.NoError(t, err)

	r := setupAuthTestRouter()
	r.GET("/api/v1/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"userId":   GetUserID(c),
			"username": GetUsername(c),
			"tenantId": GetTenantID(c),
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, `"userId":"u-1"`)
	assert.Contains(t, body, `"username":"alice"`)
	assert.Contains(t, body, `"tenantId":"tenant-1"`)
}

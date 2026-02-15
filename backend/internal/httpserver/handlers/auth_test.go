package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/auth"
	"github.com/gujiaweiguo/goreport/internal/config"
	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getTestDSNForAuth() string {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	return dsn
}

func skipIfNoDBForAuthHandler(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := getTestDSNForAuth()
	if dsn == "" {
		t.Skip("TEST_DB_DSN or DB_DSN not set")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	return db
}

func TestNewAuthHandler(t *testing.T) {
	handler := NewAuthHandler(nil)
	assert.NotNil(t, handler)
}

func TestAuthHandler_Login_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(nil)

	router := gin.New()
	router.POST("/login", handler.Login)

	body := `{"username":""}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request")
}

func TestAuthHandler_Login_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(nil)

	router := gin.New()
	router.POST("/login", handler.Login)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request")
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(nil)

	router := gin.New()
	router.POST("/login", handler.Login)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request")
}

func TestAuthHandler_Logout_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(nil)

	router := gin.New()
	router.POST("/logout", handler.Logout)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing authorization token")
}

func TestAuthHandler_Logout_EmptyAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(nil)

	router := gin.New()
	router.POST("/logout", handler.Logout)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.Header.Set("Authorization", "")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing authorization token")
}

func TestAuthHandler_Logout_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(nil)

	router := gin.New()
	router.POST("/logout", handler.Logout)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid or expired token")
}

func TestAuthHandler_Logout_TokenInQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(nil)

	router := gin.New()
	router.POST("/logout", handler.Logout)

	req := httptest.NewRequest(http.MethodPost, "/logout?token=invalid-token", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid or expired token")
}

func TestAuthHandler_Login_Success(t *testing.T) {
	db := skipIfNoDBForAuthHandler(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	auth.InitJWT(&config.JWTConfig{
		Secret:   "test-secret",
		Issuer:   "test",
		Audience: "test",
	})

	now := time.Now()
	ts := fmt.Sprintf("%08x", now.UnixNano()&0xFFFFFFFF)
	username := fmt.Sprintf("testuser-%s", ts)

	hashedPassword, err := auth.HashPassword("password123")
	require.NoError(t, err)

	user := &models.User{
		ID:        fmt.Sprintf("u-%s", ts),
		Username:  username,
		Password:  hashedPassword,
		TenantID:  "test-auth-handler",
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id = ?", user.ID)
	})

	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(db)

	router := gin.New()
	router.POST("/login", handler.Login)

	body := fmt.Sprintf(`{"username":"%s","password":"password123"}`, username)
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "login success")
	assert.Contains(t, w.Body.String(), "token")
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	db := skipIfNoDBForAuthHandler(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	auth.InitJWT(&config.JWTConfig{
		Secret:   "test-secret",
		Issuer:   "test",
		Audience: "test",
	})

	now := time.Now()
	ts := fmt.Sprintf("%08x", now.UnixNano()&0xFFFFFFFF)
	username := fmt.Sprintf("testuser-logout-%s", ts)

	hashedPassword, err := auth.HashPassword("password123")
	require.NoError(t, err)

	user := &models.User{
		ID:        fmt.Sprintf("u-logout-%s", ts),
		Username:  username,
		Password:  hashedPassword,
		TenantID:  "test-auth-handler",
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id = ?", user.ID)
	})

	token, err := auth.GenerateToken(user)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(db)

	router := gin.New()
	router.POST("/logout", handler.Logout)

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "logout success")
}

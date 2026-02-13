package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

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

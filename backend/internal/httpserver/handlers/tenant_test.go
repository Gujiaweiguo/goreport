package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTenantRepository struct {
	mock.Mock
}

func (m *mockTenantRepository) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func (m *mockTenantRepository) ListByUserID(ctx context.Context, userID string) ([]*models.Tenant, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Tenant), args.Error(1)
}

func TestNewTenantHandler(t *testing.T) {
	handler := NewTenantHandler(nil)
	assert.NotNil(t, handler)
}

func TestTenantHandler_GetCurrent_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &mockTenantRepository{}
	handler := NewTenantHandler(mockRepo)

	mockRepo.On("GetByID", mock.Anything, "tenant-1").Return(&models.Tenant{
		ID:   "tenant-1",
		Name: "Test Tenant",
	}, nil)

	router := gin.New()
	router.GET("/current", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.GetCurrent(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/current", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	assert.Contains(t, w.Body.String(), "Test Tenant")
	mockRepo.AssertExpectations(t)
}

func TestTenantHandler_GetCurrent_NoTenant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &mockTenantRepository{}
	handler := NewTenantHandler(mockRepo)

	router := gin.New()
	router.GET("/current", handler.GetCurrent)

	req := httptest.NewRequest(http.MethodGet, "/current", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "tenant not found in context")
}

func TestTenantHandler_GetCurrent_TenantNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &mockTenantRepository{}
	handler := NewTenantHandler(mockRepo)

	mockRepo.On("GetByID", mock.Anything, "tenant-1").Return(nil, errors.New("not found"))

	router := gin.New()
	router.GET("/current", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.GetCurrent(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/current", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "tenant not found")
	mockRepo.AssertExpectations(t)
}

func TestTenantHandler_List_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &mockTenantRepository{}
	handler := NewTenantHandler(mockRepo)

	mockRepo.On("ListByUserID", mock.Anything, "user-1").Return([]*models.Tenant{
		{ID: "tenant-1", Name: "Tenant 1"},
		{ID: "tenant-2", Name: "Tenant 2"},
	}, nil)

	router := gin.New()
	router.GET("/list", func(c *gin.Context) {
		c.Set("userId", "user-1")
		handler.List(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	assert.Contains(t, w.Body.String(), "Tenant 1")
	mockRepo.AssertExpectations(t)
}

func TestTenantHandler_List_NoUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &mockTenantRepository{}
	handler := NewTenantHandler(mockRepo)

	router := gin.New()
	router.GET("/list", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "user not found in context")
}

func TestTenantHandler_List_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &mockTenantRepository{}
	handler := NewTenantHandler(mockRepo)

	mockRepo.On("ListByUserID", mock.Anything, "user-1").Return(nil, errors.New("db error"))

	router := gin.New()
	router.GET("/list", func(c *gin.Context) {
		c.Set("userId", "user-1")
		handler.List(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to list tenants")
	mockRepo.AssertExpectations(t)
}

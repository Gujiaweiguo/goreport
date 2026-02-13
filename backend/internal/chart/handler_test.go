package chart

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockChartService struct {
	mock.Mock
}

func (m *mockChartService) Create(ctx context.Context, req *CreateRequest) (*models.Chart, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Chart), args.Error(1)
}

func (m *mockChartService) Update(ctx context.Context, req *UpdateRequest) (*models.Chart, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Chart), args.Error(1)
}

func (m *mockChartService) Delete(ctx context.Context, id, tenantID string) error {
	args := m.Called(ctx, id, tenantID)
	return args.Error(0)
}

func (m *mockChartService) Get(ctx context.Context, id, tenantID string) (*models.Chart, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Chart), args.Error(1)
}

func (m *mockChartService) List(ctx context.Context, tenantID string) ([]*models.Chart, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Chart), args.Error(1)
}

func (m *mockChartService) Render(ctx context.Context, id, tenantID string) (string, error) {
	args := m.Called(ctx, id, tenantID)
	return args.String(0), args.Error(1)
}

func setupChartTestHandler() (*Handler, *mockChartService) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockChartService{}
	handler := NewHandler(mockSvc)
	return handler, mockSvc
}

func TestChartNewHandler(t *testing.T) {
	handler := NewHandler(nil)
	assert.NotNil(t, handler)
}

func TestChartHandler_Create_Success(t *testing.T) {
	handler, mockSvc := setupChartTestHandler()

	mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(req *CreateRequest) bool {
		return req.Name == "Test Chart" && req.TenantID == "tenant-1"
	})).Return(&models.Chart{
		ID:       "chart-1",
		Name:     "Test Chart",
		TenantID: "tenant-1",
	}, nil)

	body := `{"name":"Test Chart","type":"bar","datasetId":"ds-1"}`
	router := gin.New()
	router.POST("/create", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Create(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	mockSvc.AssertExpectations(t)
}

func TestChartHandler_Create_NoTenant(t *testing.T) {
	handler, _ := setupChartTestHandler()

	body := `{"name":"Test Chart","type":"bar"}`
	router := gin.New()
	router.POST("/create", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestChartHandler_Create_InvalidRequest(t *testing.T) {
	handler, _ := setupChartTestHandler()

	body := `{"invalid json`
	router := gin.New()
	router.POST("/create", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Create(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChartHandler_Update_Success(t *testing.T) {
	handler, mockSvc := setupChartTestHandler()

	mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(req *UpdateRequest) bool {
		return req.ID == "chart-1" && req.TenantID == "tenant-1"
	})).Return(&models.Chart{
		ID:       "chart-1",
		Name:     "Updated Chart",
		TenantID: "tenant-1",
	}, nil)

	body := `{"id":"chart-1","name":"Updated Chart"}`
	router := gin.New()
	router.PUT("/update", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Update(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestChartHandler_Delete_Success(t *testing.T) {
	handler, mockSvc := setupChartTestHandler()

	mockSvc.On("Delete", mock.Anything, "chart-1", "tenant-1").Return(nil)

	router := gin.New()
	router.DELETE("/delete", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Delete(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/delete?id=chart-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestChartHandler_Delete_NoID(t *testing.T) {
	handler, _ := setupChartTestHandler()

	router := gin.New()
	router.DELETE("/delete", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Delete(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChartHandler_Get_Success(t *testing.T) {
	handler, mockSvc := setupChartTestHandler()

	mockSvc.On("Get", mock.Anything, "chart-1", "tenant-1").Return(&models.Chart{
		ID:       "chart-1",
		Name:     "Test Chart",
		TenantID: "tenant-1",
	}, nil)

	router := gin.New()
	router.GET("/get", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Get(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/get?id=chart-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestChartHandler_List_Success(t *testing.T) {
	handler, mockSvc := setupChartTestHandler()

	mockSvc.On("List", mock.Anything, "tenant-1").Return([]*models.Chart{
		{ID: "chart-1", Name: "Chart 1", TenantID: "tenant-1"},
		{ID: "chart-2", Name: "Chart 2", TenantID: "tenant-1"},
	}, nil)

	router := gin.New()
	router.GET("/list", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.List(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	mockSvc.AssertExpectations(t)
}

func TestChartHandler_List_Error(t *testing.T) {
	handler, mockSvc := setupChartTestHandler()

	mockSvc.On("List", mock.Anything, "tenant-1").Return(nil, errors.New("database error"))

	router := gin.New()
	router.GET("/list", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.List(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

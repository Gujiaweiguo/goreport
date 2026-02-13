package dashboard

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

type mockDashboardService struct {
	mock.Mock
}

func (m *mockDashboardService) Create(ctx context.Context, req *CreateRequest) (*models.Dashboard, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Dashboard), args.Error(1)
}

func (m *mockDashboardService) Update(ctx context.Context, req *UpdateRequest) (*models.Dashboard, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Dashboard), args.Error(1)
}

func (m *mockDashboardService) Delete(ctx context.Context, id, tenantID string) error {
	args := m.Called(ctx, id, tenantID)
	return args.Error(0)
}

func (m *mockDashboardService) Get(ctx context.Context, id, tenantID string) (*models.Dashboard, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Dashboard), args.Error(1)
}

func (m *mockDashboardService) List(ctx context.Context, tenantID string) ([]*models.Dashboard, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Dashboard), args.Error(1)
}

func setupDashboardTestHandler() (*Handler, *mockDashboardService) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockDashboardService{}
	handler := NewHandler(mockSvc)
	return handler, mockSvc
}

func TestDashboardNewHandler(t *testing.T) {
	handler := NewHandler(nil)
	assert.NotNil(t, handler)
}

func TestDashboardHandler_List_Success(t *testing.T) {
	handler, mockSvc := setupDashboardTestHandler()

	mockSvc.On("List", mock.Anything, "tenant-1").Return([]*models.Dashboard{
		{ID: "db-1", Name: "Dashboard 1", TenantID: "tenant-1"},
		{ID: "db-2", Name: "Dashboard 2", TenantID: "tenant-1"},
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

func TestDashboardHandler_List_NoTenant(t *testing.T) {
	handler, _ := setupDashboardTestHandler()

	router := gin.New()
	router.GET("/list", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDashboardHandler_Create_Success(t *testing.T) {
	handler, mockSvc := setupDashboardTestHandler()

	mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(req *CreateRequest) bool {
		return req.Name == "New Dashboard" && req.TenantID == "tenant-1"
	})).Return(&models.Dashboard{
		ID:       "db-new",
		Name:     "New Dashboard",
		TenantID: "tenant-1",
	}, nil)

	body := `{"name":"New Dashboard","code":"new-dash"}`
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
	mockSvc.AssertExpectations(t)
}

func TestDashboardHandler_Create_Error(t *testing.T) {
	handler, mockSvc := setupDashboardTestHandler()

	mockSvc.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))

	body := `{"name":"New Dashboard"}`
	router := gin.New()
	router.POST("/create", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Create(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDashboardHandler_Get_Success(t *testing.T) {
	handler, mockSvc := setupDashboardTestHandler()

	mockSvc.On("Get", mock.Anything, "db-1", "tenant-1").Return(&models.Dashboard{
		ID:       "db-1",
		Name:     "Test Dashboard",
		TenantID: "tenant-1",
	}, nil)

	router := gin.New()
	router.GET("/:id", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Get(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/db-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDashboardHandler_Delete_Success(t *testing.T) {
	handler, mockSvc := setupDashboardTestHandler()

	mockSvc.On("Delete", mock.Anything, "db-1", "tenant-1").Return(nil)

	router := gin.New()
	router.DELETE("/:id", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Delete(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/db-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDashboardHandler_Delete_Error(t *testing.T) {
	handler, mockSvc := setupDashboardTestHandler()

	mockSvc.On("Delete", mock.Anything, "db-1", "tenant-1").Return(errors.New("not found"))

	router := gin.New()
	router.DELETE("/:id", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Delete(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/db-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

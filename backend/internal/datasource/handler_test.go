package datasource

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

type mockService struct {
	mock.Mock
}

func (m *mockService) Create(ctx context.Context, req *CreateRequest) (*models.DataSource, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DataSource), args.Error(1)
}

func (m *mockService) GetByID(ctx context.Context, id string) (*models.DataSource, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DataSource), args.Error(1)
}

func (m *mockService) List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.DataSource, int64, error) {
	args := m.Called(ctx, tenantID, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.DataSource), args.Get(1).(int64), args.Error(2)
}

func (m *mockService) Update(ctx context.Context, req *UpdateRequest) (*models.DataSource, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DataSource), args.Error(1)
}

func (m *mockService) Delete(ctx context.Context, id, tenantID string) error {
	args := m.Called(ctx, id, tenantID)
	return args.Error(0)
}

func (m *mockService) Search(ctx context.Context, tenantID, keyword string, page, pageSize int) ([]*models.DataSource, int64, error) {
	args := m.Called(ctx, tenantID, keyword, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.DataSource), args.Get(1).(int64), args.Error(2)
}

func (m *mockService) Copy(ctx context.Context, id, tenantID string) (*models.DataSource, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DataSource), args.Error(1)
}

func (m *mockService) Move(ctx context.Context, id, tenantID string) error {
	args := m.Called(ctx, id, tenantID)
	return args.Error(0)
}

func (m *mockService) Rename(ctx context.Context, id, tenantID, newName string) (*models.DataSource, error) {
	args := m.Called(ctx, id, tenantID, newName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DataSource), args.Error(1)
}

func setupTestHandler() (*Handler, *mockService, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockService{}
	handler := NewHandler(mockSvc)
	router := gin.New()
	return handler, mockSvc, router
}

func performRequest(handler *Handler, mockSvc *mockService, method, path, body string, setupRoutes func(*gin.Engine, *Handler)) *httptest.ResponseRecorder {
	router := gin.New()
	setupRoutes(router, handler)

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestNewHandler(t *testing.T) {
	handler := NewHandler(nil)
	assert.NotNil(t, handler)
}

func TestHandler_Get_Success(t *testing.T) {
	handler, mockSvc, _ := setupTestHandler()

	mockSvc.On("GetByID", mock.Anything, "ds-1").Return(&models.DataSource{
		ID:       "ds-1",
		Name:     "Test DS",
		TenantID: "tenant-1",
		Type:     "mysql",
	}, nil)

	w := performRequest(handler, mockSvc, http.MethodGet, "/ds-1", "", func(r *gin.Engine, h *Handler) {
		r.GET("/:id", func(c *gin.Context) {
			c.Set("tenantId", "tenant-1")
			h.Get(c)
		})
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Get_NotFound(t *testing.T) {
	handler, mockSvc, _ := setupTestHandler()

	mockSvc.On("GetByID", mock.Anything, "not-exist").Return(nil, errors.New("not found"))

	w := performRequest(handler, mockSvc, http.MethodGet, "/not-exist", "", func(r *gin.Engine, h *Handler) {
		r.GET("/:id", func(c *gin.Context) {
			c.Set("tenantId", "tenant-1")
			h.Get(c)
		})
	})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Get_CrossTenant(t *testing.T) {
	handler, mockSvc, _ := setupTestHandler()

	mockSvc.On("GetByID", mock.Anything, "ds-1").Return(&models.DataSource{
		ID:       "ds-1",
		TenantID: "tenant-2",
	}, nil)

	w := performRequest(handler, mockSvc, http.MethodGet, "/ds-1", "", func(r *gin.Engine, h *Handler) {
		r.GET("/:id", func(c *gin.Context) {
			c.Set("tenantId", "tenant-1")
			h.Get(c)
		})
	})

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Create_Success(t *testing.T) {
	handler, mockSvc, _ := setupTestHandler()

	mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(req *CreateRequest) bool {
		return req.Name == "New DS" && req.TenantID == "tenant-1"
	})).Return(&models.DataSource{
		ID:       "ds-new",
		Name:     "New DS",
		TenantID: "tenant-1",
	}, nil)

	body := `{"name":"New DS","type":"mysql","host":"localhost","port":3306}`
	w := performRequest(handler, mockSvc, http.MethodPost, "/create", body, func(r *gin.Engine, h *Handler) {
		r.POST("/create", func(c *gin.Context) {
			c.Set("tenantId", "tenant-1")
			c.Set("userId", "user-1")
			h.Create(c)
		})
	})

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Create_NoTenant(t *testing.T) {
	handler, mockSvc, _ := setupTestHandler()

	body := `{"name":"New DS","type":"mysql"}`
	w := performRequest(handler, mockSvc, http.MethodPost, "/create", body, func(r *gin.Engine, h *Handler) {
		r.POST("/create", h.Create)
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_List_Success(t *testing.T) {
	handler, mockSvc, _ := setupTestHandler()

	mockSvc.On("List", mock.Anything, "tenant-1", 1, 10).Return([]*models.DataSource{
		{ID: "ds-1", Name: "DS1", TenantID: "tenant-1"},
		{ID: "ds-2", Name: "DS2", TenantID: "tenant-1"},
	}, int64(2), nil)

	w := performRequest(handler, mockSvc, http.MethodGet, "/list?page=1&pageSize=10", "", func(r *gin.Engine, h *Handler) {
		r.GET("/list", func(c *gin.Context) {
			c.Set("tenantId", "tenant-1")
			h.List(c)
		})
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Delete_Success(t *testing.T) {
	handler, mockSvc, _ := setupTestHandler()

	mockSvc.On("Delete", mock.Anything, "ds-1", "tenant-1").Return(nil)

	w := performRequest(handler, mockSvc, http.MethodDelete, "/ds-1", "", func(r *gin.Engine, h *Handler) {
		r.DELETE("/:id", func(c *gin.Context) {
			c.Set("tenantId", "tenant-1")
			h.Delete(c)
		})
	})

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Search_Success(t *testing.T) {
	handler, mockSvc, _ := setupTestHandler()

	mockSvc.On("Search", mock.Anything, "tenant-1", "mysql", 1, 10).Return([]*models.DataSource{
		{ID: "ds-1", Name: "MySQL DS", TenantID: "tenant-1"},
	}, int64(1), nil)

	w := performRequest(handler, mockSvc, http.MethodGet, "/search?keyword=mysql&page=1&pageSize=10", "", func(r *gin.Engine, h *Handler) {
		r.GET("/search", func(c *gin.Context) {
			c.Set("tenantId", "tenant-1")
			h.Search(c)
		})
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Copy_Success(t *testing.T) {
	handler, mockSvc, _ := setupTestHandler()

	mockSvc.On("Copy", mock.Anything, "ds-1", "tenant-1").Return(&models.DataSource{
		ID:       "ds-copy",
		Name:     "DS (副本)",
		TenantID: "tenant-1",
	}, nil)

	w := performRequest(handler, mockSvc, http.MethodPost, "/ds-1/copy", "", func(r *gin.Engine, h *Handler) {
		r.POST("/:id/copy", func(c *gin.Context) {
			c.Set("tenantId", "tenant-1")
			h.Copy(c)
		})
	})

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

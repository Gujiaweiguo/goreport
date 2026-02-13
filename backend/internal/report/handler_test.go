package report

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockReportService struct {
	mock.Mock
}

func (m *mockReportService) Create(ctx context.Context, req *CreateRequest) (*Report, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Report), args.Error(1)
}

func (m *mockReportService) Update(ctx context.Context, req *UpdateRequest) (*Report, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Report), args.Error(1)
}

func (m *mockReportService) Delete(ctx context.Context, id, tenantID string) error {
	args := m.Called(ctx, id, tenantID)
	return args.Error(0)
}

func (m *mockReportService) Get(ctx context.Context, id, tenantID string) (*Report, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Report), args.Error(1)
}

func (m *mockReportService) List(ctx context.Context, tenantID string) ([]*Report, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Report), args.Error(1)
}

func (m *mockReportService) Preview(ctx context.Context, req *PreviewRequest) (*PreviewResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PreviewResponse), args.Error(1)
}

func setupReportTestHandler() (*Handler, *mockReportService) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockReportService{}
	handler := NewHandler(mockSvc)
	return handler, mockSvc
}

func TestReportNewHandler(t *testing.T) {
	handler := NewHandler(nil)
	assert.NotNil(t, handler)
}

func TestReportHandler_Create_Success(t *testing.T) {
	handler, mockSvc := setupReportTestHandler()

	mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(req *CreateRequest) bool {
		return req.Name == "Test Report" && req.TenantID == "tenant-1"
	})).Return(&Report{
		ID:       "r-1",
		Name:     "Test Report",
		TenantID: "tenant-1",
	}, nil)

	body := `{"name":"Test Report","config":"{}"}`
	router := gin.New()
	router.POST("/create", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Set("userId", "user-1")
		c.Set("roles", []string{"admin"})
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

func TestReportHandler_Create_NoTenant(t *testing.T) {
	handler, _ := setupReportTestHandler()

	body := `{"name":"Test Report","config":"{}"}`
	router := gin.New()
	router.POST("/create", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestReportHandler_Create_InvalidRequest(t *testing.T) {
	handler, _ := setupReportTestHandler()

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

func TestReportHandler_Update_Success(t *testing.T) {
	handler, mockSvc := setupReportTestHandler()

	mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(req *UpdateRequest) bool {
		return req.ID == "r-1" && req.Name == "Updated Report" && req.TenantID == "tenant-1"
	})).Return(&Report{
		ID:       "r-1",
		Name:     "Updated Report",
		TenantID: "tenant-1",
	}, nil)

	body := `{"id":"r-1","name":"Updated Report"}`
	router := gin.New()
	router.PUT("/update", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Set("roles", []string{"admin"})
		handler.Update(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	mockSvc.AssertExpectations(t)
}

func TestReportHandler_Update_NotFound(t *testing.T) {
	handler, mockSvc := setupReportTestHandler()

	mockSvc.On("Update", mock.Anything, mock.Anything).Return(nil, ErrNotFound)

	body := `{"id":"not-exist","name":"Updated Report"}`
	router := gin.New()
	router.PUT("/update", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Set("roles", []string{"admin"})
		handler.Update(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestReportHandler_Delete_Success(t *testing.T) {
	handler, mockSvc := setupReportTestHandler()

	mockSvc.On("Delete", mock.Anything, "r-1", "tenant-1").Return(nil)

	router := gin.New()
	router.DELETE("/delete", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Set("roles", []string{"admin"})
		handler.Delete(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/delete?id=r-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestReportHandler_Delete_MissingID(t *testing.T) {
	handler, _ := setupReportTestHandler()

	router := gin.New()
	router.DELETE("/delete", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Set("roles", []string{"admin"})
		handler.Delete(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReportHandler_Get_Success(t *testing.T) {
	handler, mockSvc := setupReportTestHandler()

	mockSvc.On("Get", mock.Anything, "r-1", "tenant-1").Return(&Report{
		ID:       "r-1",
		Name:     "Test Report",
		TenantID: "tenant-1",
	}, nil)

	router := gin.New()
	router.GET("/get", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Get(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/get?id=r-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	mockSvc.AssertExpectations(t)
}

func TestReportHandler_Get_NotFound(t *testing.T) {
	handler, mockSvc := setupReportTestHandler()

	mockSvc.On("Get", mock.Anything, "not-exist", "tenant-1").Return(nil, ErrNotFound)

	router := gin.New()
	router.GET("/get", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Get(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/get?id=not-exist", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestReportHandler_List_Success(t *testing.T) {
	handler, mockSvc := setupReportTestHandler()

	mockSvc.On("List", mock.Anything, "tenant-1").Return([]*Report{
		{ID: "r-1", Name: "Report 1", TenantID: "tenant-1"},
		{ID: "r-2", Name: "Report 2", TenantID: "tenant-1"},
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

func TestReportHandler_List_NoTenant(t *testing.T) {
	handler, _ := setupReportTestHandler()

	router := gin.New()
	router.GET("/list", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestReportHandler_Preview_Success(t *testing.T) {
	handler, mockSvc := setupReportTestHandler()

	mockSvc.On("Preview", mock.Anything, mock.MatchedBy(func(req *PreviewRequest) bool {
		return req.ID == "r-1" && req.TenantID == "tenant-1"
	})).Return(&PreviewResponse{HTML: "<html>test</html>"}, nil)

	body := `{"id":"r-1","params":{}}`
	router := gin.New()
	router.POST("/preview", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Preview(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/preview", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	mockSvc.AssertExpectations(t)
}

func TestReportHandler_Preview_NoTenant(t *testing.T) {
	handler, _ := setupReportTestHandler()

	body := `{"id":"r-1"}`
	router := gin.New()
	router.POST("/preview", handler.Preview)

	req := httptest.NewRequest(http.MethodPost, "/preview", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestReportHandler_Preview_Error(t *testing.T) {
	handler, mockSvc := setupReportTestHandler()

	mockSvc.On("Preview", mock.Anything, mock.Anything).Return(nil, errors.New("preview failed"))

	body := `{"id":"r-1"}`
	router := gin.New()
	router.POST("/preview", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Preview(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/preview", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

package dashboard

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/auth"
	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewHandler(t *testing.T) {
	service := NewService(nil)
	handler := NewHandler(service)
	assert.NotNil(t, handler)
}

func TestNewRepository(t *testing.T) {
	repo := NewRepository(nil)
	assert.NotNil(t, repo)
}

type mockService struct {
	dashboard  *models.Dashboard
	dashboards []*models.Dashboard
	createErr  error
	updateErr  error
	deleteErr  error
	getErr     error
	listErr    error
}

func (m *mockService) Create(ctx context.Context, req *CreateRequest) (*models.Dashboard, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return m.dashboard, nil
}

func (m *mockService) Update(ctx context.Context, req *UpdateRequest) (*models.Dashboard, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	return m.dashboard, nil
}

func (m *mockService) Delete(ctx context.Context, id, tenantID string) error {
	return m.deleteErr
}

func (m *mockService) Get(ctx context.Context, id, tenantID string) (*models.Dashboard, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.dashboard, nil
}

func (m *mockService) List(ctx context.Context, tenantID string) ([]*models.Dashboard, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.dashboards, nil
}

func setTenantID(c *gin.Context, tenantID string) {
	c.Set(string(auth.TenantIDKey), tenantID)
}

func TestHandler_List_Success(t *testing.T) {
	service := &mockService{
		dashboards: []*models.Dashboard{
			{ID: "dashboard-1", Name: "Dashboard 1"},
		},
	}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/dashboards", nil)
	setTenantID(c, "tenant-1")

	handler.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestHandler_List_NoTenant(t *testing.T) {
	service := &mockService{}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/dashboards", nil)

	handler.List(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "tenant not found")
}

func TestHandler_Create_Success(t *testing.T) {
	service := &mockService{
		dashboard: &models.Dashboard{ID: "dashboard-1", Name: "Test Dashboard"},
	}
	handler := NewHandler(service)

	body := map[string]interface{}{
		"name": "Test Dashboard",
		"code": "test-dashboard",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/dashboards", bytes.NewReader(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	setTenantID(c, "tenant-1")

	handler.Create(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "dashboard created")
}

func TestHandler_Create_InvalidRequest(t *testing.T) {
	service := &mockService{}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/dashboards", bytes.NewReader([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")
	setTenantID(c, "tenant-1")

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request")
}

func TestHandler_Create_NoTenant(t *testing.T) {
	service := &mockService{}
	handler := NewHandler(service)

	body := map[string]interface{}{"name": "Test"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/dashboards", bytes.NewReader(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestHandler_Get_Success(t *testing.T) {
	service := &mockService{
		dashboard: &models.Dashboard{ID: "dashboard-1", Name: "Test Dashboard"},
	}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "dashboard-1"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/dashboards/dashboard-1", nil)
	setTenantID(c, "tenant-1")

	handler.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Get_NoID(t *testing.T) {
	service := &mockService{}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{}
	c.Request = httptest.NewRequest(http.MethodGet, "/dashboards/", nil)
	setTenantID(c, "tenant-1")

	handler.Get(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "id is required")
}

func TestHandler_Get_NotFound(t *testing.T) {
	service := &mockService{getErr: assert.AnError}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "not-exist"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/dashboards/not-exist", nil)
	setTenantID(c, "tenant-1")

	handler.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandler_Update_Success(t *testing.T) {
	service := &mockService{
		dashboard: &models.Dashboard{ID: "dashboard-1", Name: "Updated"},
	}
	handler := NewHandler(service)

	body := map[string]interface{}{"name": "Updated"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "dashboard-1"}}
	c.Request = httptest.NewRequest(http.MethodPut, "/dashboards/dashboard-1", bytes.NewReader(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	setTenantID(c, "tenant-1")

	handler.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Update_NoID(t *testing.T) {
	service := &mockService{}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{}
	c.Request = httptest.NewRequest(http.MethodPut, "/dashboards/", nil)
	setTenantID(c, "tenant-1")

	handler.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Delete_Success(t *testing.T) {
	service := &mockService{}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "dashboard-1"}}
	c.Request = httptest.NewRequest(http.MethodDelete, "/dashboards/dashboard-1", nil)
	setTenantID(c, "tenant-1")

	handler.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "dashboard deleted")
}

func TestHandler_Delete_NoID(t *testing.T) {
	service := &mockService{}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{}
	c.Request = httptest.NewRequest(http.MethodDelete, "/dashboards/", nil)
	setTenantID(c, "tenant-1")

	handler.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

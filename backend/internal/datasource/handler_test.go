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
)

type mockService struct {
	datasource  *models.DataSource
	datasources []*models.DataSource
	createErr   error
	getErr      error
	listErr     error
	updateErr   error
	deleteErr   error
	searchErr   error
	copyErr     error
	moveErr     error
	renameErr   error
}

func (m *mockService) Create(ctx context.Context, req *CreateRequest) (*models.DataSource, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return m.datasource, nil
}

func (m *mockService) GetByID(ctx context.Context, id string) (*models.DataSource, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.datasource != nil {
		return m.datasource, nil
	}
	return nil, errors.New("not found")
}

func (m *mockService) List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.DataSource, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	return m.datasources, int64(len(m.datasources)), nil
}

func (m *mockService) Update(ctx context.Context, req *UpdateRequest) (*models.DataSource, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	return m.datasource, nil
}

func (m *mockService) Delete(ctx context.Context, id, tenantID string) error {
	return m.deleteErr
}

func (m *mockService) Search(ctx context.Context, tenantID, keyword string, page, pageSize int) ([]*models.DataSource, int64, error) {
	if m.searchErr != nil {
		return nil, 0, m.searchErr
	}
	return m.datasources, int64(len(m.datasources)), nil
}

func (m *mockService) Copy(ctx context.Context, id, tenantID string) (*models.DataSource, error) {
	if m.copyErr != nil {
		return nil, m.copyErr
	}
	return m.datasource, nil
}

func (m *mockService) Move(ctx context.Context, id, tenantID string) error {
	return m.moveErr
}

func (m *mockService) Rename(ctx context.Context, id, tenantID string, newName string) (*models.DataSource, error) {
	if m.renameErr != nil {
		return nil, m.renameErr
	}
	if m.datasource != nil {
		m.datasource.Name = newName
		return m.datasource, nil
	}
	return nil, errors.New("not found")
}

func setupHandlerTest(t *testing.T, svc *mockService) (*Handler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(svc)
	router := gin.New()
	return handler, router
}

func TestNewDatasourceHandler(t *testing.T) {
	handler := NewHandler(nil)
	assert.NotNil(t, handler)
}

func TestNewHandlerWithMetadata(t *testing.T) {
	handler := NewHandlerWithMetadata(nil, nil)
	assert.NotNil(t, handler)
}

func TestDatasourceHandler_Get_Success(t *testing.T) {
	svc := &mockService{
		datasource: &models.DataSource{
			ID:       "ds-1",
			Name:     "Test DS",
			TenantID: "tenant-1",
		},
	}
	handler, router := setupHandlerTest(t, svc)

	router.GET("/:id", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Get(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/ds-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
}

func TestDatasourceHandler_Get_NoID(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.GET("/", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Get(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "id is required")
}

func TestDatasourceHandler_Get_ServiceError(t *testing.T) {
	svc := &mockService{getErr: errors.New("db error")}
	handler, router := setupHandlerTest(t, svc)

	router.GET("/:id", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Get(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/ds-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDatasourceHandler_Get_WrongTenant(t *testing.T) {
	svc := &mockService{
		datasource: &models.DataSource{
			ID:       "ds-1",
			TenantID: "tenant-2",
		},
	}
	handler, router := setupHandlerTest(t, svc)

	router.GET("/:id", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Get(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/ds-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "permission denied")
}

func TestDatasourceHandler_Create_Success(t *testing.T) {
	svc := &mockService{
		datasource: &models.DataSource{ID: "ds-1", Name: "Test"},
	}
	handler, router := setupHandlerTest(t, svc)

	router.POST("/create", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Set("userId", "user-1")
		handler.Create(c)
	})

	body := `{"name":"Test","type":"mysql","host":"localhost","port":3306,"database":"testdb"}`
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
}

func TestDatasourceHandler_Create_NoTenant(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.POST("/create", handler.Create)

	body := `{"name":"Test","type":"mysql","host":"localhost","port":3306}`
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "tenant not found")
}

func TestDatasourceHandler_Create_InvalidRequest(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.POST("/create", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Create(c)
	})

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDatasourceHandler_Create_ServiceError(t *testing.T) {
	svc := &mockService{createErr: errors.New("db error")}
	handler, router := setupHandlerTest(t, svc)

	router.POST("/create", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Create(c)
	})

	body := `{"name":"Test","type":"mysql","host":"localhost","port":3306,"database":"testdb"}`
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDatasourceHandler_List_Success(t *testing.T) {
	svc := &mockService{
		datasources: []*models.DataSource{
			{ID: "ds-1", Name: "DS1", TenantID: "tenant-1"},
		},
	}
	handler, router := setupHandlerTest(t, svc)

	router.GET("/list", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.List(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/list?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
}

func TestDatasourceHandler_List_NoTenant(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.GET("/list", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDatasourceHandler_List_ServiceError(t *testing.T) {
	svc := &mockService{listErr: errors.New("db error")}
	handler, router := setupHandlerTest(t, svc)

	router.GET("/list", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.List(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDatasourceHandler_Update_Success(t *testing.T) {
	svc := &mockService{
		datasource: &models.DataSource{ID: "ds-1", Name: "Updated", TenantID: "tenant-1"},
	}
	handler, router := setupHandlerTest(t, svc)

	router.PUT("/:id", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Update(c)
	})

	body := `{"id":"ds-1","name":"Updated","tenantId":"tenant-1"}`
	req := httptest.NewRequest(http.MethodPut, "/ds-1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
}

func TestDatasourceHandler_Update_NoID(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.PUT("/", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Update(c)
	})

	body := `{"name":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDatasourceHandler_Update_NoTenant(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.PUT("/:id", handler.Update)

	body := `{"id":"ds-1","name":"Updated","tenantId":"tenant-1"}`
	req := httptest.NewRequest(http.MethodPut, "/ds-1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDatasourceHandler_Update_InvalidRequest(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.PUT("/:id", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Update(c)
	})

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPut, "/ds-1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDatasourceHandler_Update_ServiceError(t *testing.T) {
	svc := &mockService{
		datasource: &models.DataSource{ID: "ds-1", TenantID: "tenant-1"},
		updateErr:  errors.New("db error"),
	}
	handler, router := setupHandlerTest(t, svc)

	router.PUT("/:id", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Update(c)
	})

	body := `{"id":"ds-1","tenantId":"tenant-1"}`
	req := httptest.NewRequest(http.MethodPut, "/ds-1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDatasourceHandler_Delete_Success(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.DELETE("/:id", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Delete(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/ds-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDatasourceHandler_Delete_NoID(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.DELETE("/", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Delete(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDatasourceHandler_Delete_NoTenant(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.DELETE("/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/ds-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDatasourceHandler_Delete_ServiceError(t *testing.T) {
	svc := &mockService{deleteErr: errors.New("db error")}
	handler, router := setupHandlerTest(t, svc)

	router.DELETE("/:id", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Delete(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/ds-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDatasourceHandler_Search_Success(t *testing.T) {
	svc := &mockService{
		datasources: []*models.DataSource{
			{ID: "ds-1", Name: "MySQL DS", TenantID: "tenant-1"},
		},
	}
	handler, router := setupHandlerTest(t, svc)

	router.GET("/search", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Search(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/search?keyword=mysql", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
}

func TestDatasourceHandler_Search_NoTenant(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.GET("/search", handler.Search)

	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDatasourceHandler_Search_ServiceError(t *testing.T) {
	svc := &mockService{searchErr: errors.New("db error")}
	handler, router := setupHandlerTest(t, svc)

	router.GET("/search", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Search(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/search?keyword=test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDatasourceHandler_Copy_Success(t *testing.T) {
	svc := &mockService{
		datasource: &models.DataSource{ID: "ds-2", Name: "Copy of DS"},
	}
	handler, router := setupHandlerTest(t, svc)

	router.POST("/:id/copy", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Copy(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/ds-1/copy", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
}

func TestDatasourceHandler_Copy_NoID(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.POST("/", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Copy(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDatasourceHandler_Copy_NoTenant(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.POST("/:id/copy", handler.Copy)

	req := httptest.NewRequest(http.MethodPost, "/ds-1/copy", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDatasourceHandler_Copy_ServiceError(t *testing.T) {
	svc := &mockService{copyErr: errors.New("db error")}
	handler, router := setupHandlerTest(t, svc)

	router.POST("/:id/copy", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Copy(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/ds-1/copy", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDatasourceHandler_Move_Success(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.POST("/move", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Move(c)
	})

	body := `{"id":"ds-1","target":"folder-1"}`
	req := httptest.NewRequest(http.MethodPost, "/move", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDatasourceHandler_Move_InvalidRequest(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.POST("/move", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Move(c)
	})

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/move", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDatasourceHandler_Move_NoTenant(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.POST("/move", handler.Move)

	body := `{"id":"ds-1","target":"folder-1"}`
	req := httptest.NewRequest(http.MethodPost, "/move", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDatasourceHandler_Move_ServiceError(t *testing.T) {
	svc := &mockService{moveErr: errors.New("db error")}
	handler, router := setupHandlerTest(t, svc)

	router.POST("/move", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Move(c)
	})

	body := `{"id":"ds-1","target":"folder-1"}`
	req := httptest.NewRequest(http.MethodPost, "/move", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDatasourceHandler_Rename_Success(t *testing.T) {
	svc := &mockService{
		datasource: &models.DataSource{ID: "ds-1", Name: "New Name", TenantID: "tenant-1"},
	}
	handler, router := setupHandlerTest(t, svc)

	router.PUT("/:id/rename", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Rename(c)
	})

	body := `{"name":"New Name"}`
	req := httptest.NewRequest(http.MethodPut, "/ds-1/rename", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDatasourceHandler_Rename_NoID(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.PUT("/rename", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Rename(c)
	})

	body := `{"name":"New Name"}`
	req := httptest.NewRequest(http.MethodPut, "/rename", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDatasourceHandler_Rename_InvalidRequest(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.PUT("/:id/rename", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Rename(c)
	})

	body := `{}`
	req := httptest.NewRequest(http.MethodPut, "/ds-1/rename", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDatasourceHandler_Rename_NoTenant(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.PUT("/:id/rename", handler.Rename)

	body := `{"name":"New Name"}`
	req := httptest.NewRequest(http.MethodPut, "/ds-1/rename", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDatasourceHandler_Rename_ServiceError(t *testing.T) {
	svc := &mockService{renameErr: errors.New("db error")}
	handler, router := setupHandlerTest(t, svc)

	router.PUT("/:id/rename", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Rename(c)
	})

	body := `{"name":"New Name"}`
	req := httptest.NewRequest(http.MethodPut, "/ds-1/rename", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDatasourceHandler_ListProfiles(t *testing.T) {
	svc := &mockService{}
	handler, router := setupHandlerTest(t, svc)

	router.GET("/profiles", handler.ListProfiles)

	req := httptest.NewRequest(http.MethodGet, "/profiles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
}

package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func getTestDSN() string {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	return dsn
}

type mockDataSourceRepo struct {
	mock.Mock
}

func (m *mockDataSourceRepo) Create(ctx context.Context, ds *models.DataSource) error {
	args := m.Called(ctx, ds)
	return args.Error(0)
}

func (m *mockDataSourceRepo) GetByID(ctx context.Context, id string) (*models.DataSource, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DataSource), args.Error(1)
}

func (m *mockDataSourceRepo) List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.DataSource, int64, error) {
	args := m.Called(ctx, tenantID, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.DataSource), int64Arg(args, 1), args.Error(2)
}

func (m *mockDataSourceRepo) Update(ctx context.Context, ds *models.DataSource) error {
	args := m.Called(ctx, ds)
	return args.Error(0)
}

func (m *mockDataSourceRepo) Delete(ctx context.Context, id, tenantID string) error {
	args := m.Called(ctx, id, tenantID)
	return args.Error(0)
}

func (m *mockDataSourceRepo) Search(ctx context.Context, tenantID, keyword string, page, pageSize int) ([]*models.DataSource, int64, error) {
	args := m.Called(ctx, tenantID, keyword, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.DataSource), int64Arg(args, 1), args.Error(2)
}

func (m *mockDataSourceRepo) Copy(ctx context.Context, id, tenantID string) (*models.DataSource, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DataSource), args.Error(1)
}

func (m *mockDataSourceRepo) Move(ctx context.Context, id, tenantID string) error {
	args := m.Called(ctx, id, tenantID)
	return args.Error(0)
}

func (m *mockDataSourceRepo) Rename(ctx context.Context, id, tenantID string, newName string) error {
	args := m.Called(ctx, id, tenantID, newName)
	return args.Error(0)
}

func int64Arg(args mock.Arguments, idx int) int64 {
	v := args.Get(idx)
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	default:
		return 0
	}
}

func newTestDataSourceHandler(repo *mockDataSourceRepo) *DataSourceHandler {
	return &DataSourceHandler{repo: repo}
}

func performRequest(t *testing.T, method, path, body, tenantID string, register func(r *gin.Engine, h *DataSourceHandler), h *DataSourceHandler) *httptest.ResponseRecorder {
	t.Helper()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if tenantID != "" {
			c.Set("tenantId", tenantID)
		}
		c.Next()
	})
	register(r, h)

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestNewDataSourceHandler(t *testing.T) {
	handler := NewDataSourceHandler(nil, nil)
	assert.NotNil(t, handler)
}

func TestCreateDatasource_RequiresTenant(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	body := `{"name":"ds","type":"mysql","host":"localhost","port":3306,"database":"goreport","username":"root","password":"root"}`
	w := performRequest(t, http.MethodPost, "/datasource/create", body, "", func(r *gin.Engine, h *DataSourceHandler) {
		r.POST("/datasource/create", h.CreateDatasource)
	}, h)

	assert.Equal(t, http.StatusForbidden, w.Code)
	repo.AssertExpectations(t)
}

func TestCreateDatasource_Success(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	repo.On("Create", mock.Anything, mock.MatchedBy(func(ds *models.DataSource) bool {
		return ds.Name == "ds" && ds.TenantID == "tenant-1" && ds.Database == "goreport"
	})).Return(nil).Once()

	body := `{"name":"ds","type":"mysql","host":"localhost","port":3306,"database":"goreport","username":"root","password":"root"}`
	w := performRequest(t, http.MethodPost, "/datasource/create", body, "tenant-1", func(r *gin.Engine, h *DataSourceHandler) {
		r.POST("/datasource/create", h.CreateDatasource)
	}, h)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	repo.AssertExpectations(t)
}

func TestListDatasources_TenantIsolation(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	repo.On("List", mock.Anything, "tenant-1", 1, 10).Return([]*models.DataSource{{
		ID:       "ds-1",
		Name:     "demo",
		TenantID: "tenant-1",
	}}, int64(1), nil).Once()

	w := performRequest(t, http.MethodGet, "/datasource/list", "", "tenant-1", func(r *gin.Engine, h *DataSourceHandler) {
		r.GET("/datasource/list", h.ListDatasources)
	}, h)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	repo.AssertExpectations(t)
}

func TestUpdateDatasource_DeniesCrossTenant(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	repo.On("GetByID", mock.Anything, "ds-1").Return(&models.DataSource{
		ID:       "ds-1",
		Name:     "demo",
		TenantID: "tenant-2",
	}, nil).Once()

	body := `{"name":"new-name"}`
	w := performRequest(t, http.MethodPut, "/datasource/ds-1", body, "tenant-1", func(r *gin.Engine, h *DataSourceHandler) {
		r.PUT("/datasource/:id", h.UpdateDatasource)
	}, h)

	assert.Equal(t, http.StatusForbidden, w.Code)
	repo.AssertExpectations(t)
}

func TestDeleteDatasource_Success(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	repo.On("GetByID", mock.Anything, "ds-1").Return(&models.DataSource{
		ID:       "ds-1",
		Name:     "demo",
		TenantID: "tenant-1",
	}, nil).Once()
	repo.On("Delete", mock.Anything, "ds-1", "tenant-1").Return(nil).Once()

	w := performRequest(t, http.MethodDelete, "/datasource/ds-1", "", "tenant-1", func(r *gin.Engine, h *DataSourceHandler) {
		r.DELETE("/datasource/:id", h.DeleteDatasource)
	}, h)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}

func TestListDatasources_RequiresTenant(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	w := performRequest(t, http.MethodGet, "/datasource/list", "", "", func(r *gin.Engine, h *DataSourceHandler) {
		r.GET("/datasource/list", h.ListDatasources)
	}, h)

	assert.Equal(t, http.StatusForbidden, w.Code)
	repo.AssertExpectations(t)
}

func TestCreateDatasource_InvalidRequest(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	w := performRequest(t, http.MethodPost, "/datasource/create", `{}`, "tenant-1", func(r *gin.Engine, h *DataSourceHandler) {
		r.POST("/datasource/create", h.CreateDatasource)
	}, h)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	repo.AssertExpectations(t)
}

func TestListDatasources_ResponseBodyShape(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	repo.On("List", mock.Anything, "tenant-1", 1, 10).Return([]*models.DataSource{}, int64(0), nil).Once()

	w := performRequest(t, http.MethodGet, "/datasource/list", "", "tenant-1", func(r *gin.Engine, h *DataSourceHandler) {
		r.GET("/datasource/list", h.ListDatasources)
	}, h)
	require.Equal(t, http.StatusOK, w.Code)

	type response struct {
		Success bool `json:"success"`
	}

	var resp response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
	repo.AssertExpectations(t)
}

func TestBuildDSN(t *testing.T) {
	ds := &models.DataSource{
		Username: "root",
		Password: "password",
		Host:     "localhost",
		Port:     3306,
		Database: "testdb",
	}

	dsn := buildDSN(ds)
	assert.Contains(t, dsn, "root:password")
	assert.Contains(t, dsn, "localhost:3306")
	assert.Contains(t, dsn, "testdb")
}

func TestUpdateDatasource_Success(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	repo.On("GetByID", mock.Anything, "ds-1").Return(&models.DataSource{
		ID:       "ds-1",
		Name:     "old-name",
		TenantID: "tenant-1",
		Host:     "localhost",
		Port:     3306,
		Database: "olddb",
	}, nil).Once()
	repo.On("Update", mock.Anything, mock.MatchedBy(func(ds *models.DataSource) bool {
		return ds.Name == "new-name" && ds.Database == "newdb"
	})).Return(nil).Once()

	body := `{"name":"new-name","database":"newdb"}`
	w := performRequest(t, http.MethodPut, "/datasource/ds-1", body, "tenant-1", func(r *gin.Engine, h *DataSourceHandler) {
		r.PUT("/datasource/:id", h.UpdateDatasource)
	}, h)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	repo.AssertExpectations(t)
}

func TestUpdateDatasource_NotFound(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	repo.On("GetByID", mock.Anything, "ds-1").Return(nil, assert.AnError).Once()

	body := `{"name":"new-name"}`
	w := performRequest(t, http.MethodPut, "/datasource/ds-1", body, "tenant-1", func(r *gin.Engine, h *DataSourceHandler) {
		r.PUT("/datasource/:id", h.UpdateDatasource)
	}, h)

	assert.Equal(t, http.StatusNotFound, w.Code)
	repo.AssertExpectations(t)
}

func TestUpdateDatasource_InvalidRequest(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	repo.On("GetByID", mock.Anything, "ds-1").Return(&models.DataSource{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}, nil).Once()

	body := `{invalid json}`
	w := performRequest(t, http.MethodPut, "/datasource/ds-1", body, "tenant-1", func(r *gin.Engine, h *DataSourceHandler) {
		r.PUT("/datasource/:id", h.UpdateDatasource)
	}, h)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	repo.AssertExpectations(t)
}

func TestUpdateDatasource_UpdateError(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	repo.On("GetByID", mock.Anything, "ds-1").Return(&models.DataSource{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}, nil).Once()
	repo.On("Update", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	body := `{"name":"new-name"}`
	w := performRequest(t, http.MethodPut, "/datasource/ds-1", body, "tenant-1", func(r *gin.Engine, h *DataSourceHandler) {
		r.PUT("/datasource/:id", h.UpdateDatasource)
	}, h)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	repo.AssertExpectations(t)
}

func TestDeleteDatasource_NotFound(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	repo.On("GetByID", mock.Anything, "ds-1").Return(nil, assert.AnError).Once()

	w := performRequest(t, http.MethodDelete, "/datasource/ds-1", "", "tenant-1", func(r *gin.Engine, h *DataSourceHandler) {
		r.DELETE("/datasource/:id", h.DeleteDatasource)
	}, h)

	assert.Equal(t, http.StatusNotFound, w.Code)
	repo.AssertExpectations(t)
}

func TestDeleteDatasource_CrossTenant(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	repo.On("GetByID", mock.Anything, "ds-1").Return(&models.DataSource{
		ID:       "ds-1",
		TenantID: "tenant-2",
	}, nil).Once()

	w := performRequest(t, http.MethodDelete, "/datasource/ds-1", "", "tenant-1", func(r *gin.Engine, h *DataSourceHandler) {
		r.DELETE("/datasource/:id", h.DeleteDatasource)
	}, h)

	assert.Equal(t, http.StatusForbidden, w.Code)
	repo.AssertExpectations(t)
}

func TestDeleteDatasource_DeleteError(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	repo.On("GetByID", mock.Anything, "ds-1").Return(&models.DataSource{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}, nil).Once()
	repo.On("Delete", mock.Anything, "ds-1", "tenant-1").Return(assert.AnError).Once()

	w := performRequest(t, http.MethodDelete, "/datasource/ds-1", "", "tenant-1", func(r *gin.Engine, h *DataSourceHandler) {
		r.DELETE("/datasource/:id", h.DeleteDatasource)
	}, h)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	repo.AssertExpectations(t)
}

func TestCreateDatasources_RepoError(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	repo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	body := `{"name":"ds","type":"mysql","host":"localhost","port":3306,"database":"goreport","username":"root","password":"root"}`
	w := performRequest(t, http.MethodPost, "/datasource/create", body, "tenant-1", func(r *gin.Engine, h *DataSourceHandler) {
		r.POST("/datasource/create", h.CreateDatasource)
	}, h)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	repo.AssertExpectations(t)
}

func TestListDatasources_Error(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := newTestDataSourceHandler(repo)

	repo.On("List", mock.Anything, "tenant-1", 1, 10).Return(nil, int64(0), assert.AnError).Once()

	w := performRequest(t, http.MethodGet, "/datasource/list", "", "tenant-1", func(r *gin.Engine, h *DataSourceHandler) {
		r.GET("/datasource/list", h.ListDatasources)
	}, h)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	repo.AssertExpectations(t)
}

type mockMetadataProvider struct {
	mock.Mock
}

func (m *mockMetadataProvider) GetTables(ctx context.Context, tenantID, datasourceID, dsn string) ([]string, error) {
	args := m.Called(ctx, tenantID, datasourceID, dsn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockMetadataProvider) GetFields(ctx context.Context, tenantID, datasourceID, dsn, tableName string) ([]map[string]interface{}, error) {
	args := m.Called(ctx, tenantID, datasourceID, dsn, tableName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func TestGetTables_MissingID(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := &DataSourceHandler{repo: repo}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Next()
	})
	r.GET("/datasource/:id/tables", h.GetTables)

	req := httptest.NewRequest(http.MethodGet, "/datasource//tables", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetTables_NotFound(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := &DataSourceHandler{repo: repo}

	repo.On("GetByID", mock.Anything, "ds-1").Return(nil, assert.AnError).Once()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Next()
	})
	r.GET("/datasource/:id/tables", h.GetTables)

	req := httptest.NewRequest(http.MethodGet, "/datasource/ds-1/tables", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	repo.AssertExpectations(t)
}

func TestGetTables_CrossTenant(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := &DataSourceHandler{repo: repo}

	repo.On("GetByID", mock.Anything, "ds-1").Return(&models.DataSource{
		ID:       "ds-1",
		TenantID: "tenant-2",
	}, nil).Once()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Next()
	})
	r.GET("/datasource/:id/tables", h.GetTables)

	req := httptest.NewRequest(http.MethodGet, "/datasource/ds-1/tables", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	repo.AssertExpectations(t)
}

func TestGetFields_MissingID(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := &DataSourceHandler{repo: repo}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Next()
	})
	r.GET("/datasource/:id/fields/:table", h.GetFields)

	req := httptest.NewRequest(http.MethodGet, "/datasource//fields/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetFields_MissingTable(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := &DataSourceHandler{repo: repo}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Next()
	})
	r.GET("/datasource/:id/fields/:table", h.GetFields)

	req := httptest.NewRequest(http.MethodGet, "/datasource/ds-1/fields/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestGetFields_NotFound(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := &DataSourceHandler{repo: repo}

	repo.On("GetByID", mock.Anything, "ds-1").Return(nil, assert.AnError).Once()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Next()
	})
	r.GET("/datasource/:id/fields/:table", h.GetFields)

	req := httptest.NewRequest(http.MethodGet, "/datasource/ds-1/fields/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	repo.AssertExpectations(t)
}

func TestGetFields_CrossTenant(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := &DataSourceHandler{repo: repo}

	repo.On("GetByID", mock.Anything, "ds-1").Return(&models.DataSource{
		ID:       "ds-1",
		TenantID: "tenant-2",
	}, nil).Once()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Next()
	})
	r.GET("/datasource/:id/fields/:table", h.GetFields)

	req := httptest.NewRequest(http.MethodGet, "/datasource/ds-1/fields/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	repo.AssertExpectations(t)
}

func TestTestDatasource_InvalidRequest(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := &DataSourceHandler{repo: repo}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Next()
	})
	r.POST("/datasource/test", h.TestDatasource)

	req := httptest.NewRequest(http.MethodPost, "/datasource/test", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTestDatasource_MissingFields(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := &DataSourceHandler{repo: repo}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Next()
	})
	r.POST("/datasource/test", h.TestDatasource)

	req := httptest.NewRequest(http.MethodPost, "/datasource/test", strings.NewReader(`{"host": "localhost"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetFields_DatasourceNotFound(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := &DataSourceHandler{repo: repo}

	repo.On("GetByID", mock.Anything, "nonexistent").Return(nil, assert.AnError).Once()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Next()
	})
	r.GET("/datasource/:id/fields/:table", h.GetFields)

	req := httptest.NewRequest(http.MethodGet, "/datasource/nonexistent/fields/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	repo.AssertExpectations(t)
}

func TestUpdateDatasource_EmptyID(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := &DataSourceHandler{repo: repo}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Next()
	})
	r.PUT("/datasource/:id", h.UpdateDatasource)

	req := httptest.NewRequest(http.MethodPut, "/datasource/", strings.NewReader(`{"name": "test"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteDatasource_EmptyID(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := &DataSourceHandler{repo: repo}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Next()
	})
	r.DELETE("/datasource/:id", h.DeleteDatasource)

	req := httptest.NewRequest(http.MethodDelete, "/datasource/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTestDatasource_DB_Integration(t *testing.T) {
	dsn := getTestDSN()
	if dsn == "" {
		t.Skip("TEST_DB_DSN or DB_DSN not set")
	}

	repo := &mockDataSourceRepo{}
	h := &DataSourceHandler{repo: repo}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Next()
	})
	r.POST("/datasource/test", h.TestDatasource)

	body := `{"name":"test","type":"mysql","host":"127.0.0.1","port":3306,"username":"root","password":"root","database":"goreport"}`
	req := httptest.NewRequest(http.MethodPost, "/datasource/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestTestDatasource_DB_InvalidHost(t *testing.T) {
	repo := &mockDataSourceRepo{}
	h := &DataSourceHandler{repo: repo}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Next()
	})
	r.POST("/datasource/test", h.TestDatasource)

	body := `{"name":"test","type":"mysql","host":"invalid-host-xyz","port":3306,"username":"root","password":"root","database":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/datasource/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "connection failed")
}

package dataset

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

type mockDatasetService struct {
	mock.Mock
}

func (m *mockDatasetService) Create(ctx context.Context, req *CreateRequest) (*models.Dataset, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Dataset), args.Error(1)
}

func (m *mockDatasetService) Get(ctx context.Context, id, tenantID string) (*models.Dataset, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Dataset), args.Error(1)
}

func (m *mockDatasetService) GetWithFields(ctx context.Context, id, tenantID string) (*models.Dataset, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Dataset), args.Error(1)
}

func (m *mockDatasetService) List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.Dataset, int64, error) {
	args := m.Called(ctx, tenantID, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.Dataset), args.Get(1).(int64), args.Error(2)
}

func (m *mockDatasetService) Update(ctx context.Context, req *UpdateRequest) (*models.Dataset, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Dataset), args.Error(1)
}

func (m *mockDatasetService) Delete(ctx context.Context, id, tenantID string) error {
	args := m.Called(ctx, id, tenantID)
	return args.Error(0)
}

func (m *mockDatasetService) Preview(ctx context.Context, id, tenantID string) ([]map[string]interface{}, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *mockDatasetService) GetSchema(ctx context.Context, id, tenantID string) (*SchemaResponse, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SchemaResponse), args.Error(1)
}

func (m *mockDatasetService) CreateComputedField(ctx context.Context, req *CreateFieldRequest) (*models.DatasetField, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DatasetField), args.Error(1)
}

func (m *mockDatasetService) UpdateField(ctx context.Context, req *UpdateFieldRequest) (*models.DatasetField, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DatasetField), args.Error(1)
}

func (m *mockDatasetService) BatchUpdateFields(ctx context.Context, datasetID, tenantID string, req *BatchUpdateFieldsRequest) (*BatchUpdateFieldsResponse, error) {
	args := m.Called(ctx, datasetID, tenantID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BatchUpdateFieldsResponse), args.Error(1)
}

func (m *mockDatasetService) DeleteField(ctx context.Context, fieldID, tenantID string) error {
	args := m.Called(ctx, fieldID, tenantID)
	return args.Error(0)
}

func (m *mockDatasetService) ListDimensions(ctx context.Context, datasetID, tenantID string) ([]*models.DatasetField, error) {
	args := m.Called(ctx, datasetID, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DatasetField), args.Error(1)
}

func (m *mockDatasetService) ListMeasures(ctx context.Context, datasetID, tenantID string) ([]*models.DatasetField, error) {
	args := m.Called(ctx, datasetID, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DatasetField), args.Error(1)
}

func (m *mockDatasetService) ListFields(ctx context.Context, datasetID, tenantID string) ([]*models.DatasetField, error) {
	args := m.Called(ctx, datasetID, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DatasetField), args.Error(1)
}

type mockQueryExecutor struct {
	mock.Mock
}

func (m *mockQueryExecutor) Query(ctx context.Context, req *QueryRequest) (*QueryResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*QueryResponse), args.Error(1)
}

func setupDatasetTestHandler() (*Handler, *mockDatasetService, *mockQueryExecutor) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockDatasetService{}
	mockExec := &mockQueryExecutor{}
	handler := NewHandler(mockSvc, mockExec)
	return handler, mockSvc, mockExec
}

func TestDatasetNewHandler(t *testing.T) {
	handler := NewHandler(nil, nil)
	assert.NotNil(t, handler)
}

func TestDatasetHandler_Create_Success(t *testing.T) {
	handler, mockSvc, _ := setupDatasetTestHandler()

	mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(req *CreateRequest) bool {
		return req.Name == "Test Dataset" && req.TenantID == "tenant-1"
	})).Return(&models.Dataset{
		ID:       "ds-1",
		Name:     "Test Dataset",
		TenantID: "tenant-1",
	}, nil)

	body := `{"name":"Test Dataset","type":"sql"}`
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

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDatasetHandler_Create_NoTenant(t *testing.T) {
	handler, _, _ := setupDatasetTestHandler()

	body := `{"name":"Test Dataset","type":"sql"}`
	router := gin.New()
	router.POST("/create", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDatasetHandler_List_Success(t *testing.T) {
	handler, mockSvc, _ := setupDatasetTestHandler()

	mockSvc.On("List", mock.Anything, "tenant-1", 1, 10).Return([]*models.Dataset{
		{ID: "ds-1", Name: "Dataset 1", TenantID: "tenant-1"},
		{ID: "ds-2", Name: "Dataset 2", TenantID: "tenant-1"},
	}, int64(2), nil)

	router := gin.New()
	router.GET("/list", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.List(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/list?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	mockSvc.AssertExpectations(t)
}

func TestDatasetHandler_List_NoTenant(t *testing.T) {
	handler, _, _ := setupDatasetTestHandler()

	router := gin.New()
	router.GET("/list", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDatasetHandler_Get_Success(t *testing.T) {
	handler, mockSvc, _ := setupDatasetTestHandler()

	mockSvc.On("GetWithFields", mock.Anything, "ds-1", "tenant-1").Return(&models.Dataset{
		ID:       "ds-1",
		Name:     "Test Dataset",
		TenantID: "tenant-1",
	}, nil)

	router := gin.New()
	router.GET("/:id", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Get(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/ds-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDatasetHandler_Get_NotFound(t *testing.T) {
	handler, mockSvc, _ := setupDatasetTestHandler()

	mockSvc.On("GetWithFields", mock.Anything, "not-exist", "tenant-1").Return(nil, errors.New("not found"))

	router := gin.New()
	router.GET("/:id", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Get(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/not-exist", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDatasetHandler_Delete_Success(t *testing.T) {
	handler, mockSvc, _ := setupDatasetTestHandler()

	mockSvc.On("Delete", mock.Anything, "ds-1", "tenant-1").Return(nil)

	router := gin.New()
	router.DELETE("/:id", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Set("roles", []string{"admin"})
		handler.Delete(c)
	})

	req := httptest.NewRequest(http.MethodDelete, "/ds-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDatasetHandler_GetSchema_Success(t *testing.T) {
	handler, mockSvc, _ := setupDatasetTestHandler()

	dimName := "ID"
	measureName := "Amount"
	mockSvc.On("GetSchema", mock.Anything, "ds-1", "tenant-1").Return(&SchemaResponse{
		Dimensions: []*models.DatasetField{{ID: "f1", Name: "id", DisplayName: &dimName}},
		Measures:   []*models.DatasetField{{ID: "f2", Name: "amount", DisplayName: &measureName}},
		Computed:   []*models.DatasetField{},
	}, nil)

	router := gin.New()
	router.GET("/:id/schema", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.GetSchema(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/ds-1/schema", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDatasetHandler_Update_Success(t *testing.T) {
	handler, mockSvc, _ := setupDatasetTestHandler()

	mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(req *UpdateRequest) bool {
		return req.ID == "ds-1" && *req.Name == "Updated Name"
	})).Return(&models.Dataset{
		ID:       "ds-1",
		Name:     "Updated Name",
		TenantID: "tenant-1",
	}, nil)

	body := `{"name":"Updated Name"}`
	router := gin.New()
	router.PUT("/:id", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		c.Set("roles", []string{"admin"})
		handler.Update(c)
	})

	req := httptest.NewRequest(http.MethodPut, "/ds-1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDatasetHandler_Update_NoTenant(t *testing.T) {
	handler, _, _ := setupDatasetTestHandler()

	body := `{"name":"Updated Name"}`
	router := gin.New()
	router.PUT("/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/ds-1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDatasetHandler_Preview_Success(t *testing.T) {
	handler, mockSvc, _ := setupDatasetTestHandler()

	mockSvc.On("Preview", mock.Anything, "ds-1", "tenant-1").Return([]map[string]interface{}{
		{"id": 1, "name": "test"},
	}, nil)

	router := gin.New()
	router.GET("/:id/preview", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Preview(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/ds-1/preview", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDatasetHandler_Preview_Error(t *testing.T) {
	handler, mockSvc, _ := setupDatasetTestHandler()

	mockSvc.On("Preview", mock.Anything, "ds-1", "tenant-1").Return(nil, errors.New("preview failed"))

	router := gin.New()
	router.GET("/:id/preview", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.Preview(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/ds-1/preview", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDatasetHandler_Preview_NoTenant(t *testing.T) {
	handler, _, _ := setupDatasetTestHandler()

	router := gin.New()
	router.GET("/:id/preview", handler.Preview)

	req := httptest.NewRequest(http.MethodGet, "/ds-1/preview", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDatasetHandler_QueryData_Success(t *testing.T) {
	handler, _, mockExec := setupDatasetTestHandler()

	mockExec.On("Query", mock.Anything, mock.Anything).Return(&QueryResponse{
		Data:     []map[string]interface{}{{"id": 1, "name": "test"}},
		Total:    1,
		Page:     1,
		PageSize: 10,
	}, nil)

	body := `{"fields":["id","name"]}`
	router := gin.New()
	router.POST("/:id/query", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.QueryData(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/ds-1/query", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockExec.AssertExpectations(t)
}

func TestDatasetHandler_QueryData_Error(t *testing.T) {
	handler, _, mockExec := setupDatasetTestHandler()

	mockExec.On("Query", mock.Anything, mock.Anything).Return(nil, errors.New("query failed"))

	body := `{"fields":["id","name"]}`
	router := gin.New()
	router.POST("/:id/query", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.QueryData(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/ds-1/query", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockExec.AssertExpectations(t)
}

func TestDatasetHandler_QueryData_NoTenant(t *testing.T) {
	handler, _, _ := setupDatasetTestHandler()

	body := `{"fields":["id","name"]}`
	router := gin.New()
	router.POST("/:id/query", handler.QueryData)

	req := httptest.NewRequest(http.MethodPost, "/ds-1/query", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDatasetHandler_GetDimensions_Success(t *testing.T) {
	handler, mockSvc, _ := setupDatasetTestHandler()

	dimName := "Name"
	mockSvc.On("ListDimensions", mock.Anything, "ds-1", "tenant-1").Return([]*models.DatasetField{
		{ID: "f1", Name: "name", Type: "dimension", DisplayName: &dimName},
	}, nil)

	router := gin.New()
	router.GET("/:id/dimensions", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.GetDimensions(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/ds-1/dimensions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDatasetHandler_GetDimensions_Error(t *testing.T) {
	handler, mockSvc, _ := setupDatasetTestHandler()

	mockSvc.On("ListDimensions", mock.Anything, "ds-1", "tenant-1").Return(nil, errors.New("not found"))

	router := gin.New()
	router.GET("/:id/dimensions", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.GetDimensions(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/ds-1/dimensions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDatasetHandler_GetMeasures_Success(t *testing.T) {
	handler, mockSvc, _ := setupDatasetTestHandler()

	measureName := "Amount"
	mockSvc.On("ListMeasures", mock.Anything, "ds-1", "tenant-1").Return([]*models.DatasetField{
		{ID: "f1", Name: "amount", Type: "measure", DisplayName: &measureName},
	}, nil)

	router := gin.New()
	router.GET("/:id/measures", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.GetMeasures(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/ds-1/measures", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDatasetHandler_GetMeasures_Error(t *testing.T) {
	handler, mockSvc, _ := setupDatasetTestHandler()

	mockSvc.On("ListMeasures", mock.Anything, "ds-1", "tenant-1").Return(nil, errors.New("not found"))

	router := gin.New()
	router.GET("/:id/measures", func(c *gin.Context) {
		c.Set("tenantId", "tenant-1")
		handler.GetMeasures(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/ds-1/measures", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}

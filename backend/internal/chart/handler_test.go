package chart

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/auth"
	"github.com/gujiaweiguo/goreport/internal/dataset"
	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewHandler(t *testing.T) {
	handler := NewHandler(nil)
	assert.NotNil(t, handler)
}

func TestNewRepository(t *testing.T) {
	repo := NewRepository(nil)
	assert.NotNil(t, repo)
}

type mockChartService struct {
	chart        *models.Chart
	charts       []*models.Chart
	createErr    error
	updateErr    error
	deleteErr    error
	getErr       error
	listErr      error
	renderErr    error
	renderResult string
}

func (m *mockChartService) Create(ctx context.Context, req *CreateRequest) (*models.Chart, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return m.chart, nil
}

func (m *mockChartService) Update(ctx context.Context, req *UpdateRequest) (*models.Chart, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	return m.chart, nil
}

func (m *mockChartService) Delete(ctx context.Context, id, tenantID string) error {
	return m.deleteErr
}

func (m *mockChartService) Get(ctx context.Context, id, tenantID string) (*models.Chart, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.chart, nil
}

func (m *mockChartService) List(ctx context.Context, tenantID string) ([]*models.Chart, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.charts, nil
}

func (m *mockChartService) Render(ctx context.Context, id, tenantID string) (string, error) {
	if m.renderErr != nil {
		return "", m.renderErr
	}
	return m.renderResult, nil
}

func setChartTenantID(c *gin.Context, tenantID string) {
	c.Set(string(auth.TenantIDKey), tenantID)
}

func TestHandler_Create_Success(t *testing.T) {
	service := &mockChartService{
		chart: &models.Chart{ID: "chart-1", Name: "Test Chart"},
	}
	handler := NewHandler(service)

	body := map[string]interface{}{
		"name": "Test Chart",
		"type": "bar",
		"config": map[string]interface{}{
			"series": []map[string]interface{}{
				{"name": "Series 1", "type": "bar", "data": []int{1, 2, 3}},
			},
		},
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/charts", bytes.NewReader(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	setChartTenantID(c, "tenant-1")

	handler.Create(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "chart created")
}

func TestHandler_Create_InvalidRequest(t *testing.T) {
	service := &mockChartService{}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/charts", bytes.NewReader([]byte("invalid")))
	c.Request.Header.Set("Content-Type", "application/json")
	setChartTenantID(c, "tenant-1")

	handler.Create(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Create_NoTenant(t *testing.T) {
	service := &mockChartService{}
	handler := NewHandler(service)

	body := map[string]interface{}{"name": "Test", "type": "bar", "config": map[string]interface{}{}}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/charts", bytes.NewReader(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestHandler_Update_Success(t *testing.T) {
	service := &mockChartService{
		chart: &models.Chart{ID: "chart-1", Name: "Updated"},
	}
	handler := NewHandler(service)

	body := map[string]interface{}{"id": "chart-1", "name": "Updated"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/charts", bytes.NewReader(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	setChartTenantID(c, "tenant-1")

	handler.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Update_NotFound(t *testing.T) {
	service := &mockChartService{updateErr: ErrNotFound}
	handler := NewHandler(service)

	body := map[string]interface{}{"id": "not-exist", "name": "Updated"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/charts", bytes.NewReader(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	setChartTenantID(c, "tenant-1")

	handler.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandler_Delete_Success(t *testing.T) {
	service := &mockChartService{}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/charts?id=chart-1", nil)
	setChartTenantID(c, "tenant-1")

	handler.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "chart deleted")
}

func TestHandler_Delete_NoID(t *testing.T) {
	service := &mockChartService{}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/charts", nil)
	setChartTenantID(c, "tenant-1")

	handler.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Get_Success(t *testing.T) {
	service := &mockChartService{
		chart: &models.Chart{ID: "chart-1", Name: "Test Chart"},
	}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/charts?id=chart-1", nil)
	setChartTenantID(c, "tenant-1")

	handler.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Get_NoID(t *testing.T) {
	service := &mockChartService{}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/charts", nil)
	setChartTenantID(c, "tenant-1")

	handler.Get(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Get_NotFound(t *testing.T) {
	service := &mockChartService{getErr: ErrNotFound}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/charts?id=not-exist", nil)
	setChartTenantID(c, "tenant-1")

	handler.Get(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandler_List_Success(t *testing.T) {
	service := &mockChartService{
		charts: []*models.Chart{
			{ID: "chart-1", Name: "Chart 1"},
		},
	}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/charts/list", nil)
	setChartTenantID(c, "tenant-1")

	handler.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_List_NoTenant(t *testing.T) {
	service := &mockChartService{}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/charts/list", nil)

	handler.List(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestHandler_Render_Success(t *testing.T) {
	service := &mockChartService{
		renderResult: `{"series":[{"name":"Series1","data":[1,2,3]}]}`,
	}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/charts/render?id=chart-1", nil)
	setChartTenantID(c, "tenant-1")

	handler.Render(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Render_NoID(t *testing.T) {
	service := &mockChartService{}
	handler := NewHandler(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/charts/render", nil)
	setChartTenantID(c, "tenant-1")

	handler.Render(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChartConfig_Marshal(t *testing.T) {
	config := ChartConfig{
		Title: "Test Chart",
		XAxis: &AxisConfig{
			Type: "category",
			Data: []string{"A", "B", "C"},
			Name: "X",
		},
		YAxis: &AxisConfig{
			Type: "value",
			Name: "Y",
		},
		Series: []SeriesConfig{
			{
				Name: "Series 1",
				Type: "bar",
				Data: []any{10, 20, 30},
			},
		},
		Params: map[string]interface{}{
			"param1": "value1",
		},
	}

	jsonBytes, err := json.Marshal(config)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonBytes), "Test Chart")
	assert.Contains(t, string(jsonBytes), "Series 1")
}

func TestService_Update_CodeAndType(t *testing.T) {
	queryExec := &mockQueryExecutor{}
	repo := &mockRepository{}
	svc := NewService(repo, queryExec)

	existingChart := &models.Chart{ID: "c1", Name: "Old", Code: "old-code", Type: "bar", TenantID: "tenant-1"}
	repo.On("Get", mock.Anything, "c1", "tenant-1").Return(existingChart, nil).Once()
	repo.On("Update", mock.Anything, mock.MatchedBy(func(c *models.Chart) bool {
		return c.Code == "new-code" && c.Type == "line"
	})).Return(nil).Once()

	chart, err := svc.Update(context.Background(), &UpdateRequest{
		TenantID: "tenant-1",
		ID:       "c1",
		Code:     "new-code",
		Type:     "line",
		Config: ChartConfig{
			Series: []SeriesConfig{{Name: "S1", Type: "line", Data: []any{1}}},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "new-code", chart.Code)
	assert.Equal(t, "line", chart.Type)
	repo.AssertExpectations(t)
}

func TestService_Render_NoDatasetID(t *testing.T) {
	queryExec := &mockQueryExecutor{}
	repo := &mockRepository{}
	svc := NewService(repo, queryExec)

	configJSON, _ := json.Marshal(ChartConfig{
		Series: []SeriesConfig{
			{
				Name: "Series 1",
				Type: "bar",
				Data: []any{1, 2, 3},
			},
		},
	})

	existingChart := &models.Chart{
		ID:       "c1",
		Name:     "Test Chart",
		TenantID: "tenant-1",
		Config:   string(configJSON),
	}

	repo.On("Get", mock.Anything, "c1", "tenant-1").Return(existingChart, nil).Once()

	result, err := svc.Render(context.Background(), "c1", "tenant-1")
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Series 1")

	repo.AssertExpectations(t)
	queryExec.AssertNotCalled(t, "Query")
}

func TestService_Render_InvalidJSON(t *testing.T) {
	queryExec := &mockQueryExecutor{}
	repo := &mockRepository{}
	svc := NewService(repo, queryExec)

	existingChart := &models.Chart{
		ID:       "c1",
		Name:     "Test Chart",
		TenantID: "tenant-1",
		Config:   "invalid json",
	}

	repo.On("Get", mock.Anything, "c1", "tenant-1").Return(existingChart, nil).Once()

	_, err := svc.Render(context.Background(), "c1", "tenant-1")
	assert.Error(t, err)

	repo.AssertExpectations(t)
}

func TestDatasetQueryRequest(t *testing.T) {
	req := dataset.QueryRequest{
		DatasetID: "ds-1",
		Page:      1,
		PageSize:  100,
	}
	assert.Equal(t, "ds-1", req.DatasetID)
	assert.Equal(t, 1, req.Page)
	assert.Equal(t, 100, req.PageSize)
}

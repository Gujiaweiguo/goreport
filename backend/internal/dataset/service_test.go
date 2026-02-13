package dataset

import (
	"context"
	"errors"
	"testing"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDatasetRepository struct {
	mock.Mock
}

func (m *mockDatasetRepository) Create(ctx context.Context, dataset *models.Dataset) error {
	args := m.Called(ctx, dataset)
	return args.Error(0)
}

func (m *mockDatasetRepository) GetByID(ctx context.Context, id string) (*models.Dataset, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Dataset), args.Error(1)
}

func (m *mockDatasetRepository) GetByIDWithFields(ctx context.Context, id string) (*models.Dataset, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Dataset), args.Error(1)
}

func (m *mockDatasetRepository) List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.Dataset, int64, error) {
	args := m.Called(ctx, tenantID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.Dataset), args.Get(1).(int64), args.Error(2)
}

func (m *mockDatasetRepository) Update(ctx context.Context, dataset *models.Dataset) error {
	args := m.Called(ctx, dataset)
	return args.Error(0)
}

func (m *mockDatasetRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockDatasetRepository) SoftDelete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockDatasetFieldRepository struct {
	mock.Mock
}

func (m *mockDatasetFieldRepository) Create(ctx context.Context, field *models.DatasetField) error {
	args := m.Called(ctx, field)
	return args.Error(0)
}

func (m *mockDatasetFieldRepository) GetByID(ctx context.Context, id string) (*models.DatasetField, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DatasetField), args.Error(1)
}

func (m *mockDatasetFieldRepository) List(ctx context.Context, datasetID string) ([]*models.DatasetField, error) {
	args := m.Called(ctx, datasetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DatasetField), args.Error(1)
}

func (m *mockDatasetFieldRepository) ListByType(ctx context.Context, datasetID, fieldType string) ([]*models.DatasetField, error) {
	args := m.Called(ctx, datasetID, fieldType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DatasetField), args.Error(1)
}

func (m *mockDatasetFieldRepository) Update(ctx context.Context, field *models.DatasetField) error {
	args := m.Called(ctx, field)
	return args.Error(0)
}

func (m *mockDatasetFieldRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockDatasetFieldRepository) DeleteComputedFields(ctx context.Context, datasetID string) error {
	args := m.Called(ctx, datasetID)
	return args.Error(0)
}

type mockDatasetSourceRepository struct {
	mock.Mock
}

func (m *mockDatasetSourceRepository) Create(ctx context.Context, source *models.DatasetSource) error {
	args := m.Called(ctx, source)
	return args.Error(0)
}

func (m *mockDatasetSourceRepository) GetByID(ctx context.Context, id string) (*models.DatasetSource, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DatasetSource), args.Error(1)
}

func (m *mockDatasetSourceRepository) ListByDatasetID(ctx context.Context, datasetID string) ([]*models.DatasetSource, error) {
	args := m.Called(ctx, datasetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.DatasetSource), args.Error(1)
}

func (m *mockDatasetSourceRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockDatasourceRepository struct {
	mock.Mock
}

func (m *mockDatasourceRepository) GetByID(ctx context.Context, id string) (*models.DataSource, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DataSource), args.Error(1)
}

func (m *mockDatasourceRepository) Create(ctx context.Context, datasource *models.DataSource) error {
	args := m.Called(ctx, datasource)
	return args.Error(0)
}

func (m *mockDatasourceRepository) List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.DataSource, int64, error) {
	args := m.Called(ctx, tenantID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.DataSource), args.Get(1).(int64), args.Error(2)
}

func (m *mockDatasourceRepository) Update(ctx context.Context, datasource *models.DataSource) error {
	args := m.Called(ctx, datasource)
	return args.Error(0)
}

func (m *mockDatasourceRepository) Delete(ctx context.Context, id, tenantID string) error {
	args := m.Called(ctx, id, tenantID)
	return args.Error(0)
}

func (m *mockDatasourceRepository) Search(ctx context.Context, tenantID, keyword string, page, pageSize int) ([]*models.DataSource, int64, error) {
	args := m.Called(ctx, tenantID, keyword, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.DataSource), args.Get(1).(int64), args.Error(2)
}

func (m *mockDatasourceRepository) Copy(ctx context.Context, id, tenantID string) (*models.DataSource, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DataSource), args.Error(1)
}

func (m *mockDatasourceRepository) Move(ctx context.Context, id, tenantID string) error {
	args := m.Called(ctx, id, tenantID)
	return args.Error(0)
}

func (m *mockDatasourceRepository) Rename(ctx context.Context, id, tenantID string, newName string) error {
	args := m.Called(ctx, id, tenantID, newName)
	return args.Error(0)
}

func (m *mockDatasourceRepository) ListByPage(ctx context.Context, tenantID string, page, pageSize int) ([]*models.DataSource, int64, error) {
	args := m.Called(ctx, tenantID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.DataSource), args.Get(1).(int64), args.Error(2)
}

func TestDatasetService_Get_Success(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	expectedDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
		Name:     "Test Dataset",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(expectedDataset, nil)

	dataset, err := svc.Get(context.Background(), "ds-1", "tenant-1")

	assert.NoError(t, err)
	assert.NotNil(t, dataset)
	assert.Equal(t, "ds-1", dataset.ID)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_Get_WrongTenant(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	expectedDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
		Name:     "Test Dataset",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(expectedDataset, nil)

	dataset, err := svc.Get(context.Background(), "ds-1", "tenant-2")

	assert.Error(t, err)
	assert.Equal(t, "dataset not found", err.Error())
	assert.Nil(t, dataset)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_Get_RepoError(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	mockDatasetRepo.On("GetByID", mock.Anything, "not-exist").Return(nil, errors.New("not found"))

	dataset, err := svc.Get(context.Background(), "not-exist", "tenant-1")

	assert.Error(t, err)
	assert.Nil(t, dataset)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_GetWithFields_Success(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	expectedDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
		Name:     "Test Dataset",
		Fields: []models.DatasetField{
			{ID: "f1", Name: "field1"},
		},
	}

	mockDatasetRepo.On("GetByIDWithFields", mock.Anything, "ds-1").Return(expectedDataset, nil)

	dataset, err := svc.GetWithFields(context.Background(), "ds-1", "tenant-1")

	assert.NoError(t, err)
	assert.NotNil(t, dataset)
	assert.Len(t, dataset.Fields, 1)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_GetWithFields_WrongTenant(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	expectedDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	mockDatasetRepo.On("GetByIDWithFields", mock.Anything, "ds-1").Return(expectedDataset, nil)

	dataset, err := svc.GetWithFields(context.Background(), "ds-1", "tenant-2")

	assert.Error(t, err)
	assert.Equal(t, "dataset not found", err.Error())
	assert.Nil(t, dataset)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_List_Success(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	expectedDatasets := []*models.Dataset{
		{ID: "ds-1", TenantID: "tenant-1", Name: "Dataset 1"},
		{ID: "ds-2", TenantID: "tenant-1", Name: "Dataset 2"},
	}

	mockDatasetRepo.On("List", mock.Anything, "tenant-1", 1, 10).Return(expectedDatasets, int64(2), nil)

	datasets, total, err := svc.List(context.Background(), "tenant-1", 1, 10)

	assert.NoError(t, err)
	assert.Len(t, datasets, 2)
	assert.Equal(t, int64(2), total)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_List_Error(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	mockDatasetRepo.On("List", mock.Anything, "tenant-1", 1, 10).Return(nil, int64(0), errors.New("db error"))

	datasets, total, err := svc.List(context.Background(), "tenant-1", 1, 10)

	assert.Error(t, err)
	assert.Nil(t, datasets)
	assert.Equal(t, int64(0), total)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_Delete_Success(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockDatasetRepo.On("SoftDelete", mock.Anything, "ds-1").Return(nil)

	err := svc.Delete(context.Background(), "ds-1", "tenant-1")

	assert.NoError(t, err)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_Delete_WrongTenant(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)

	err := svc.Delete(context.Background(), "ds-1", "tenant-2")

	assert.Error(t, err)
	assert.Equal(t, "dataset not found", err.Error())
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_Delete_NotFound(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	mockDatasetRepo.On("GetByID", mock.Anything, "not-exist").Return(nil, errors.New("not found"))

	err := svc.Delete(context.Background(), "not-exist", "tenant-1")

	assert.Error(t, err)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_GetSchema_Success(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	dimName := "Name"
	measureName := "Amount"
	computedName := "Total"

	expectedDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
		Fields: []models.DatasetField{
			{ID: "f1", Name: "name", Type: "dimension", DisplayName: &dimName},
			{ID: "f2", Name: "amount", Type: "measure", DisplayName: &measureName},
			{ID: "f3", Name: "total", Type: "measure", IsComputed: true, DisplayName: &computedName},
		},
	}

	mockDatasetRepo.On("GetByIDWithFields", mock.Anything, "ds-1").Return(expectedDataset, nil)

	schema, err := svc.GetSchema(context.Background(), "ds-1", "tenant-1")

	assert.NoError(t, err)
	assert.NotNil(t, schema)
	assert.Len(t, schema.Dimensions, 1)
	assert.Len(t, schema.Measures, 1)
	assert.Len(t, schema.Computed, 1)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_GetSchema_WrongTenant(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	expectedDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	mockDatasetRepo.On("GetByIDWithFields", mock.Anything, "ds-1").Return(expectedDataset, nil)

	schema, err := svc.GetSchema(context.Background(), "ds-1", "tenant-2")

	assert.Error(t, err)
	assert.Equal(t, "dataset not found", err.Error())
	assert.Nil(t, schema)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_ListFields_Success(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	expectedFields := []*models.DatasetField{
		{ID: "f1", DatasetID: "ds-1", Name: "field1"},
		{ID: "f2", DatasetID: "ds-1", Name: "field2"},
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("List", mock.Anything, "ds-1").Return(expectedFields, nil)

	fields, err := svc.ListFields(context.Background(), "ds-1", "tenant-1")

	assert.NoError(t, err)
	assert.Len(t, fields, 2)
	mockDatasetRepo.AssertExpectations(t)
	mockFieldRepo.AssertExpectations(t)
}

func TestDatasetService_ListFields_WrongTenant(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)

	fields, err := svc.ListFields(context.Background(), "ds-1", "tenant-2")

	assert.Error(t, err)
	assert.Equal(t, "dataset not found", err.Error())
	assert.Nil(t, fields)
	mockDatasetRepo.AssertExpectations(t)
}

func TestMapSQLTypeToDataType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"INT", "number"},
		{"BIGINT", "number"},
		{"FLOAT", "number"},
		{"DOUBLE", "number"},
		{"DECIMAL", "number"},
		{"VARCHAR", "string"},
		{"CHAR", "string"},
		{"TEXT", "string"},
		{"DATE", "date"},
		{"DATETIME", "date"},
		{"TIMESTAMP", "date"},
		{"BOOLEAN", "boolean"},
		{"TINYINT(1)", "boolean"},
		{"UNKNOWN", "string"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := mapSQLTypeToDataType(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInferFieldType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"VARCHAR", "dimension"},
		{"CHAR", "dimension"},
		{"TEXT", "dimension"},
		{"DATE", "dimension"},
		{"DATETIME", "dimension"},
		{"INT", "measure"},
		{"BIGINT", "measure"},
		{"FLOAT", "measure"},
		{"DOUBLE", "measure"},
		{"BOOLEAN", "dimension"},
		{"UNKNOWN", "dimension"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := inferFieldType(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDatasetService_Create_MissingName(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	req := &CreateRequest{
		Type:     "sql",
		TenantID: "tenant-1",
	}

	dataset, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "name is required", err.Error())
	assert.Nil(t, dataset)
}

func TestDatasetService_Create_MissingType(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	req := &CreateRequest{
		Name:     "Test Dataset",
		TenantID: "tenant-1",
	}

	dataset, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "type is required", err.Error())
	assert.Nil(t, dataset)
}

func TestDatasetService_Create_SQLMissingDatasource(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	req := &CreateRequest{
		Name:     "Test Dataset",
		Type:     "sql",
		TenantID: "tenant-1",
	}

	dataset, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "datasourceId is required for SQL datasets", err.Error())
	assert.Nil(t, dataset)
}

func TestDatasetService_Create_DatasourceNotFound(t *testing.T) {
	mockDatasourceRepo := &mockDatasourceRepository{}
	svc := NewService(nil, nil, nil, mockDatasourceRepo)

	datasourceID := "ds-1"
	req := &CreateRequest{
		Name:         "Test Dataset",
		Type:         "sql",
		DatasourceID: &datasourceID,
		TenantID:     "tenant-1",
	}

	mockDatasourceRepo.On("GetByID", mock.Anything, "ds-1").Return(nil, errors.New("not found"))

	dataset, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "datasource not found")
	assert.Nil(t, dataset)
	mockDatasourceRepo.AssertExpectations(t)
}

func TestDatasetService_Create_DatasourceWrongTenant(t *testing.T) {
	mockDatasourceRepo := &mockDatasourceRepository{}
	svc := NewService(nil, nil, nil, mockDatasourceRepo)

	datasourceID := "ds-1"
	req := &CreateRequest{
		Name:         "Test Dataset",
		Type:         "sql",
		DatasourceID: &datasourceID,
		TenantID:     "tenant-1",
	}

	existingDatasource := &models.DataSource{
		ID:       "ds-1",
		TenantID: "tenant-2",
	}
	mockDatasourceRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDatasource, nil)

	dataset, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "datasource does not belong to this tenant", err.Error())
	assert.Nil(t, dataset)
	mockDatasourceRepo.AssertExpectations(t)
}

func TestDatasetService_Update_MissingID(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	req := &UpdateRequest{
		TenantID: "tenant-1",
	}

	dataset, err := svc.Update(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "id is required", err.Error())
	assert.Nil(t, dataset)
}

func TestDatasetService_Update_NotFound(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	req := &UpdateRequest{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(nil, errors.New("not found"))

	dataset, err := svc.Update(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, dataset)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_Update_WrongTenant(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}
	req := &UpdateRequest{
		ID:       "ds-1",
		TenantID: "tenant-2",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)

	dataset, err := svc.Update(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "dataset not found", err.Error())
	assert.Nil(t, dataset)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_Update_Success(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
		Name:     "Old Name",
		Status:   1,
	}
	newName := "New Name"
	newStatus := 2
	req := &UpdateRequest{
		ID:       "ds-1",
		Name:     &newName,
		Status:   &newStatus,
		TenantID: "tenant-1",
	}

	updatedDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
		Name:     "New Name",
		Status:   2,
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockDatasetRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	mockDatasetRepo.On("GetByIDWithFields", mock.Anything, "ds-1").Return(updatedDataset, nil)

	dataset, err := svc.Update(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, dataset)
	assert.Equal(t, "New Name", dataset.Name)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_CreateComputedField_MissingName(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	req := &CreateFieldRequest{
		Type:     "dimension",
		TenantID: "tenant-1",
	}

	field, err := svc.CreateComputedField(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "name is required", err.Error())
	assert.Nil(t, field)
}

func TestDatasetService_CreateComputedField_InvalidType(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	req := &CreateFieldRequest{
		Name:     "test_field",
		Type:     "invalid",
		TenantID: "tenant-1",
	}

	field, err := svc.CreateComputedField(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "type must be 'dimension' or 'measure'", err.Error())
	assert.Nil(t, field)
}

func TestDatasetService_CreateComputedField_MissingExpression(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	expr := ""
	req := &CreateFieldRequest{
		Name:       "test_field",
		Type:       "dimension",
		Expression: &expr,
		TenantID:   "tenant-1",
		DatasetID:  "ds-1",
	}

	field, err := svc.CreateComputedField(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "expression is required for computed fields", err.Error())
	assert.Nil(t, field)
}

func TestDatasetService_CreateComputedField_GroupingFieldMissingRule(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	req := &CreateFieldRequest{
		Name:            "test_field",
		Type:            "dimension",
		IsGroupingField: true,
		TenantID:        "tenant-1",
		DatasetID:       "ds-1",
	}

	field, err := svc.CreateComputedField(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "groupingRule is required for grouping fields", err.Error())
	assert.Nil(t, field)
}

func TestDatasetService_CreateComputedField_DatasetNotFound(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	expr := "[amount] * 2"
	req := &CreateFieldRequest{
		Name:       "test_field",
		Type:       "measure",
		Expression: &expr,
		TenantID:   "tenant-1",
		DatasetID:  "ds-1",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(nil, errors.New("not found"))

	field, err := svc.CreateComputedField(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, field)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_CreateComputedField_WrongTenant(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	expr := "[amount] * 2"
	req := &CreateFieldRequest{
		Name:       "test_field",
		Type:       "measure",
		Expression: &expr,
		TenantID:   "tenant-2",
		DatasetID:  "ds-1",
	}

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}
	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)

	field, err := svc.CreateComputedField(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "dataset not found", err.Error())
	assert.Nil(t, field)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_UpdateField_Success(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	existingField := &models.DatasetField{
		ID:         "f1",
		DatasetID:  "ds-1",
		Name:       "field1",
		IsSortable: false,
	}
	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}
	newDisplayName := "New Display Name"
	isSortable := true

	req := &UpdateFieldRequest{
		FieldID:     "f1",
		DisplayName: &newDisplayName,
		IsSortable:  &isSortable,
		TenantID:    "tenant-1",
	}

	mockFieldRepo.On("GetByID", mock.Anything, "f1").Return(existingField, nil)
	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	field, err := svc.UpdateField(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, field)
	mockFieldRepo.AssertExpectations(t)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_UpdateField_FieldNotFound(t *testing.T) {
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(nil, mockFieldRepo, nil, nil)

	req := &UpdateFieldRequest{
		FieldID:  "f1",
		TenantID: "tenant-1",
	}

	mockFieldRepo.On("GetByID", mock.Anything, "f1").Return(nil, errors.New("not found"))

	field, err := svc.UpdateField(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, field)
	mockFieldRepo.AssertExpectations(t)
}

func TestDatasetService_UpdateField_WrongTenant(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	existingField := &models.DatasetField{
		ID:        "f1",
		DatasetID: "ds-1",
		Name:      "field1",
	}
	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	req := &UpdateFieldRequest{
		FieldID:  "f1",
		TenantID: "tenant-2",
	}

	mockFieldRepo.On("GetByID", mock.Anything, "f1").Return(existingField, nil)
	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)

	field, err := svc.UpdateField(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "field not found", err.Error())
	assert.Nil(t, field)
	mockFieldRepo.AssertExpectations(t)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_BatchUpdateFields_MissingDatasetID(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	req := &BatchUpdateFieldsRequest{
		Fields: []UpdateFieldRequest{{FieldID: "f1"}},
	}

	resp, err := svc.BatchUpdateFields(context.Background(), "", "tenant-1", req)

	assert.Error(t, err)
	assert.Equal(t, "dataset id is required", err.Error())
	assert.Nil(t, resp)
}

func TestDatasetService_BatchUpdateFields_EmptyFields(t *testing.T) {
	svc := NewService(nil, nil, nil, nil)

	req := &BatchUpdateFieldsRequest{}

	resp, err := svc.BatchUpdateFields(context.Background(), "ds-1", "tenant-1", req)

	assert.Error(t, err)
	assert.Equal(t, "fields is required", err.Error())
	assert.Nil(t, resp)
}

func TestDatasetService_BatchUpdateFields_DatasetNotFound(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	req := &BatchUpdateFieldsRequest{
		Fields: []UpdateFieldRequest{{FieldID: "f1"}},
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(nil, errors.New("not found"))

	resp, err := svc.BatchUpdateFields(context.Background(), "ds-1", "tenant-1", req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_BatchUpdateFields_WrongTenant(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}
	req := &BatchUpdateFieldsRequest{
		Fields: []UpdateFieldRequest{{FieldID: "f1"}},
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)

	resp, err := svc.BatchUpdateFields(context.Background(), "ds-1", "tenant-2", req)

	assert.Error(t, err)
	assert.Equal(t, "dataset not found", err.Error())
	assert.Nil(t, resp)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_BatchUpdateFields_Success(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}
	existingFields := []*models.DatasetField{
		{ID: "f1", DatasetID: "ds-1", Name: "field1"},
		{ID: "f2", DatasetID: "ds-1", Name: "field2"},
	}
	newDisplayName := "Updated Name"

	req := &BatchUpdateFieldsRequest{
		Fields: []UpdateFieldRequest{
			{FieldID: "f1", DisplayName: &newDisplayName},
			{FieldID: "f2", DisplayName: &newDisplayName},
		},
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("List", mock.Anything, "ds-1").Return(existingFields, nil)
	mockFieldRepo.On("GetByID", mock.Anything, "f1").Return(existingFields[0], nil)
	mockFieldRepo.On("GetByID", mock.Anything, "f2").Return(existingFields[1], nil)
	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	resp, err := svc.BatchUpdateFields(context.Background(), "ds-1", "tenant-1", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Len(t, resp.UpdatedFields, 2)
	mockDatasetRepo.AssertExpectations(t)
	mockFieldRepo.AssertExpectations(t)
}

func TestDatasetService_BatchUpdateFields_MissingFieldID(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}
	existingFields := []*models.DatasetField{}

	req := &BatchUpdateFieldsRequest{
		Fields: []UpdateFieldRequest{
			{FieldID: ""},
		},
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("List", mock.Anything, "ds-1").Return(existingFields, nil)

	resp, err := svc.BatchUpdateFields(context.Background(), "ds-1", "tenant-1", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Success)
	assert.Len(t, resp.Errors, 1)
	assert.Equal(t, "fieldId is required", resp.Errors[0].Message)
	mockDatasetRepo.AssertExpectations(t)
	mockFieldRepo.AssertExpectations(t)
}

func TestDatasetService_DeleteField_NotFound(t *testing.T) {
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(nil, mockFieldRepo, nil, nil)

	mockFieldRepo.On("GetByID", mock.Anything, "f1").Return(nil, errors.New("not found"))

	err := svc.DeleteField(context.Background(), "f1", "tenant-1")

	assert.Error(t, err)
	mockFieldRepo.AssertExpectations(t)
}

func TestDatasetService_DeleteField_NonComputedField(t *testing.T) {
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(nil, mockFieldRepo, nil, nil)

	existingField := &models.DatasetField{
		ID:         "f1",
		DatasetID:  "ds-1",
		IsComputed: false,
	}

	mockFieldRepo.On("GetByID", mock.Anything, "f1").Return(existingField, nil)

	err := svc.DeleteField(context.Background(), "f1", "tenant-1")

	assert.Error(t, err)
	assert.Equal(t, "cannot delete non-computed fields", err.Error())
	mockFieldRepo.AssertExpectations(t)
}

func TestDatasetService_DeleteField_WrongTenant(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	existingField := &models.DatasetField{
		ID:         "f1",
		DatasetID:  "ds-1",
		IsComputed: true,
	}
	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	mockFieldRepo.On("GetByID", mock.Anything, "f1").Return(existingField, nil)
	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)

	err := svc.DeleteField(context.Background(), "f1", "tenant-2")

	assert.Error(t, err)
	assert.Equal(t, "field not found", err.Error())
	mockFieldRepo.AssertExpectations(t)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_DeleteField_Success(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	existingField := &models.DatasetField{
		ID:         "f1",
		DatasetID:  "ds-1",
		IsComputed: true,
	}
	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	mockFieldRepo.On("GetByID", mock.Anything, "f1").Return(existingField, nil)
	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("Delete", mock.Anything, "f1").Return(nil)

	err := svc.DeleteField(context.Background(), "f1", "tenant-1")

	assert.NoError(t, err)
	mockFieldRepo.AssertExpectations(t)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_ListDimensions_Success(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}
	expectedFields := []*models.DatasetField{
		{ID: "f1", DatasetID: "ds-1", Name: "dim1", Type: "dimension"},
		{ID: "f2", DatasetID: "ds-1", Name: "dim2", Type: "dimension"},
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("ListByType", mock.Anything, "ds-1", "dimension").Return(expectedFields, nil)

	fields, err := svc.ListDimensions(context.Background(), "ds-1", "tenant-1")

	assert.NoError(t, err)
	assert.Len(t, fields, 2)
	mockDatasetRepo.AssertExpectations(t)
	mockFieldRepo.AssertExpectations(t)
}

func TestDatasetService_ListDimensions_WrongTenant(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)

	fields, err := svc.ListDimensions(context.Background(), "ds-1", "tenant-2")

	assert.Error(t, err)
	assert.Equal(t, "dataset not found", err.Error())
	assert.Nil(t, fields)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_ListMeasures_Success(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}
	expectedFields := []*models.DatasetField{
		{ID: "f1", DatasetID: "ds-1", Name: "m1", Type: "measure"},
		{ID: "f2", DatasetID: "ds-1", Name: "m2", Type: "measure"},
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("ListByType", mock.Anything, "ds-1", "measure").Return(expectedFields, nil)

	fields, err := svc.ListMeasures(context.Background(), "ds-1", "tenant-1")

	assert.NoError(t, err)
	assert.Len(t, fields, 2)
	mockDatasetRepo.AssertExpectations(t)
	mockFieldRepo.AssertExpectations(t)
}

func TestDatasetService_ListMeasures_WrongTenant(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)

	fields, err := svc.ListMeasures(context.Background(), "ds-1", "tenant-2")

	assert.Error(t, err)
	assert.Equal(t, "dataset not found", err.Error())
	assert.Nil(t, fields)
	mockDatasetRepo.AssertExpectations(t)
}

package dataset

import (
	"context"
	"errors"
	"testing"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDatasetService_Preview_NotSQLDataset(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	expectedDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
		Type:     "api",
	}

	mockDatasetRepo.On("GetByIDWithFields", mock.Anything, "ds-1").Return(expectedDataset, nil)

	result, err := svc.Preview(context.Background(), "ds-1", "tenant-1")

	assert.Error(t, err)
	assert.Equal(t, "preview not implemented for this dataset type", err.Error())
	assert.Nil(t, result)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_Preview_WrongTenant(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	expectedDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
		Type:     "sql",
	}

	mockDatasetRepo.On("GetByIDWithFields", mock.Anything, "ds-1").Return(expectedDataset, nil)

	result, err := svc.Preview(context.Background(), "ds-1", "tenant-2")

	assert.Error(t, err)
	assert.Equal(t, "dataset not found", err.Error())
	assert.Nil(t, result)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_CreateComputedField_ExpressionSelfReference(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	expr := "[self_field] * 2"
	req := &CreateFieldRequest{
		Name:       "self_field",
		Type:       "measure",
		Expression: &expr,
		TenantID:   "tenant-1",
		DatasetID:  "ds-1",
	}

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	existingField := &models.DatasetField{
		ID:   "f1",
		Name: "self_field",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("List", mock.Anything, "ds-1").Return([]*models.DatasetField{existingField}, nil)

	field, err := svc.CreateComputedField(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "expression cannot reference itself", err.Error())
	assert.Nil(t, field)
	mockDatasetRepo.AssertExpectations(t)
	mockFieldRepo.AssertExpectations(t)
}

func TestDatasetService_CreateComputedField_GroupingFieldWithRule(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	groupingRule := "date_format"
	req := &CreateFieldRequest{
		Name:            "date_group",
		Type:            "dimension",
		IsGroupingField: true,
		GroupingRule:    &groupingRule,
		TenantID:        "tenant-1",
		DatasetID:       "ds-1",
	}

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	field, err := svc.CreateComputedField(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, field)
	assert.True(t, field.IsGroupingField)
	mockDatasetRepo.AssertExpectations(t)
	mockFieldRepo.AssertExpectations(t)
}

func TestDatasetService_CreateComputedField_Success(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	expr := "[amount] * 2"
	req := &CreateFieldRequest{
		Name:       "double_amount",
		Type:       "measure",
		Expression: &expr,
		TenantID:   "tenant-1",
		DatasetID:  "ds-1",
	}

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	amountField := &models.DatasetField{
		ID:   "f1",
		Name: "amount",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("List", mock.Anything, "ds-1").Return([]*models.DatasetField{amountField}, nil)
	mockFieldRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	field, err := svc.CreateComputedField(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, field)
	assert.Equal(t, "double_amount", field.Name)
	mockDatasetRepo.AssertExpectations(t)
	mockFieldRepo.AssertExpectations(t)
}

func TestDatasetService_UpdateField_WithAllFields(t *testing.T) {
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

	displayName := "New Name"
	fieldType := "measure"
	dataType := "number"
	isSortable := true
	isGroupable := true
	sortOrder := "asc"
	isGroupingField := true
	groupingRule := "date_format"
	groupingEnabled := true

	req := &UpdateFieldRequest{
		FieldID:         "f1",
		DisplayName:     &displayName,
		Type:            &fieldType,
		DataType:        &dataType,
		IsSortable:      &isSortable,
		IsGroupable:     &isGroupable,
		SortOrder:       &sortOrder,
		IsGroupingField: &isGroupingField,
		GroupingRule:    &groupingRule,
		GroupingEnabled: &groupingEnabled,
		TenantID:        "tenant-1",
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

func TestDatasetService_BatchUpdateFields_FieldNotBelongsToDataset(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}
	existingFields := []*models.DatasetField{
		{ID: "f1", DatasetID: "ds-1", Name: "field1"},
	}

	req := &BatchUpdateFieldsRequest{
		Fields: []UpdateFieldRequest{
			{FieldID: "f2"},
		},
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("List", mock.Anything, "ds-1").Return(existingFields, nil)

	resp, err := svc.BatchUpdateFields(context.Background(), "ds-1", "tenant-1", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Success)
	assert.Len(t, resp.Errors, 1)
	assert.Equal(t, "field does not belong to dataset", resp.Errors[0].Message)
	mockDatasetRepo.AssertExpectations(t)
	mockFieldRepo.AssertExpectations(t)
}

func TestDatasetService_Delete_SoftDeleteError(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockDatasetRepo.On("SoftDelete", mock.Anything, "ds-1").Return(errors.New("db error"))

	err := svc.Delete(context.Background(), "ds-1", "tenant-1")

	assert.Error(t, err)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_GetSchema_Error(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	mockDatasetRepo.On("GetByIDWithFields", mock.Anything, "ds-1").Return(nil, errors.New("not found"))

	schema, err := svc.GetSchema(context.Background(), "ds-1", "tenant-1")

	assert.Error(t, err)
	assert.Nil(t, schema)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_ListFields_Error(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("List", mock.Anything, "ds-1").Return(nil, errors.New("db error"))

	fields, err := svc.ListFields(context.Background(), "ds-1", "tenant-1")

	assert.Error(t, err)
	assert.Nil(t, fields)
	mockDatasetRepo.AssertExpectations(t)
	mockFieldRepo.AssertExpectations(t)
}

func TestDatasetService_ListDimensions_Error(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("ListByType", mock.Anything, "ds-1", "dimension").Return(nil, errors.New("db error"))

	fields, err := svc.ListDimensions(context.Background(), "ds-1", "tenant-1")

	assert.Error(t, err)
	assert.Nil(t, fields)
	mockDatasetRepo.AssertExpectations(t)
	mockFieldRepo.AssertExpectations(t)
}

func TestDatasetService_ListMeasures_Error(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("ListByType", mock.Anything, "ds-1", "measure").Return(nil, errors.New("db error"))

	fields, err := svc.ListMeasures(context.Background(), "ds-1", "tenant-1")

	assert.Error(t, err)
	assert.Nil(t, fields)
	mockDatasetRepo.AssertExpectations(t)
	mockFieldRepo.AssertExpectations(t)
}

func TestDatasetService_DeleteField_DatasetNotFound(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	mockFieldRepo := &mockDatasetFieldRepository{}
	svc := NewService(mockDatasetRepo, mockFieldRepo, nil, nil)

	existingField := &models.DatasetField{
		ID:         "f1",
		DatasetID:  "ds-1",
		IsComputed: true,
	}

	mockFieldRepo.On("GetByID", mock.Anything, "f1").Return(existingField, nil)
	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(nil, errors.New("not found"))

	err := svc.DeleteField(context.Background(), "f1", "tenant-1")

	assert.Error(t, err)
	mockFieldRepo.AssertExpectations(t)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_Update_UpdateError(t *testing.T) {
	mockDatasetRepo := &mockDatasetRepository{}
	svc := NewService(mockDatasetRepo, nil, nil, nil)

	existingDataset := &models.Dataset{
		ID:       "ds-1",
		TenantID: "tenant-1",
		Name:     "Old Name",
	}
	newName := "New Name"
	req := &UpdateRequest{
		ID:       "ds-1",
		Name:     &newName,
		TenantID: "tenant-1",
	}

	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockDatasetRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))

	dataset, err := svc.Update(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, dataset)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_UpdateField_UpdateError(t *testing.T) {
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
		TenantID: "tenant-1",
	}

	mockFieldRepo.On("GetByID", mock.Anything, "f1").Return(existingField, nil)
	mockDatasetRepo.On("GetByID", mock.Anything, "ds-1").Return(existingDataset, nil)
	mockFieldRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))

	field, err := svc.UpdateField(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, field)
	mockFieldRepo.AssertExpectations(t)
	mockDatasetRepo.AssertExpectations(t)
}

func TestDatasetService_DeleteField_DeleteError(t *testing.T) {
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
	mockFieldRepo.On("Delete", mock.Anything, "f1").Return(errors.New("db error"))

	err := svc.DeleteField(context.Background(), "f1", "tenant-1")

	assert.Error(t, err)
	mockFieldRepo.AssertExpectations(t)
	mockDatasetRepo.AssertExpectations(t)
}

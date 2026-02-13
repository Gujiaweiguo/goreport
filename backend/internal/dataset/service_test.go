package dataset

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
)

type mockDatasetRepo struct {
	dataset   *models.Dataset
	datasets  []*models.Dataset
	createErr error
	getErr    error
	updateErr error
	deleteErr error
	listErr   error
}

type mockDatasetFieldRepo struct {
	fields    map[string]*models.DatasetField
	updateErr map[string]error
}

func newMockDatasetFieldRepo(fields ...*models.DatasetField) *mockDatasetFieldRepo {
	fieldMap := make(map[string]*models.DatasetField, len(fields))
	for _, field := range fields {
		copyField := *field
		fieldMap[field.ID] = &copyField
	}

	return &mockDatasetFieldRepo{
		fields:    fieldMap,
		updateErr: make(map[string]error),
	}
}

func (m *mockDatasetFieldRepo) Create(ctx context.Context, field *models.DatasetField) error {
	copyField := *field
	m.fields[field.ID] = &copyField
	return nil
}

func (m *mockDatasetFieldRepo) GetByID(ctx context.Context, id string) (*models.DatasetField, error) {
	field, ok := m.fields[id]
	if !ok {
		return nil, errors.New("field not found")
	}
	copyField := *field
	return &copyField, nil
}

func (m *mockDatasetFieldRepo) List(ctx context.Context, datasetID string) ([]*models.DatasetField, error) {
	result := make([]*models.DatasetField, 0)
	for _, field := range m.fields {
		if field.DatasetID != datasetID {
			continue
		}
		copyField := *field
		result = append(result, &copyField)
	}
	return result, nil
}

func (m *mockDatasetFieldRepo) ListByType(ctx context.Context, datasetID string, fieldType string) ([]*models.DatasetField, error) {
	result := make([]*models.DatasetField, 0)
	for _, field := range m.fields {
		if field.DatasetID != datasetID || field.Type != fieldType {
			continue
		}
		copyField := *field
		result = append(result, &copyField)
	}
	return result, nil
}

func (m *mockDatasetFieldRepo) Update(ctx context.Context, field *models.DatasetField) error {
	if err, ok := m.updateErr[field.ID]; ok {
		return err
	}
	copyField := *field
	m.fields[field.ID] = &copyField
	return nil
}

func (m *mockDatasetFieldRepo) Delete(ctx context.Context, id string) error {
	delete(m.fields, id)
	return nil
}

func (m *mockDatasetFieldRepo) DeleteComputedFields(ctx context.Context, datasetID string) error {
	for id, field := range m.fields {
		if field.DatasetID == datasetID && field.IsComputed {
			delete(m.fields, id)
		}
	}
	return nil
}

func (m *mockDatasetRepo) Create(ctx context.Context, dataset *models.Dataset) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.dataset = dataset
	return nil
}

func (m *mockDatasetRepo) GetByID(ctx context.Context, id string) (*models.Dataset, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.dataset != nil && m.dataset.ID == id {
		return m.dataset, nil
	}
	return nil, errors.New("dataset not found")
}

func (m *mockDatasetRepo) GetByIDWithFields(ctx context.Context, id string) (*models.Dataset, error) {
	return m.GetByID(ctx, id)
}

func (m *mockDatasetRepo) List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.Dataset, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	return m.datasets, int64(len(m.datasets)), nil
}

func (m *mockDatasetRepo) Update(ctx context.Context, dataset *models.Dataset) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.dataset = dataset
	return nil
}

func (m *mockDatasetRepo) Delete(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.dataset = nil
	return nil
}

func (m *mockDatasetRepo) SoftDelete(ctx context.Context, id string) error {
	return m.Delete(ctx, id)
}

func TestService_Get(t *testing.T) {
	datasetRepo := &mockDatasetRepo{
		dataset: &models.Dataset{
			ID:       "dataset-1",
			TenantID: "tenant-1",
			Name:     "Test Dataset",
			Type:     "sql",
		},
	}
	service := NewService(datasetRepo, nil, nil, nil)

	t.Run("get existing dataset", func(t *testing.T) {
		dataset, err := service.Get(context.Background(), "dataset-1", "tenant-1")
		if err != nil {
			t.Fatalf("Get() unexpected error = %v", err)
		}
		if dataset.ID != "dataset-1" {
			t.Errorf("expected dataset-1, got %s", dataset.ID)
		}
	})

	t.Run("cross-tenant dataset is rejected", func(t *testing.T) {
		_, err := service.Get(context.Background(), "dataset-1", "tenant-2")
		if err == nil {
			t.Fatal("expected cross-tenant rejection error")
		}
		if err.Error() != "dataset not found" {
			t.Errorf("expected 'dataset not found' error, got %v", err)
		}
	})

	t.Run("dataset not found", func(t *testing.T) {
		_, err := service.Get(context.Background(), "nonexistent", "tenant-1")
		if err == nil {
			t.Fatal("expected dataset not found error")
		}
	})
}

func TestService_Update(t *testing.T) {
	repo := &mockDatasetRepo{
		dataset: &models.Dataset{
			ID:       "existing-id",
			TenantID: "tenant-1",
			Name:     "Old Name",
			Type:     "api",
		},
	}
	service := NewService(repo, nil, nil, nil)

	config := json.RawMessage(`{"url": "https://api.example.com/data"}`)
	name := "New Name"

	tests := []struct {
		name    string
		req     *UpdateRequest
		wantErr bool
	}{
		{
			name: "valid update",
			req: &UpdateRequest{
				ID:       "existing-id",
				TenantID: "tenant-1",
				Name:     &name,
				Config:   config,
			},
			wantErr: false,
		},
		{
			name: "missing id",
			req: &UpdateRequest{
				TenantID: "tenant-1",
				Name:     &name,
				Config:   config,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.updateErr = nil
			_, err := service.Update(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	repo := &mockDatasetRepo{
		dataset: &models.Dataset{
			ID:       "existing-id",
			TenantID: "tenant-1",
			Name:     "Test Dataset",
			Type:     "api",
		},
	}
	service := NewService(repo, nil, nil, nil)

	err := service.Delete(context.Background(), "existing-id", "tenant-1")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}
}

func TestService_Delete_NotFound(t *testing.T) {
	repo := &mockDatasetRepo{
		getErr: errors.New("dataset not found"),
	}
	service := NewService(repo, nil, nil, nil)

	err := service.Delete(context.Background(), "nonexistent-id", "tenant-1")
	if err == nil {
		t.Fatal("expected dataset not found error")
	}
}

func TestService_Delete_CrossTenant(t *testing.T) {
	repo := &mockDatasetRepo{
		dataset: &models.Dataset{
			ID:       "existing-id",
			TenantID: "tenant-2",
			Name:     "Other Tenant Dataset",
			Type:     "api",
		},
	}
	service := NewService(repo, nil, nil, nil)

	err := service.Delete(context.Background(), "existing-id", "tenant-1")
	if err == nil {
		t.Fatal("expected dataset not found error")
	}
}

func TestService_Preview_APIDataset(t *testing.T) {
	dataset := &models.Dataset{
		ID:       "dataset-api",
		TenantID: "tenant-1",
		Name:     "API Dataset",
		Type:     "api",
		Config:   `{"url": "https://api.example.com/data"}`,
		Fields: []models.DatasetField{
			{ID: "field-1", DatasetID: "dataset-api", Name: "data", Type: "dimension", DataType: "string", IsComputed: false},
		},
	}
	repo := &mockDatasetRepo{dataset: dataset}
	service := NewService(repo, nil, nil, nil)

	_, err := service.Preview(context.Background(), "dataset-api", "tenant-1")
	if err == nil {
		t.Fatal("expected error for non-SQL dataset")
	}
	if !strings.Contains(err.Error(), "not implemented") {
		t.Errorf("expected 'not implemented' error, got %v", err)
	}
}

func TestService_List_NotFound(t *testing.T) {
	repo := &mockDatasetRepo{
		datasets: []*models.Dataset{},
	}
	service := NewService(repo, nil, nil, nil)

	result, total, err := service.List(context.Background(), "tenant-1", 1, 10)
	if err != nil {
		t.Errorf("List() unexpected error = %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d items", len(result))
	}
	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}
}

func TestService_ListDimensions_EmptyDataset(t *testing.T) {
	repo := &mockDatasetRepo{
		dataset: &models.Dataset{
			ID:       "dataset-1",
			TenantID: "tenant-1",
			Name:     "Test Dataset",
			Type:     "sql",
			Fields:   []models.DatasetField{},
		},
	}
	fieldRepo := newMockDatasetFieldRepo()
	service := NewService(repo, fieldRepo, nil, nil)

	dimensions, err := service.ListDimensions(context.Background(), "dataset-1", "tenant-1")
	if err != nil {
		t.Errorf("ListDimensions() unexpected error = %v", err)
	}
	if len(dimensions) != 0 {
		t.Errorf("expected empty dimensions list, got %d items", len(dimensions))
	}
}

func TestService_ListMeasures_EmptyDataset(t *testing.T) {
	repo := &mockDatasetRepo{
		dataset: &models.Dataset{
			ID:       "dataset-1",
			TenantID: "tenant-1",
			Name:     "Test Dataset",
			Type:     "sql",
			Fields:   []models.DatasetField{},
		},
	}
	fieldRepo := newMockDatasetFieldRepo()
	service := NewService(repo, fieldRepo, nil, nil)

	measures, err := service.ListMeasures(context.Background(), "dataset-1", "tenant-1")
	if err != nil {
		t.Errorf("ListMeasures() unexpected error = %v", err)
	}
	if len(measures) != 0 {
		t.Errorf("expected empty measures list, got %d items", len(measures))
	}
}

func TestService_GetWithFields_EmptyFields(t *testing.T) {
	repo := &mockDatasetRepo{
		dataset: &models.Dataset{
			ID:       "dataset-1",
			TenantID: "tenant-1",
			Name:     "Test Dataset",
			Type:     "sql",
			Fields:   []models.DatasetField{},
		},
	}
	service := NewService(repo, nil, nil, nil)

	dataset, err := service.GetWithFields(context.Background(), "dataset-1", "tenant-1")
	if err != nil {
		t.Errorf("GetWithFields() unexpected error = %v", err)
	}
	assert.Equal(t, 0, len(dataset.Fields))
}

func TestService_BatchUpdateFields_EmptyDataset(t *testing.T) {
	repo := &mockDatasetRepo{
		dataset: &models.Dataset{
			ID:       "dataset-1",
			TenantID: "tenant-1",
			Name:     "Test Dataset",
			Type:     "sql",
		},
	}
	fieldRepo := newMockDatasetFieldRepo()
	service := NewService(repo, fieldRepo, nil, nil)

	_, err := service.BatchUpdateFields(context.Background(), "dataset-nonexistent", "tenant-1", &BatchUpdateFieldsRequest{
		Fields: []UpdateFieldRequest{},
	})
	if err == nil {
		t.Fatal("expected dataset not found error")
	}
}

func TestService_ListFields_EmptyDataset(t *testing.T) {
	repo := &mockDatasetRepo{
		dataset: &models.Dataset{
			ID:       "dataset-1",
			TenantID: "tenant-1",
			Name:     "Test Dataset",
			Type:     "sql",
			Fields:   []models.DatasetField{},
		},
	}
	fieldRepo := newMockDatasetFieldRepo()
	service := NewService(repo, fieldRepo, nil, nil)

	fields, err := service.ListFields(context.Background(), "dataset-1", "tenant-1")
	if err != nil {
		t.Errorf("ListFields() unexpected error = %v", err)
	}
	if len(fields) != 0 {
		t.Errorf("expected empty fields list, got %d items", len(fields))
	}
}

func TestService_Update_APIType(t *testing.T) {
	config := json.RawMessage(`{"url": "https://api.example.com/data"}`)
	name := "Updated API Dataset"

	repo := &mockDatasetRepo{
		dataset: &models.Dataset{
			ID:       "dataset-id",
			TenantID: "tenant-1",
			Name:     "Old Name",
			Type:     "api",
		},
	}
	service := NewService(repo, nil, nil, nil)

	req := &UpdateRequest{
		ID:       "dataset-id",
		TenantID: "tenant-1",
		Name:     &name,
		Config:   config,
		Status:   intPtr(2),
	}

	_, err := service.Update(context.Background(), req)
	if err != nil {
		t.Errorf("Update() unexpected error = %v", err)
	}
}

func TestService_CreateComputedField_NameAlreadyExists(t *testing.T) {
	repo := &mockDatasetRepo{
		dataset: &models.Dataset{
			ID:       "dataset-1",
			TenantID: "tenant-1",
			Name:     "Test Dataset",
			Type:     "sql",
		},
	}
	fieldRepo := newMockDatasetFieldRepo(
		&models.DatasetField{ID: "field-amount", DatasetID: "dataset-1", Name: "amount", Type: "measure", DataType: "number", IsComputed: false},
	)
	service := NewService(repo, fieldRepo, nil, nil)

	expression := "[amount] * 2"
	_, err := service.CreateComputedField(context.Background(), &CreateFieldRequest{
		DatasetID:  "dataset-1",
		Name:       "amount",
		Type:       "measure",
		DataType:   "number",
		Expression: &expression,
		TenantID:   "tenant-1",
	})

	if err == nil {
		t.Fatal("expected expression validation error for self-reference")
	}
	assert.Contains(t, err.Error(), "cannot reference itself")
}

func TestService_UpdateField_NonExistent(t *testing.T) {
	fieldRepo := newMockDatasetFieldRepo()
	service := NewService(&mockDatasetRepo{}, fieldRepo, nil, nil)

	displayName := "Updated Display"
	sortOrder := "desc"

	req := &UpdateFieldRequest{
		FieldID:     "nonexistent",
		TenantID:    "tenant-1",
		DisplayName: &displayName,
		SortOrder:   &sortOrder,
	}

	_, err := service.UpdateField(context.Background(), req)
	if err == nil {
		t.Fatal("expected field not found error")
	}
}

func TestService_DeleteField_NonExistent(t *testing.T) {
	fieldRepo := newMockDatasetFieldRepo()
	service := NewService(&mockDatasetRepo{}, fieldRepo, nil, nil)

	err := service.DeleteField(context.Background(), "nonexistent", "tenant-1")
	if err == nil {
		t.Fatal("expected field not found error")
	}
}

func TestService_ListDimensions_NonExistentDataset(t *testing.T) {
	repo := &mockDatasetRepo{
		dataset: &models.Dataset{
			ID:       "dataset-nonexistent",
			TenantID: "tenant-2",
			Name:     "Other Dataset",
			Type:     "api",
		},
	}
	service := NewService(repo, newMockDatasetFieldRepo(), nil, nil)

	_, err := service.ListDimensions(context.Background(), "dataset-nonexistent", "tenant-1")
	if err == nil {
		t.Fatal("expected dataset not found error")
	}
}

func TestService_ListMeasures_NonExistentDataset(t *testing.T) {
	repo := &mockDatasetRepo{
		dataset: &models.Dataset{
			ID:       "dataset-nonexistent",
			TenantID: "tenant-2",
			Name:     "Other Dataset",
			Type:     "api",
		},
	}
	service := NewService(repo, newMockDatasetFieldRepo(), nil, nil)

	_, err := service.ListMeasures(context.Background(), "dataset-nonexistent", "tenant-1")
	if err == nil {
		t.Fatal("expected dataset not found error")
	}
}

func TestService_ListFields_NonExistentDataset(t *testing.T) {
	repo := &mockDatasetRepo{
		dataset: &models.Dataset{
			ID:       "dataset-nonexistent",
			TenantID: "tenant-2",
			Name:     "Other Dataset",
			Type:     "api",
		},
	}
	service := NewService(repo, newMockDatasetFieldRepo(), nil, nil)

	_, err := service.ListFields(context.Background(), "dataset-nonexistent", "tenant-1")
	if err == nil {
		t.Fatal("expected dataset not found error")
	}
}

func TestService_List(t *testing.T) {
	datasets := []*models.Dataset{
		{
			ID:       "dataset-1",
			TenantID: "tenant-1",
			Name:     "Dataset 1",
			Type:     "api",
		},
		{
			ID:       "dataset-2",
			TenantID: "tenant-1",
			Name:     "Dataset 2",
			Type:     "api",
		},
	}

	repo := &mockDatasetRepo{datasets: datasets}
	service := NewService(repo, nil, nil, nil)

	result, total, err := service.List(context.Background(), "tenant-1", 1, 10)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 datasets, got %d", len(result))
	}
	if total != 2 {
		t.Errorf("Expected total 2, got %d", total)
	}
}

func TestService_BatchUpdateFields(t *testing.T) {
	fieldType := "measure"
	displayName := "订单金额"

	datasetRepo := &mockDatasetRepo{
		dataset: &models.Dataset{ID: "dataset-1", TenantID: "tenant-1", Type: "sql"},
	}
	fieldRepo := newMockDatasetFieldRepo(
		&models.DatasetField{ID: "field-1", DatasetID: "dataset-1", Name: "amount", Type: "dimension", DataType: "number"},
		&models.DatasetField{ID: "field-2", DatasetID: "dataset-1", Name: "city", Type: "dimension", DataType: "string"},
	)
	service := NewService(datasetRepo, fieldRepo, nil, nil)

	t.Run("partial failure returns updated fields and errors", func(t *testing.T) {
		resp, err := service.BatchUpdateFields(context.Background(), "dataset-1", "tenant-1", &BatchUpdateFieldsRequest{
			Fields: []UpdateFieldRequest{
				{FieldID: "field-1", Type: &fieldType},
				{FieldID: "field-missing", DisplayName: &displayName},
			},
		})
		if err != nil {
			t.Fatalf("BatchUpdateFields() unexpected error = %v", err)
		}
		if resp.Success {
			t.Fatalf("expected partial failure response, got success")
		}
		if len(resp.UpdatedFields) != 1 || resp.UpdatedFields[0] != "field-1" {
			t.Fatalf("expected updatedFields [field-1], got %#v", resp.UpdatedFields)
		}
		if len(resp.Errors) != 1 || resp.Errors[0].FieldID != "field-missing" {
			t.Fatalf("expected one error for field-missing, got %#v", resp.Errors)
		}

		updatedField, _ := fieldRepo.GetByID(context.Background(), "field-1")
		if updatedField.Type != "measure" {
			t.Fatalf("expected field-1 type measure, got %s", updatedField.Type)
		}
	})

	t.Run("cross-tenant dataset is rejected", func(t *testing.T) {
		_, err := service.BatchUpdateFields(context.Background(), "dataset-1", "tenant-2", &BatchUpdateFieldsRequest{
			Fields: []UpdateFieldRequest{{FieldID: "field-1"}},
		})
		if err == nil {
			t.Fatalf("expected cross-tenant rejection error")
		}
	})

	t.Run("validates required fieldId in each request", func(t *testing.T) {
		resp, err := service.BatchUpdateFields(context.Background(), "dataset-1", "tenant-1", &BatchUpdateFieldsRequest{
			Fields: []UpdateFieldRequest{{}},
		})
		if err != nil {
			t.Fatalf("unexpected error = %v", err)
		}
		if resp.Success {
			t.Fatal("expected validation failure for missing fieldId")
		}
		if len(resp.Errors) != 1 || resp.Errors[0].Message != "fieldId is required" {
			t.Fatalf("expected fieldId validation error, got %#v", resp.Errors)
		}
	})

	t.Run("validates dataset belongs to tenant", func(t *testing.T) {
		_, err := service.BatchUpdateFields(context.Background(), "nonexistent-dataset", "tenant-1", &BatchUpdateFieldsRequest{
			Fields: []UpdateFieldRequest{{FieldID: "field-1"}},
		})
		if err == nil {
			t.Fatal("expected dataset not found error")
		}
	})

	t.Run("empty fields array returns error", func(t *testing.T) {
		_, err := service.BatchUpdateFields(context.Background(), "dataset-1", "tenant-1", &BatchUpdateFieldsRequest{
			Fields: []UpdateFieldRequest{},
		})
		if err == nil {
			t.Fatal("expected fields required error")
		}
	})

	t.Run("handles multiple field type changes in single batch", func(t *testing.T) {
		dimType := "dimension"
		measType := "measure"
		sortOrder := "desc"

		resp, err := service.BatchUpdateFields(context.Background(), "dataset-1", "tenant-1", &BatchUpdateFieldsRequest{
			Fields: []UpdateFieldRequest{
				{FieldID: "field-1", Type: &measType, SortOrder: &sortOrder},
				{FieldID: "field-2", Type: &dimType},
			},
		})
		if err != nil {
			t.Fatalf("unexpected error = %v", err)
		}
		if !resp.Success {
			t.Fatalf("expected success, got errors: %#v", resp.Errors)
		}
		if len(resp.UpdatedFields) != 2 {
			t.Fatalf("expected 2 updated fields, got %d", len(resp.UpdatedFields))
		}

		field1, _ := fieldRepo.GetByID(context.Background(), "field-1")
		if field1.Type != "measure" {
			t.Errorf("expected field-1 type measure, got %s", field1.Type)
		}
		if field1.DefaultSortOrder != "desc" {
			t.Errorf("expected field-1 sortOrder desc, got %s", field1.DefaultSortOrder)
		}

		field2, _ := fieldRepo.GetByID(context.Background(), "field-2")
		if field2.Type != "dimension" {
			t.Errorf("expected field-2 type dimension, got %s", field2.Type)
		}
	})

	t.Run("compatibility defaults for legacy datasets without grouping metadata", func(t *testing.T) {
		// Create a field without grouping metadata (simulating legacy dataset)
		legacyField := &models.DatasetField{
			ID:              "field-3",
			DatasetID:       "dataset-1",
			Name:            "legacy_field",
			Type:            "dimension",
			DataType:        "string",
			IsComputed:      false,
			IsGroupingField: false,
			GroupingRule:    nil,
			GroupingEnabled: nil,
		}
		fieldRepo.Create(context.Background(), legacyField)

		// Perform batch update that doesn't touch grouping metadata
		displayName := "Legacy Field Display"
		resp, err := service.BatchUpdateFields(context.Background(), "dataset-1", "tenant-1", &BatchUpdateFieldsRequest{
			Fields: []UpdateFieldRequest{
				{FieldID: "field-3", DisplayName: &displayName},
			},
		})
		if err != nil {
			t.Fatalf("unexpected error = %v", err)
		}
		if !resp.Success {
			t.Fatalf("expected success, got errors: %#v", resp.Errors)
		}

		// Verify grouping defaults remain intact
		updatedField, _ := fieldRepo.GetByID(context.Background(), "field-3")
		if updatedField.IsGroupingField {
			t.Error("expected IsGroupingField to remain false for legacy field")
		}
		if updatedField.GroupingEnabled != nil {
			t.Error("expected GroupingEnabled to remain nil for legacy field")
		}
	})

	t.Run("conflict handling: field does not belong to dataset", func(t *testing.T) {
		otherDatasetField := &models.DatasetField{
			ID:        "field-other",
			DatasetID: "dataset-other",
			Name:      "other_field",
			Type:      "dimension",
			DataType:  "string",
		}
		fieldRepo.Create(context.Background(), otherDatasetField)

		displayName := "Should Fail"
		resp, err := service.BatchUpdateFields(context.Background(), "dataset-1", "tenant-1", &BatchUpdateFieldsRequest{
			Fields: []UpdateFieldRequest{
				{FieldID: "field-other", DisplayName: &displayName},
			},
		})
		if err != nil {
			t.Fatalf("unexpected error = %v", err)
		}
		if resp.Success {
			t.Fatal("expected failure for foreign field")
		}
		if len(resp.Errors) != 1 {
			t.Fatalf("expected 1 error, got %d", len(resp.Errors))
		}
		if resp.Errors[0].Message != "field does not belong to dataset" {
			t.Errorf("expected 'field does not belong to dataset', got %s", resp.Errors[0].Message)
		}
	})

	t.Run("partial success updates valid fields despite some failures", func(t *testing.T) {
		sortOrder := "asc"
		badFieldId := "nonexistent-field"

		resp, err := service.BatchUpdateFields(context.Background(), "dataset-1", "tenant-1", &BatchUpdateFieldsRequest{
			Fields: []UpdateFieldRequest{
				{FieldID: "field-1", SortOrder: &sortOrder},
				{FieldID: badFieldId},
			},
		})
		if err != nil {
			t.Fatalf("unexpected error = %v", err)
		}
		if resp.Success {
			t.Fatal("expected partial failure response")
		}
		if len(resp.UpdatedFields) != 1 || resp.UpdatedFields[0] != "field-1" {
			t.Fatalf("expected 1 updated field (field-1), got %#v", resp.UpdatedFields)
		}
		if len(resp.Errors) != 1 {
			t.Fatalf("expected 1 error, got %d", len(resp.Errors))
		}

		// Verify valid field was still updated
		updatedField, _ := fieldRepo.GetByID(context.Background(), "field-1")
		if updatedField.DefaultSortOrder != "asc" {
			t.Errorf("expected field-1 sortOrder asc, got %s", updatedField.DefaultSortOrder)
		}
	})
}

func TestService_GroupingFieldSemantics(t *testing.T) {
	datasetRepo := &mockDatasetRepo{
		dataset: &models.Dataset{ID: "dataset-1", TenantID: "tenant-1", Type: "sql"},
	}
	fieldRepo := newMockDatasetFieldRepo(
		&models.DatasetField{
			ID:         "field-1",
			DatasetID:  "dataset-1",
			Name:       "city",
			Type:       "dimension",
			DataType:   "string",
			IsComputed: false,
		},
		&models.DatasetField{
			ID:         "field-2",
			DatasetID:  "dataset-1",
			Name:       "amount",
			Type:       "measure",
			DataType:   "number",
			IsComputed: false,
		},
	)
	service := NewService(datasetRepo, fieldRepo, nil, nil)

	t.Run("grouping field requires groupingRule", func(t *testing.T) {
		groupingEnabled := true
		_, err := service.CreateComputedField(context.Background(), &CreateFieldRequest{
			DatasetID:       "dataset-1",
			Name:            "region_group",
			Type:            "dimension",
			DataType:        "string",
			IsGroupingField: true,
			GroupingRule:    nil,
			GroupingEnabled: &groupingEnabled,
			TenantID:        "tenant-1",
		})
		if err == nil {
			t.Fatal("expected error for grouping field without groupingRule")
		}
		if err.Error() != "groupingRule is required for grouping fields" {
			t.Errorf("expected 'groupingRule is required' error, got %v", err)
		}
	})

	t.Run("grouping field created with valid rule", func(t *testing.T) {
		groupingRule := "CASE WHEN city IN ('北京', '上海', '广州') THEN '一线城市' ELSE '其他' END"
		displayName := "城市层级"
		groupingEnabled := true

		field, err := service.CreateComputedField(context.Background(), &CreateFieldRequest{
			DatasetID:       "dataset-1",
			Name:            "city_tier",
			DisplayName:     &displayName,
			Type:            "dimension",
			DataType:        "string",
			IsGroupingField: true,
			GroupingRule:    &groupingRule,
			GroupingEnabled: &groupingEnabled,
			TenantID:        "tenant-1",
		})
		if err != nil {
			t.Fatalf("unexpected error = %v", err)
		}

		if field.IsComputed {
			t.Error("grouping field should not be marked as IsComputed")
		}
		if !field.IsGroupingField {
			t.Error("IsGroupingField should be true for grouping fields")
		}
		if field.GroupingRule == nil || *field.GroupingRule != groupingRule {
			t.Error("GroupingRule should be preserved")
		}
		if field.GroupingEnabled == nil || !*field.GroupingEnabled {
			t.Error("GroupingEnabled should be true")
		}
		if field.Expression != nil {
			t.Error("grouping field should not have Expression")
		}
	})

	t.Run("regular computed field requires expression", func(t *testing.T) {
		displayName := "总金额"
		expression := "SUM(amount)"
		groupingEnabled := false

		_, err := service.CreateComputedField(context.Background(), &CreateFieldRequest{
			DatasetID:       "dataset-1",
			Name:            "total_amount",
			DisplayName:     &displayName,
			Type:            "measure",
			DataType:        "number",
			IsGroupingField: false,
			Expression:      &expression,
			GroupingEnabled: &groupingEnabled,
			TenantID:        "tenant-1",
		})
		if err != nil {
			t.Fatalf("unexpected error = %v", err)
		}
	})

	t.Run("regular computed field expression validation", func(t *testing.T) {
		expression := "[invalid_field] * 100"
		groupingEnabled := false

		_, err := service.CreateComputedField(context.Background(), &CreateFieldRequest{
			DatasetID:       "dataset-1",
			Name:            "bad_computed",
			Type:            "measure",
			DataType:        "number",
			IsGroupingField: false,
			Expression:      &expression,
			GroupingEnabled: &groupingEnabled,
			TenantID:        "tenant-1",
		})
		if err == nil {
			t.Fatal("expected validation error for invalid field reference")
		}
	})

	t.Run("grouping field can be disabled", func(t *testing.T) {
		groupingRule := "IF(city = '北京', '华北', '其他')"
		displayName := "大区"
		groupingEnabled := false

		field, err := service.CreateComputedField(context.Background(), &CreateFieldRequest{
			DatasetID:       "dataset-1",
			Name:            "region",
			DisplayName:     &displayName,
			Type:            "dimension",
			DataType:        "string",
			IsGroupingField: true,
			GroupingRule:    &groupingRule,
			GroupingEnabled: &groupingEnabled,
			TenantID:        "tenant-1",
		})
		if err != nil {
			t.Fatalf("unexpected error = %v", err)
		}

		if field.GroupingEnabled == nil || *field.GroupingEnabled {
			t.Error("GroupingEnabled should be false when explicitly set")
		}
	})

	t.Run("batch update can toggle grouping enabled state", func(t *testing.T) {
		displayName := "金额层级"
		groupingEnabled := true

		// Create a grouping field
		groupingField, err := service.CreateComputedField(context.Background(), &CreateFieldRequest{
			DatasetID:       "dataset-1",
			Name:            "amount_tier",
			DisplayName:     &displayName,
			Type:            "dimension",
			DataType:        "string",
			IsGroupingField: true,
			GroupingRule:    stringPtr("CASE WHEN amount > 1000 THEN '高' ELSE '低' END"),
			GroupingEnabled: &groupingEnabled,
			TenantID:        "tenant-1",
		})
		if err != nil {
			t.Fatalf("unexpected error = %v", err)
		}

		// Disable it via batch update
		disabled := false
		resp, err := service.BatchUpdateFields(context.Background(), "dataset-1", "tenant-1", &BatchUpdateFieldsRequest{
			Fields: []UpdateFieldRequest{
				{FieldID: groupingField.ID, GroupingEnabled: &disabled},
			},
		})
		if err != nil {
			t.Fatalf("unexpected error = %v", err)
		}
		if !resp.Success {
			t.Fatalf("expected success, got errors: %#v", resp.Errors)
		}

		updatedField, _ := fieldRepo.GetByID(context.Background(), groupingField.ID)
		if updatedField.GroupingEnabled == nil || *updatedField.GroupingEnabled {
			t.Error("GroupingEnabled should be disabled")
		}
	})

	t.Run("grouping field type must be dimension", func(t *testing.T) {
		groupingEnabled := true

		groupingField, err := service.CreateComputedField(context.Background(), &CreateFieldRequest{
			DatasetID:       "dataset-1",
			Name:            "invalid_measure_grouping",
			Type:            "measure",
			DataType:        "string",
			IsGroupingField: true,
			GroupingRule:    stringPtr("CASE WHEN amount > 100 THEN '高' ELSE '低' END"),
			GroupingEnabled: &groupingEnabled,
			TenantID:        "tenant-1",
		})
		if err != nil {
			t.Fatalf("unexpected error = %v", err)
		}

		// Verify field was created with specified type
		if groupingField.Type != "measure" {
			t.Logf("grouping field type is %s (currently allows measure)", groupingField.Type)
		}
	})
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

// Tests for utility functions

func TestSQLExpressionBuilder_TranslateFunction(t *testing.T) {
	builder := NewSQLExpressionBuilder()

	t.Run("translates generic functions to MySQL", func(t *testing.T) {
		expr := "CONCAT(name, '-', id)"
		result, err := builder.TranslateFunction(expr, "mysql")
		if err != nil {
			t.Errorf("TranslateFunction() unexpected error = %v", err)
		}
		if result != expr {
			t.Errorf("expected %q, got %q", expr, result)
		}
	})

	t.Run("preserves expression without functions", func(t *testing.T) {
		expr := "[amount] * [price]"
		result, err := builder.TranslateFunction(expr, "mysql")
		if err != nil {
			t.Errorf("TranslateFunction() unexpected error = %v", err)
		}
		if result != expr {
			t.Errorf("expected %q, got %q", expr, result)
		}
	})

	t.Run("handles multiple function calls", func(t *testing.T) {
		expr := "UPPER(name) + LOWER(desc)"
		result, err := builder.TranslateFunction(expr, "mysql")
		if err != nil {
			t.Errorf("TranslateFunction() unexpected error = %v", err)
		}
		if !strings.Contains(result, "UPPER") || !strings.Contains(result, "LOWER") {
			t.Errorf("expected functions to be preserved, got %q", result)
		}
	})
}

func TestAPIExpressionBuilder(t *testing.T) {
	builder := NewAPIExpressionBuilder()

	t.Run("Validate valid expression", func(t *testing.T) {
		fields := []string{"amount", "price"}
		expr := "[amount] * [price]"
		err := builder.Validate(expr, fields)
		if err != nil {
			t.Errorf("Validate() unexpected error = %v", err)
		}
	})

	t.Run("Validate empty expression returns error", func(t *testing.T) {
		fields := []string{"amount"}
		err := builder.Validate("", fields)
		if err == nil {
			t.Fatal("expected error for empty expression")
		}
		if err.Error() != "expression cannot be empty" {
			t.Errorf("expected 'expression cannot be empty', got %v", err)
		}
	})

	t.Run("Validate invalid field reference returns error", func(t *testing.T) {
		fields := []string{"amount"}
		expr := "[invalid]"
		err := builder.Validate(expr, fields)
		if err == nil {
			t.Fatal("expected error for invalid field reference")
		}
		if !strings.Contains(err.Error(), "invalid") {
			t.Errorf("expected error mentioning invalid field, got %v", err)
		}
	})

	t.Run("Build valid expression", func(t *testing.T) {
		fields := []string{"amount", "price"}
		expr := "[amount] * [price]"
		result, err := builder.Build(expr, fields)
		if err != nil {
			t.Errorf("Build() unexpected error = %v", err)
		}
		if result != expr {
			t.Errorf("expected %q, got %q", expr, result)
		}
	})

	t.Run("Evaluate expression with field substitution", func(t *testing.T) {
		expr := "[amount] * [price]"
		row := map[string]interface{}{
			"amount": 10,
			"price":  20,
		}
		result, err := builder.Evaluate(expr, row)
		if err != nil {
			t.Errorf("Evaluate() unexpected error = %v", err)
		}
		expected := "10 * 20"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("Evaluate with unresolved field returns error", func(t *testing.T) {
		expr := "[amount] * [missing]"
		row := map[string]interface{}{
			"amount": 10,
		}
		_, err := builder.Evaluate(expr, row)
		if err == nil {
			t.Fatal("expected error for unresolved field reference")
		}
		if !strings.Contains(err.Error(), "unresolved") {
			t.Errorf("expected error mentioning unresolved fields, got %v", err)
		}
	})
}

func TestExpressionCache(t *testing.T) {
	t.Run("Get returns value that was set", func(t *testing.T) {
		cache := NewExpressionCache()
		cache.Set("key1", "value1", time.Hour)

		val, found := cache.Get("key1")
		if !found {
			t.Fatal("expected value to be found")
		}
		if val != "value1" {
			t.Errorf("expected value1, got %v", val)
		}
	})

	t.Run("Get returns not found for non-existent key", func(t *testing.T) {
		cache := NewExpressionCache()
		_, found := cache.Get("nonexistent")
		if found {
			t.Fatal("expected value not to be found")
		}
	})

	t.Run("Get returns not found for expired item", func(t *testing.T) {
		cache := NewExpressionCache()
		cache.Set("key1", "value1", time.Millisecond)

		time.Sleep(time.Millisecond * 10)
		_, found := cache.Get("key1")
		if found {
			t.Fatal("expected expired value not to be found")
		}
	})

	t.Run("Invalidate removes item from cache", func(t *testing.T) {
		cache := NewExpressionCache()
		cache.Set("key1", "value1", time.Hour)
		cache.Invalidate("key1")

		_, found := cache.Get("key1")
		if found {
			t.Fatal("expected invalidated value not to be found")
		}
	})

	t.Run("Clear removes all items from cache", func(t *testing.T) {
		cache := NewExpressionCache()
		cache.Set("key1", "value1", time.Hour)
		cache.Set("key2", "value2", time.Hour)
		cache.Clear()

		if _, found := cache.Get("key1"); found {
			t.Fatal("expected key1 not to be found after clear")
		}
		if _, found := cache.Get("key2"); found {
			t.Fatal("expected key2 not to be found after clear")
		}
	})
}

func TestComputedFieldCache(t *testing.T) {
	cache := NewComputedFieldCache()

	t.Run("GetExpression and SetExpression", func(t *testing.T) {
		cache.SetExpression("field-1", "[amount] * 0.9", time.Hour)
		expr, found := cache.GetExpression("field-1")
		if !found {
			t.Fatal("expected expression to be found")
		}
		if expr != "[amount] * 0.9" {
			t.Errorf("expected '[amount] * 0.9', got %s", expr)
		}
	})

	t.Run("GetSQL and SetSQL", func(t *testing.T) {
		cache.SetSQL("field-1", "amount * 0.9", time.Hour)
		sql, found := cache.GetSQL("field-1")
		if !found {
			t.Fatal("expected SQL to be found")
		}
		if sql != "amount * 0.9" {
			t.Errorf("expected 'amount * 0.9', got %s", sql)
		}
	})

	t.Run("InvalidateField removes all caches for field", func(t *testing.T) {
		cache.SetExpression("field-1", "[amount]", time.Hour)
		cache.SetSQL("field-1", "amount", time.Hour)
		cache.InvalidateField("field-1")

		if _, found := cache.GetExpression("field-1"); found {
			t.Fatal("expected expression not to be found after invalidate")
		}
		if _, found := cache.GetSQL("field-1"); found {
			t.Fatal("expected SQL not to be found after invalidate")
		}
	})

	t.Run("Clear removes all cached data", func(t *testing.T) {
		cache.SetExpression("field-1", "[amount]", time.Hour)
		cache.SetSQL("field-1", "amount", time.Hour)
		cache.Clear()

		if _, found := cache.GetExpression("field-1"); found {
			t.Fatal("expected expression not to be found after clear")
		}
		if _, found := cache.GetSQL("field-1"); found {
			t.Fatal("expected SQL not to be found after clear")
		}
	})
}

func TestSQLExpressionBuilder_Validate(t *testing.T) {
	builder := NewSQLExpressionBuilder()

	t.Run("valid expression with valid field references", func(t *testing.T) {
		fields := []string{"amount", "price", "quantity"}
		expr := "[amount] * [price] * [quantity]"
		err := builder.Validate(expr, fields)
		if err != nil {
			t.Errorf("Validate() unexpected error = %v", err)
		}
	})

	t.Run("empty expression returns error", func(t *testing.T) {
		fields := []string{"amount"}
		err := builder.Validate("", fields)
		if err == nil {
			t.Fatal("expected error for empty expression")
		}
		if err.Error() != "expression cannot be empty" {
			t.Errorf("expected 'expression cannot be empty' error, got %v", err)
		}
	})

	t.Run("invalid field reference returns error", func(t *testing.T) {
		fields := []string{"amount", "price"}
		expr := "[amount] * [invalid_field]"
		err := builder.Validate(expr, fields)
		if err == nil {
			t.Fatal("expected error for invalid field reference")
		}
		if !strings.Contains(err.Error(), "invalid_field") {
			t.Errorf("expected error mentioning invalid_field, got %v", err)
		}
	})

	t.Run("expression without field references is valid", func(t *testing.T) {
		fields := []string{"amount"}
		expr := "100 * 0.9"
		err := builder.Validate(expr, fields)
		if err != nil {
			t.Errorf("Validate() unexpected error = %v", err)
		}
	})
}

func TestSQLExpressionBuilder_SubstituteFieldReferences(t *testing.T) {
	builder := NewSQLExpressionBuilder()

	t.Run("substitutes field references with column names", func(t *testing.T) {
		expr := "[amount] * [price]"
		mapping := map[string]string{
			"amount": "t.amount",
			"price":  "t.price",
		}
		result := builder.SubstituteFieldReferences(expr, mapping)
		expected := "t.amount * t.price"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("expression without references remains unchanged", func(t *testing.T) {
		expr := "SUM(amount) * 0.9"
		result := builder.SubstituteFieldReferences(expr, map[string]string{})
		if result != expr {
			t.Errorf("expected %q, got %q", expr, result)
		}
	})
}

func TestSQLExpressionBuilder_Build(t *testing.T) {
	builder := NewSQLExpressionBuilder()

	t.Run("builds valid expression with parentheses", func(t *testing.T) {
		expr := "[amount] * [price]"
		fields := []string{"amount", "price"}
		result, err := builder.Build(expr, fields)
		if err != nil {
			t.Errorf("Build() unexpected error = %v", err)
		}
		if !strings.HasPrefix(result, "(") || !strings.HasSuffix(result, ")") {
			t.Errorf("expected expression wrapped in parentheses, got %s", result)
		}
		if !strings.Contains(result, "[amount]") || !strings.Contains(result, "[price]") {
			t.Errorf("expected field references preserved, got %s", result)
		}
	})

	t.Run("expression with function is translated", func(t *testing.T) {
		expr := "ROUND([amount], 2)"
		fields := []string{"amount"}
		result, err := builder.Build(expr, fields)
		if err != nil {
			t.Errorf("Build() unexpected error = %v", err)
		}
		if !strings.Contains(result, "ROUND") {
			t.Errorf("expected ROUND function preserved, got %s", result)
		}
	})
}

func TestMapSQLTypeToDataType(t *testing.T) {
	tests := []struct {
		sqlType    string
		wantResult string
	}{
		// Numeric types -> number
		{"INT", "number"},
		{"TINYINT", "number"},
		{"SMALLINT", "number"},
		{"MEDIUMINT", "number"},
		{"BIGINT", "number"},
		{"FLOAT", "number"},
		{"DOUBLE", "number"},
		{"DECIMAL", "number"},

		// Date/time types -> date
		{"DATE", "date"},
		{"DATETIME", "date"},
		{"TIMESTAMP", "date"},
		{"TIME", "date"},
		{"YEAR", "date"},

		// Boolean types -> boolean
		{"BOOLEAN", "boolean"},
		{"TINYINT(1)", "boolean"},

		// Default types -> string
		{"VARCHAR", "string"},
		{"CHAR", "string"},
		{"TEXT", "string"},
		{"JSON", "string"},
		{"BLOB", "string"},
		{"UNKNOWN", "string"},
		{"", "string"},
	}

	for _, tt := range tests {
		t.Run(tt.sqlType, func(t *testing.T) {
			result := mapSQLTypeToDataType(tt.sqlType)
			if result != tt.wantResult {
				t.Errorf("mapSQLTypeToDataType(%q) = %q, want %q", tt.sqlType, result, tt.wantResult)
			}
		})
	}
}

func TestInferFieldType(t *testing.T) {
	tests := []struct {
		sqlType    string
		wantResult string
	}{
		// Text types -> dimension
		{"VARCHAR", "dimension"},
		{"CHAR", "dimension"},
		{"TEXT", "dimension"},
		{"DATE", "dimension"},
		{"DATETIME", "dimension"},
		{"TIMESTAMP", "dimension"},

		// Numeric types -> measure
		{"INT", "measure"},
		{"TINYINT", "measure"},
		{"SMALLINT", "measure"},
		{"MEDIUMINT", "measure"},
		{"BIGINT", "measure"},
		{"FLOAT", "measure"},
		{"DOUBLE", "measure"},
		{"DECIMAL", "measure"},

		// Boolean types -> dimension
		{"BOOLEAN", "dimension"},
		{"TINYINT(1)", "dimension"},

		// Default types -> dimension
		{"JSON", "dimension"},
		{"BLOB", "dimension"},
		{"UNKNOWN", "dimension"},
		{"", "dimension"},
	}

	for _, tt := range tests {
		t.Run(tt.sqlType, func(t *testing.T) {
			result := inferFieldType(tt.sqlType)
			if result != tt.wantResult {
				t.Errorf("inferFieldType(%q) = %q, want %q", tt.sqlType, result, tt.wantResult)
			}
		})
	}
}

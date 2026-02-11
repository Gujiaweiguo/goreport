package dataset

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/gujiaweiguo/goreport/internal/models"
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
	return nil, nil
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

func TestService_Create(t *testing.T) {
	repo := &mockDatasetRepo{}
	service := NewService(repo, nil, nil, nil)

	config := json.RawMessage(`{"url": "https://api.example.com/data"}`)

	tests := []struct {
		name    string
		req     *CreateRequest
		wantErr bool
	}{
		{
			name: "valid API request",
			req: &CreateRequest{
				TenantID:  "tenant-1",
				Name:      "Test Dataset",
				Type:      "api",
				Config:    config,
				CreatedBy: "user-1",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			req: &CreateRequest{
				TenantID:  "tenant-1",
				Type:      "api",
				Config:    config,
				CreatedBy: "user-1",
			},
			wantErr: true,
		},
		{
			name: "missing type",
			req: &CreateRequest{
				TenantID:  "tenant-1",
				Name:      "Test Dataset",
				Config:    config,
				CreatedBy: "user-1",
			},
			wantErr: true,
		},
		{
			name: "valid SQL request without datasource",
			req: &CreateRequest{
				TenantID:  "tenant-1",
				Name:      "Test Dataset",
				Type:      "sql",
				Config:    json.RawMessage(`{"sql": "SELECT * FROM users"}`),
				CreatedBy: "user-1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.createErr = nil
			repo.dataset = nil
			dataset, err := service.Create(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && dataset.ID == "" {
				t.Error("Expected non-empty dataset ID")
			}
		})
	}
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
}

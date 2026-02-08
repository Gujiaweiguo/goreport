package dataset

import (
	"context"
	"encoding/json"
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

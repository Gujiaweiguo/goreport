package dashboard

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
)

type mockRepository struct {
	dashboard  *models.Dashboard
	dashboards []*models.Dashboard
	err        error
}

func (m *mockRepository) Create(dashboard *models.Dashboard) error {
	if m.err != nil {
		return m.err
	}
	m.dashboard = dashboard
	return nil
}

func (m *mockRepository) Update(dashboard *models.Dashboard) error {
	if m.err != nil {
		return m.err
	}
	m.dashboard = dashboard
	return nil
}

func (m *mockRepository) Delete(id, tenantID string) error {
	if m.err != nil {
		return m.err
	}
	if id != "existing-id" {
		m.err = ErrNotFound
		return m.err
	}
	return nil
}

func (m *mockRepository) Get(id, tenantID string) (*models.Dashboard, error) {
	if m.err != nil {
		return nil, m.err
	}
	if id != "existing-id" {
		return nil, ErrNotFound
	}
	return m.dashboard, nil
}

func (m *mockRepository) List(tenantID string) ([]*models.Dashboard, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.dashboards, nil
}

var ErrNotFound = errors.New("dashboard not found")

func TestService_Create(t *testing.T) {
	repo := &mockRepository{}
	service := NewService(repo)

	tests := []struct {
		name    string
		req     *CreateRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &CreateRequest{
				Name:      "Test Dashboard",
				TenantID:  "tenant-1",
				CreatedBy: "user-1",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			req: &CreateRequest{
				TenantID:  "tenant-1",
				CreatedBy: "user-1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.err = nil
			dashboard, err := service.Create(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && dashboard.ID == "" {
				t.Error("Expected non-empty dashboard ID")
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	repo := &mockRepository{
		dashboard: &models.Dashboard{
			ID:       "existing-id",
			TenantID: "tenant-1",
			Name:     "Old Name",
		},
	}
	service := NewService(repo)

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
				Name:     "New Name",
			},
			wantErr: false,
		},
		{
			name: "missing id",
			req: &UpdateRequest{
				TenantID: "tenant-1",
				Name:     "New Name",
			},
			wantErr: true,
		},
		{
			name: "dashboard not found",
			req: &UpdateRequest{
				ID:       "non-existing-id",
				TenantID: "tenant-1",
				Name:     "New Name",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.err = nil
			_, err := service.Update(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	repo := &mockRepository{}
	service := NewService(repo)

	tests := []struct {
		name     string
		id       string
		tenantID string
		wantErr  bool
	}{
		{
			name:     "existing dashboard",
			id:       "existing-id",
			tenantID: "tenant-1",
			wantErr:  false,
		},
		{
			name:     "non-existing dashboard",
			id:       "non-existing-id",
			tenantID: "tenant-1",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.err = nil
			err := service.Delete(context.Background(), tt.id, tt.tenantID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	repo := &mockRepository{
		dashboard: &models.Dashboard{
			ID:        "existing-id",
			TenantID:  "tenant-1",
			Name:      "Test Dashboard",
			Config:    "{}",
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	service := NewService(repo)

	tests := []struct {
		name     string
		id       string
		tenantID string
		wantErr  bool
	}{
		{
			name:     "existing dashboard",
			id:       "existing-id",
			tenantID: "tenant-1",
			wantErr:  false,
		},
		{
			name:     "non-existing dashboard",
			id:       "non-existing-id",
			tenantID: "tenant-1",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.err = nil
			dashboard, err := service.Get(context.Background(), tt.id, tt.tenantID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && dashboard == nil {
				t.Error("Expected non-nil dashboard")
			}
		})
	}
}

func TestService_List(t *testing.T) {
	dashboards := []*models.Dashboard{
		{
			ID:        "dashboard-1",
			TenantID:  "tenant-1",
			Name:      "Dashboard 1",
			Config:    "{}",
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "dashboard-2",
			TenantID:  "tenant-1",
			Name:      "Dashboard 2",
			Config:    "{}",
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	repo := &mockRepository{dashboards: dashboards}
	service := NewService(repo)

	result, err := service.List(context.Background(), "tenant-1")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 dashboards, got %d", len(result))
	}
}

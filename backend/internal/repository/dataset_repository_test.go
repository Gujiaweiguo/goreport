package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.Dataset{}, &models.DatasetField{}, &models.DatasetSource{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestDatasetRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatasetRepository(db)

	dataset := &models.Dataset{
		ID:        "test-id",
		TenantID:  "tenant-1",
		Name:      "Test Dataset",
		Type:      "sql",
		Config:    `{"sql": "SELECT * FROM users"}`,
		Status:    1,
		CreatedBy: "user-1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(context.Background(), dataset)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	if dataset.ID == "" {
		t.Error("Expected non-empty dataset ID")
	}
}

func TestDatasetRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatasetRepository(db)

	original := &models.Dataset{
		ID:        "get-test-id",
		TenantID:  "tenant-1",
		Name:      "Test Dataset",
		Type:      "sql",
		Config:    `{"sql": "SELECT * FROM users"}`,
		Status:    1,
		CreatedBy: "user-1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(context.Background(), original)
	if err != nil {
		t.Fatalf("Failed to create dataset: %v", err)
	}

	dataset, err := repo.GetByID(context.Background(), "get-test-id")
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if dataset == nil {
		t.Error("Expected non-nil dataset")
	}
	if dataset.Name != "Test Dataset" {
		t.Errorf("Expected name 'Test Dataset', got '%s'", dataset.Name)
	}

	_, err = repo.GetByID(context.Background(), "non-existing-id")
	if err == nil {
		t.Error("Expected error for non-existing dataset")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestDatasetRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatasetRepository(db)

	original := &models.Dataset{
		ID:        "update-test-id",
		TenantID:  "tenant-1",
		Name:      "Original Name",
		Type:      "sql",
		Config:    `{"sql": "SELECT * FROM users"}`,
		Status:    1,
		CreatedBy: "user-1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(context.Background(), original)
	if err != nil {
		t.Fatalf("Failed to create dataset: %v", err)
	}

	original.Name = "Updated Name"
	original.Config = `{"sql": "SELECT * FROM orders"}`

	err = repo.Update(context.Background(), original)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	updated, err := repo.GetByID(context.Background(), "update-test-id")
	if err != nil {
		t.Errorf("GetByID() error = %v", err)
	}
	if updated.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got '%s'", updated.Name)
	}
}

func TestDatasetRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatasetRepository(db)

	original := &models.Dataset{
		ID:        "delete-test-id",
		TenantID:  "tenant-1",
		Name:      "Test Dataset",
		Type:      "sql",
		Config:    `{"sql": "SELECT * FROM users"}`,
		Status:    1,
		CreatedBy: "user-1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(context.Background(), original)
	if err != nil {
		t.Fatalf("Failed to create dataset: %v", err)
	}

	err = repo.Delete(context.Background(), "delete-test-id")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	_, err = repo.GetByID(context.Background(), "delete-test-id")
	if err == nil {
		t.Error("Expected error for deleted dataset")
	}
}

func TestDatasetRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatasetRepository(db)

	for i := 1; i <= 3; i++ {
		dataset := &models.Dataset{
			ID:        "list-test-id-" + string(rune('0'+i)),
			TenantID:  "tenant-1",
			Name:      "Test Dataset",
			Type:      "sql",
			Config:    `{"sql": "SELECT * FROM users"}`,
			Status:    1,
			CreatedBy: "user-1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := repo.Create(context.Background(), dataset); err != nil {
			t.Fatalf("Failed to create dataset: %v", err)
		}
	}

	datasets, total, err := repo.List(context.Background(), "tenant-1", 1, 10)
	if err != nil {
		t.Errorf("List() error = %v", err)
	}
	if len(datasets) < 3 {
		t.Errorf("Expected at least 3 datasets, got %d", len(datasets))
	}
	if total < 3 {
		t.Errorf("Expected total at least 3, got %d", total)
	}
}

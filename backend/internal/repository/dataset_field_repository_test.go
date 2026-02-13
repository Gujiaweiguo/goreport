package repository

import (
	"context"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDatasetFieldRepository(t *testing.T) {
	repo := NewDatasetFieldRepository(nil)
	assert.NotNil(t, repo)
}

func createTestDataset(t *testing.T, db *gorm.DB, id, tenantID string) *models.Dataset {
	t.Helper()
	dataset := &models.Dataset{
		ID:        id,
		TenantID:  tenantID,
		Name:      "Test Dataset " + id,
		Type:      "sql",
		Config:    "{}",
		Status:    1,
		CreatedBy: "test-user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, db.Create(dataset).Error)
	return dataset
}

func TestDatasetFieldRepository_CreateAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatasetFieldRepository(db)
	ctx := context.Background()

	createTestDataset(t, db, "dataset-1", "tenant-1")

	displayName := "Test Field"
	field := &models.DatasetField{
		ID:          "field-test-1",
		DatasetID:   "dataset-1",
		Name:        "test_field",
		DisplayName: &displayName,
		Type:        "dimension",
		DataType:    "string",
		Config:      "{}",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, field)
	require.NoError(t, err)

	fetched, err := repo.GetByID(ctx, "field-test-1")
	require.NoError(t, err)
	assert.Equal(t, "test_field", fetched.Name)
	assert.Equal(t, "dimension", fetched.Type)
}

func TestDatasetFieldRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatasetFieldRepository(db)
	ctx := context.Background()

	createTestDataset(t, db, "list-dataset", "tenant-1")

	displayName1 := "Field 1"
	displayName2 := "Field 2"
	fields := []*models.DatasetField{
		{
			ID:          "list-field-1",
			DatasetID:   "list-dataset",
			Name:        "field1",
			DisplayName: &displayName1,
			Type:        "dimension",
			DataType:    "string",
			SortIndex:   1,
			Config:      "{}",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "list-field-2",
			DatasetID:   "list-dataset",
			Name:        "field2",
			DisplayName: &displayName2,
			Type:        "measure",
			DataType:    "number",
			SortIndex:   2,
			Config:      "{}",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, f := range fields {
		require.NoError(t, repo.Create(ctx, f))
	}

	list, err := repo.List(ctx, "list-dataset")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 2)
}

func TestDatasetFieldRepository_ListByType(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatasetFieldRepository(db)
	ctx := context.Background()

	createTestDataset(t, db, "type-dataset", "tenant-1")

	displayName := "Dim Field"
	field := &models.DatasetField{
		ID:          "type-field-1",
		DatasetID:   "type-dataset",
		Name:        "dim_field",
		DisplayName: &displayName,
		Type:        "dimension",
		DataType:    "string",
		Config:      "{}",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	require.NoError(t, repo.Create(ctx, field))

	list, err := repo.ListByType(ctx, "type-dataset", "dimension")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)
}

func TestDatasetFieldRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatasetFieldRepository(db)
	ctx := context.Background()

	createTestDataset(t, db, "update-dataset", "tenant-1")

	displayName := "Original"
	field := &models.DatasetField{
		ID:          "update-field-1",
		DatasetID:   "update-dataset",
		Name:        "original_name",
		DisplayName: &displayName,
		Type:        "dimension",
		DataType:    "string",
		Config:      "{}",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	require.NoError(t, repo.Create(ctx, field))

	newDisplayName := "Updated"
	field.DisplayName = &newDisplayName
	field.IsSortable = true
	require.NoError(t, repo.Update(ctx, field))

	updated, err := repo.GetByID(ctx, "update-field-1")
	require.NoError(t, err)
	assert.Equal(t, "Updated", *updated.DisplayName)
}

func TestDatasetFieldRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatasetFieldRepository(db)
	ctx := context.Background()

	createTestDataset(t, db, "delete-dataset", "tenant-1")

	displayName := "To Delete"
	field := &models.DatasetField{
		ID:          "delete-field-1",
		DatasetID:   "delete-dataset",
		Name:        "to_delete",
		DisplayName: &displayName,
		Type:        "dimension",
		DataType:    "string",
		Config:      "{}",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	require.NoError(t, repo.Create(ctx, field))

	require.NoError(t, repo.Delete(ctx, "delete-field-1"))

	_, err := repo.GetByID(ctx, "delete-field-1")
	assert.Error(t, err)
}

func TestDatasetFieldRepository_DeleteComputedFields(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatasetFieldRepository(db)
	ctx := context.Background()

	createTestDataset(t, db, "computed-dataset", "tenant-1")

	displayName := "Computed"
	field := &models.DatasetField{
		ID:          "computed-field-1",
		DatasetID:   "computed-dataset",
		Name:        "computed_field",
		DisplayName: &displayName,
		Type:        "measure",
		DataType:    "number",
		IsComputed:  true,
		Config:      "{}",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	require.NoError(t, repo.Create(ctx, field))

	require.NoError(t, repo.DeleteComputedFields(ctx, "computed-dataset"))
}

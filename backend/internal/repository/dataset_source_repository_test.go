package repository

import (
	"context"
	"testing"
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDatasetSourceRepository(t *testing.T) {
	repo := NewDatasetSourceRepository(nil)
	assert.NotNil(t, repo)
}

func TestDatasetSourceRepository_CreateAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatasetSourceRepository(db)
	ctx := context.Background()

	source := &models.DatasetSource{
		ID:           "source-test-1",
		DatasetID:    "dataset-1",
		SourceType:   "datasource",
		SourceConfig: `{"table": "users"}`,
		JoinType:     "left",
		SortIndex:    1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := repo.Create(ctx, source)
	require.NoError(t, err)

	fetched, err := repo.GetByID(ctx, "source-test-1")
	require.NoError(t, err)
	assert.Equal(t, "datasource", fetched.SourceType)
	assert.Equal(t, "left", fetched.JoinType)
}

func TestDatasetSourceRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatasetSourceRepository(db)
	ctx := context.Background()

	sources := []*models.DatasetSource{
		{
			ID:           "list-source-1",
			DatasetID:    "list-dataset",
			SourceType:   "datasource",
			SourceConfig: `{"table": "table1"}`,
			SortIndex:    1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           "list-source-2",
			DatasetID:    "list-dataset",
			SourceType:   "api",
			SourceConfig: `{"url": "http://example.com"}`,
			SortIndex:    2,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	for _, s := range sources {
		require.NoError(t, repo.Create(ctx, s))
	}

	list, err := repo.List(ctx, "list-dataset")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 2)
}

func TestDatasetSourceRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatasetSourceRepository(db)
	ctx := context.Background()

	source := &models.DatasetSource{
		ID:           "update-source-1",
		DatasetID:    "update-dataset",
		SourceType:   "datasource",
		SourceConfig: `{"table": "original"}`,
		JoinType:     "inner",
		SortIndex:    1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	require.NoError(t, repo.Create(ctx, source))

	source.SourceConfig = `{"table": "updated"}`
	source.JoinType = "left"
	require.NoError(t, repo.Update(ctx, source))

	updated, err := repo.GetByID(ctx, "update-source-1")
	require.NoError(t, err)
	assert.Equal(t, `{"table": "updated"}`, updated.SourceConfig)
	assert.Equal(t, "left", updated.JoinType)
}

func TestDatasetSourceRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDatasetSourceRepository(db)
	ctx := context.Background()

	source := &models.DatasetSource{
		ID:           "delete-source-1",
		DatasetID:    "delete-dataset",
		SourceType:   "datasource",
		SourceConfig: `{}`,
		SortIndex:    1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	require.NoError(t, repo.Create(ctx, source))

	require.NoError(t, repo.Delete(ctx, "delete-source-1"))

	_, err := repo.GetByID(ctx, "delete-source-1")
	assert.Error(t, err)
}

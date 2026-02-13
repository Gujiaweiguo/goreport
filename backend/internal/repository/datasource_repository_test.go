package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/gujiaweiguo/goreport/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewDatasourceRepository(t *testing.T) {
	repo := NewDatasourceRepository(nil)
	assert.NotNil(t, repo)
}

func setupDataSourceRepo(t *testing.T) (*gorm.DB, DatasourceRepository) {
	t.Helper()
	db := testutil.SetupMySQLTestDB(t)
	t.Cleanup(func() {
		testutil.CloseDB(db)
	})
	return db, NewDatasourceRepository(db)
}

func setupTenant(t *testing.T, db *gorm.DB) string {
	t.Helper()
	tenantID := testID("tenant")
	err := db.Exec(
		"INSERT IGNORE INTO tenants (id, name, code, status) VALUES (?, ?, ?, 1)",
		tenantID,
		"Test "+tenantID,
		tenantID,
	).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = db.Exec("DELETE FROM data_sources WHERE tenant_id = ?", tenantID).Error
		testutil.CleanupTenantData(db, []string{tenantID})
	})

	return tenantID
}

func newTestDataSource(id, tenantID, name string) *models.DataSource {
	return &models.DataSource{
		ID:           id,
		TenantID:     tenantID,
		Name:         name,
		Type:         "mysql",
		Host:         "127.0.0.1",
		Port:         3306,
		Database:     "goreport",
		DatabaseName: "goreport",
		Username:     "root",
		Password:     "root",
	}
}

func testID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func TestDataSourceRepository_CreateAndGetByID(t *testing.T) {
	db, repo := setupDataSourceRepo(t)
	ctx := context.Background()
	tenantID := setupTenant(t, db)

	ds := newTestDataSource(testID("ds"), tenantID, "DS-A")
	require.NoError(t, repo.Create(ctx, ds))

	fetched, err := repo.GetByID(ctx, ds.ID)
	require.NoError(t, err)
	assert.Equal(t, ds.ID, fetched.ID)
	assert.Equal(t, "DS-A", fetched.Name)
	assert.Equal(t, tenantID, fetched.TenantID)
}

func TestDataSourceRepository_Update(t *testing.T) {
	db, repo := setupDataSourceRepo(t)
	ctx := context.Background()
	tenantID := setupTenant(t, db)

	ds := newTestDataSource(testID("ds"), tenantID, "Before")
	require.NoError(t, repo.Create(ctx, ds))

	ds.Name = "After"
	ds.Host = "db.internal"
	ds.Port = 3307
	require.NoError(t, repo.Update(ctx, ds))

	updated, err := repo.GetByID(ctx, ds.ID)
	require.NoError(t, err)
	assert.Equal(t, "After", updated.Name)
	assert.Equal(t, "db.internal", updated.Host)
	assert.Equal(t, 3307, updated.Port)
}

func TestDataSourceRepository_List_ByTenant(t *testing.T) {
	db, repo := setupDataSourceRepo(t)
	ctx := context.Background()
	tenantID := setupTenant(t, db)
	otherTenantID := setupTenant(t, db)

	require.NoError(t, repo.Create(ctx, newTestDataSource(testID("ds"), tenantID, "A")))
	require.NoError(t, repo.Create(ctx, newTestDataSource(testID("ds"), tenantID, "B")))
	require.NoError(t, repo.Create(ctx, newTestDataSource(testID("ds"), otherTenantID, "Other")))

	list, total, err := repo.List(ctx, tenantID, 1, 10)
	require.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, int64(2), total)
	for _, item := range list {
		assert.Equal(t, tenantID, item.TenantID)
	}
}

func TestDataSourceRepository_Delete(t *testing.T) {
	db, repo := setupDataSourceRepo(t)
	ctx := context.Background()
	tenantID := setupTenant(t, db)

	ds := newTestDataSource(testID("ds"), tenantID, "ToDelete")
	require.NoError(t, repo.Create(ctx, ds))

	require.NoError(t, repo.Delete(ctx, ds.ID, tenantID))
	_, err := repo.GetByID(ctx, ds.ID)
	assert.Error(t, err)
}

func TestDataSourceRepository_Search(t *testing.T) {
	db, repo := setupDataSourceRepo(t)
	ctx := context.Background()
	tenantID := setupTenant(t, db)

	require.NoError(t, repo.Create(ctx, newTestDataSource(testID("ds"), tenantID, "MySQL-Production")))
	require.NoError(t, repo.Create(ctx, newTestDataSource(testID("ds"), tenantID, "MySQL-Dev")))
	require.NoError(t, repo.Create(ctx, newTestDataSource(testID("ds"), tenantID, "PostgreSQL-Main")))

	results, total, err := repo.Search(ctx, tenantID, "MySQL", 1, 10)
	require.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, int64(2), total)

	results, total, err = repo.Search(ctx, tenantID, "Production", 1, 10)
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, int64(1), total)

	results, total, err = repo.Search(ctx, tenantID, "", 1, 10)
	require.NoError(t, err)
	assert.Len(t, results, 3)
	assert.Equal(t, int64(3), total)
}

func TestDataSourceRepository_Copy(t *testing.T) {
	db, repo := setupDataSourceRepo(t)
	ctx := context.Background()
	tenantID := setupTenant(t, db)

	original := newTestDataSource(testID("ds"), tenantID, "Original-DS")
	original.Host = "db.example.com"
	original.Port = 5432
	require.NoError(t, repo.Create(ctx, original))

	copy, err := repo.Copy(ctx, original.ID, tenantID)
	require.NoError(t, err)
	assert.NotNil(t, copy)
	assert.NotEqual(t, original.ID, copy.ID)
	assert.Equal(t, "Original-DS (副本)", copy.Name)
	assert.Equal(t, "db.example.com", copy.Host)
	assert.Equal(t, 5432, copy.Port)
	assert.Equal(t, tenantID, copy.TenantID)

	list, _, err := repo.List(ctx, tenantID, 1, 10)
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestDataSourceRepository_Copy_WrongTenant(t *testing.T) {
	db, repo := setupDataSourceRepo(t)
	ctx := context.Background()
	tenantID := setupTenant(t, db)
	otherTenantID := setupTenant(t, db)

	ds := newTestDataSource(testID("ds"), tenantID, "DS")
	require.NoError(t, repo.Create(ctx, ds))

	_, err := repo.Copy(ctx, ds.ID, otherTenantID)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestDataSourceRepository_Move(t *testing.T) {
	db, repo := setupDataSourceRepo(t)
	ctx := context.Background()
	tenantID := setupTenant(t, db)

	ds := newTestDataSource(testID("ds"), tenantID, "MoveTest")
	require.NoError(t, repo.Create(ctx, ds))

	err := repo.Move(ctx, ds.ID, tenantID)
	require.NoError(t, err)
}

func TestDataSourceRepository_Rename(t *testing.T) {
	db, repo := setupDataSourceRepo(t)
	ctx := context.Background()
	tenantID := setupTenant(t, db)

	ds := newTestDataSource(testID("ds"), tenantID, "OldName")
	require.NoError(t, repo.Create(ctx, ds))

	err := repo.Rename(ctx, ds.ID, tenantID, "NewName")
	require.NoError(t, err)

	renamed, err := repo.GetByID(ctx, ds.ID)
	require.NoError(t, err)
	assert.Equal(t, "NewName", renamed.Name)
}

func TestDataSourceRepository_Rename_EmptyName(t *testing.T) {
	db, repo := setupDataSourceRepo(t)
	ctx := context.Background()
	tenantID := setupTenant(t, db)

	ds := newTestDataSource(testID("ds"), tenantID, "Name")
	require.NoError(t, repo.Create(ctx, ds))

	err := repo.Rename(ctx, ds.ID, tenantID, "")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidData, err)
}

func TestDataSourceRepository_Rename_Duplicate(t *testing.T) {
	db, repo := setupDataSourceRepo(t)
	ctx := context.Background()
	tenantID := setupTenant(t, db)

	ds1 := newTestDataSource(testID("ds"), tenantID, "DS1")
	ds2 := newTestDataSource(testID("ds"), tenantID, "DS2")
	require.NoError(t, repo.Create(ctx, ds1))
	require.NoError(t, repo.Create(ctx, ds2))

	err := repo.Rename(ctx, ds2.ID, tenantID, "DS1")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrDuplicatedKey, err)
}

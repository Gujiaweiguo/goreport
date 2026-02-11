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

func setupDataSourceRepo(t *testing.T) (*gorm.DB, DataSourceRepository) {
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

	list, err := repo.List(ctx, tenantID)
	require.NoError(t, err)
	assert.Len(t, list, 2)
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

	require.NoError(t, repo.Delete(ctx, ds.ID))
	_, err := repo.GetByID(ctx, ds.ID)
	assert.Error(t, err)
}

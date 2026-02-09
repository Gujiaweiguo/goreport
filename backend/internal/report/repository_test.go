package report

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gujiaweiguo/goreport/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewRepository(t *testing.T) {
	repo := NewRepository(nil)
	assert.NotNil(t, repo)
}

func setupRepo(t *testing.T) (*gorm.DB, Repository) {
	t.Helper()
	db := testutil.SetupMySQLTestDB(t)
	t.Cleanup(func() {
		testutil.CloseDB(db)
	})
	return db, NewRepository(db)
}

func setupTenant(t *testing.T, db *gorm.DB) string {
	t.Helper()
	tenantID := uniqueID("tenant")
	err := db.Exec(
		"INSERT IGNORE INTO tenants (id, name, code, status) VALUES (?, ?, ?, 1)",
		tenantID,
		"Test "+tenantID,
		tenantID,
	).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		testutil.CleanupTenantData(db, []string{tenantID})
	})

	return tenantID
}

func newTestReport(id, tenantID, name string) *Report {
	return &Report{
		ID:       id,
		TenantID: tenantID,
		Name:     name,
		Code:     "RPT-001",
		Type:     "report",
		Config:   `{"title":"` + name + `"}`,
		Status:   1,
	}
}

func uniqueID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func TestRepository_CreateAndGet(t *testing.T) {
	db, repo := setupRepo(t)
	ctx := context.Background()
	tenantID := setupTenant(t, db)

	report := newTestReport(uniqueID("r"), tenantID, "Test Report")
	require.NoError(t, repo.Create(ctx, report))

	fetched, err := repo.Get(ctx, report.ID, report.TenantID)
	require.NoError(t, err)
	assert.Equal(t, report.ID, fetched.ID)
	assert.Equal(t, "Test Report", fetched.Name)
}

func TestRepository_Update(t *testing.T) {
	db, repo := setupRepo(t)
	ctx := context.Background()
	tenantID := setupTenant(t, db)

	report := newTestReport(uniqueID("r"), tenantID, "Before Update")
	require.NoError(t, repo.Create(ctx, report))

	report.Name = "After Update"
	report.Code = "RPT-002"
	require.NoError(t, repo.Update(ctx, report))

	updated, err := repo.Get(ctx, report.ID, report.TenantID)
	require.NoError(t, err)
	assert.Equal(t, "After Update", updated.Name)
	assert.Equal(t, "RPT-002", updated.Code)
}

func TestRepository_Delete(t *testing.T) {
	db, repo := setupRepo(t)
	ctx := context.Background()
	tenantID := setupTenant(t, db)

	report := newTestReport(uniqueID("r"), tenantID, "To Delete")
	require.NoError(t, repo.Create(ctx, report))
	require.NoError(t, repo.Delete(ctx, report.ID, report.TenantID))

	_, err := repo.Get(ctx, report.ID, report.TenantID)
	assert.Error(t, err)
}

func TestRepository_List_ByTenant(t *testing.T) {
	db, repo := setupRepo(t)
	ctx := context.Background()
	tenantID := setupTenant(t, db)
	otherTenantID := setupTenant(t, db)

	require.NoError(t, repo.Create(ctx, newTestReport(uniqueID("r"), tenantID, "Report A")))
	require.NoError(t, repo.Create(ctx, newTestReport(uniqueID("r"), tenantID, "Report B")))
	require.NoError(t, repo.Create(ctx, newTestReport(uniqueID("r"), otherTenantID, "Other Tenant")))

	list, err := repo.List(ctx, tenantID)
	require.NoError(t, err)
	assert.Len(t, list, 2)

	for _, r := range list {
		assert.Equal(t, tenantID, r.TenantID)
	}
}

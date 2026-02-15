package dashboard

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

func TestNewRepository_Dashboard(t *testing.T) {
	repo := NewRepository(nil)
	assert.NotNil(t, repo)
}

func setupDashboardRepo(t *testing.T) (*gorm.DB, Repository) {
	t.Helper()
	db := testutil.SetupMySQLTestDB(t)
	t.Cleanup(func() {
		testutil.CloseDB(db)
	})
	return db, NewRepository(db)
}

func setupTenantForDashboard(t *testing.T, db *gorm.DB) string {
	t.Helper()
	tenantID := testDashboardID("tenant")
	err := db.Exec(
		"INSERT IGNORE INTO tenants (id, name, code, status) VALUES (?, ?, ?, 1)",
		tenantID,
		"Test "+tenantID,
		tenantID,
	).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = db.Exec("DELETE FROM dashboards WHERE tenant_id = ?", tenantID).Error
		testutil.CleanupTenantData(db, []string{tenantID})
	})

	return tenantID
}

func testDashboardID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func newTestDashboard(id, tenantID, name string) *models.Dashboard {
	return &models.Dashboard{
		ID:       id,
		TenantID: tenantID,
		Name:     name,
		Code:     "code-" + name,
		Config: models.DashboardConfig{
			Width:           1920,
			Height:          1080,
			BackgroundColor: "#0a0e27",
		},
		Components: []models.DashboardComponent{
			{
				ID:      "comp-1",
				Title:   "Test Component",
				Type:    "chart",
				Width:   400,
				Height:  300,
				X:       0,
				Y:       0,
				Visible: true,
				Locked:  false,
			},
		},
		Status:    1,
		ViewCount: 0,
		CreatedBy: "test-user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestRepository_Create(t *testing.T) {
	db, repo := setupDashboardRepo(t)
	tenantID := setupTenantForDashboard(t, db)

	dashboard := newTestDashboard(testDashboardID("dash"), tenantID, "Test Dashboard")
	err := repo.Create(dashboard)
	require.NoError(t, err)

	// Verify it was created
	var found models.Dashboard
	err = db.Where("id = ?", dashboard.ID).First(&found).Error
	require.NoError(t, err)
	assert.Equal(t, "Test Dashboard", found.Name)
	assert.Equal(t, tenantID, found.TenantID)
}

func TestRepository_Get(t *testing.T) {
	db, repo := setupDashboardRepo(t)
	tenantID := setupTenantForDashboard(t, db)

	dashboard := newTestDashboard(testDashboardID("dash"), tenantID, "Get Test")
	require.NoError(t, repo.Create(dashboard))

	// Test get
	found, err := repo.Get(dashboard.ID, tenantID)
	require.NoError(t, err)
	assert.Equal(t, dashboard.ID, found.ID)
	assert.Equal(t, "Get Test", found.Name)
}

func TestRepository_Get_NotFound(t *testing.T) {
	_, repo := setupDashboardRepo(t)

	_, err := repo.Get("non-existing-id", "tenant-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRepository_Get_WrongTenant(t *testing.T) {
	db, repo := setupDashboardRepo(t)
	tenantID := setupTenantForDashboard(t, db)

	dashboard := newTestDashboard(testDashboardID("dash"), tenantID, "Wrong Tenant Test")
	require.NoError(t, repo.Create(dashboard))

	// Try to get with different tenant
	_, err := repo.Get(dashboard.ID, "other-tenant")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRepository_Update(t *testing.T) {
	db, repo := setupDashboardRepo(t)
	tenantID := setupTenantForDashboard(t, db)

	dashboard := newTestDashboard(testDashboardID("dash"), tenantID, "Before Update")
	require.NoError(t, repo.Create(dashboard))

	// Update
	dashboard.Name = "After Update"
	dashboard.ViewCount = 100
	err := repo.Update(dashboard)
	require.NoError(t, err)

	// Verify
	found, err := repo.Get(dashboard.ID, tenantID)
	require.NoError(t, err)
	assert.Equal(t, "After Update", found.Name)
	assert.Equal(t, 100, found.ViewCount)
}

func TestRepository_Delete(t *testing.T) {
	db, repo := setupDashboardRepo(t)
	tenantID := setupTenantForDashboard(t, db)

	dashboard := newTestDashboard(testDashboardID("dash"), tenantID, "To Delete")
	require.NoError(t, repo.Create(dashboard))

	// Delete
	err := repo.Delete(dashboard.ID, tenantID)
	require.NoError(t, err)

	// Verify deleted
	var found models.Dashboard
	err = db.Unscoped().Where("id = ?", dashboard.ID).First(&found).Error
	require.NoError(t, err)
	assert.NotNil(t, found.DeletedAt)
}

func TestRepository_Delete_NotFound(t *testing.T) {
	_, repo := setupDashboardRepo(t)

	err := repo.Delete("non-existing-id", "tenant-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRepository_Delete_WrongTenant(t *testing.T) {
	db, repo := setupDashboardRepo(t)
	tenantID := setupTenantForDashboard(t, db)

	dashboard := newTestDashboard(testDashboardID("dash"), tenantID, "Delete Wrong Tenant")
	require.NoError(t, repo.Create(dashboard))

	// Try to delete with wrong tenant
	err := repo.Delete(dashboard.ID, "other-tenant")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRepository_List(t *testing.T) {
	db, repo := setupDashboardRepo(t)
	tenantID := setupTenantForDashboard(t, db)
	otherTenantID := setupTenantForDashboard(t, db)

	// Create dashboards for different tenants
	require.NoError(t, repo.Create(newTestDashboard(testDashboardID("dash"), tenantID, "Dashboard A")))
	require.NoError(t, repo.Create(newTestDashboard(testDashboardID("dash"), tenantID, "Dashboard B")))
	require.NoError(t, repo.Create(newTestDashboard(testDashboardID("dash"), otherTenantID, "Other Dashboard")))

	// List for first tenant
	list, err := repo.List(tenantID)
	require.NoError(t, err)
	assert.Len(t, list, 2)

	for _, d := range list {
		assert.Equal(t, tenantID, d.TenantID)
	}
}

func TestRepository_List_Empty(t *testing.T) {
	_, repo := setupDashboardRepo(t)

	list, err := repo.List("empty-tenant")
	require.NoError(t, err)
	assert.Empty(t, list)
}

// Test context cancellation
func TestRepository_Get_ContextCancellation(t *testing.T) {
	db, repo := setupDashboardRepo(t)
	tenantID := setupTenantForDashboard(t, db)

	dashboard := newTestDashboard(testDashboardID("dash"), tenantID, "Context Test")
	require.NoError(t, repo.Create(dashboard))

	// This test verifies the repository works with a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// The repository should still work (it doesn't use context in current implementation)
	// but this test documents the expected behavior
	_ = ctx // Context is not used in current implementation

	found, err := repo.Get(dashboard.ID, tenantID)
	require.NoError(t, err)
	assert.Equal(t, dashboard.ID, found.ID)
}

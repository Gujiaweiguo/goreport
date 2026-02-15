package testutil

import (
	"os"
	"strings"
	"testing"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	if dsn == "" {
		t.Skip("TEST_DB_DSN or DB_DSN not set")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	return db
}

func TestNewDatasetFixturesAndGetters(t *testing.T) {
	f := NewDatasetFixtures()
	require.NotNil(t, f)
	assert.Len(t, f.Datasets, 2)
	assert.Len(t, f.Fields, 8)
	assert.Len(t, f.Sources, 2)

	assert.NotNil(t, f.GetDatasetByID("dataset-001"))
	assert.Nil(t, f.GetDatasetByID("not-found"))

	assert.Len(t, f.GetFieldsByDatasetID("dataset-001"), 7)
	assert.Len(t, f.GetFieldsByDatasetID("dataset-002"), 1)

	assert.Len(t, f.GetDimensions("dataset-001"), 3)
	assert.Len(t, f.GetMeasures("dataset-001"), 4)
	assert.Len(t, f.GetComputedFields("dataset-001"), 2)
}

func TestNewDatasourceFixturesAndGetters(t *testing.T) {
	f := NewDatasourceFixtures()
	require.NotNil(t, f)
	assert.Len(t, f.Datasources, 2)
	assert.NotEmpty(t, f.TenantID)
	assert.NotEmpty(t, f.UserID)

	assert.NotNil(t, f.GetDatasourceByID("ds-test-001"))
	assert.Nil(t, f.GetDatasourceByID("not-found"))

	assert.NotNil(t, f.GetDatasourceByName("Test MySQL"))
	assert.Nil(t, f.GetDatasourceByName("not-found"))
}

func TestNewAuthFixturesAndGetUser(t *testing.T) {
	f := NewAuthFixtures()
	require.NotNil(t, f)
	assert.Len(t, f.Tenants, 1)
	assert.Len(t, f.Users, 2)

	u := f.GetUserByUsername("testuser")
	require.NotNil(t, u)
	assert.Equal(t, "user", u.Role)

	assert.Nil(t, f.GetUserByUsername("not-found"))
}

func TestSanitizeTestName(t *testing.T) {
	got := sanitizeTestName("Test_A/B:Case#1")
	assert.Equal(t, "Test-A-BCase1", got)

	long := sanitizeTestName(strings.Repeat("a", 80))
	assert.Len(t, long, 32)
}

func TestGenerateUniqueID(t *testing.T) {
	id1 := GenerateUniqueID("pref")
	id2 := GenerateUniqueID("pref")
	assert.True(t, strings.HasPrefix(id1, "pref-"))
	assert.True(t, strings.HasPrefix(id2, "pref-"))
	assert.NotEqual(t, id1, id2)
}

func TestEnsureTenant(t *testing.T) {
	db := openTestDB(t)
	defer CloseDB(db)

	tenantID := GenerateUniqueID("tenant")
	err := EnsureTenant(db, tenantID)
	require.NoError(t, err)

	var tenant models.Tenant
	err = db.Where("id = ?", tenantID).First(&tenant).Error
	require.NoError(t, err)
	assert.Equal(t, tenantID, tenant.ID)

	_ = db.Unscoped().Where("id = ?", tenantID).Delete(&models.Tenant{}).Error
}

func TestDatasourceFixturesSetupCleanup(t *testing.T) {
	db := openTestDB(t)
	defer CloseDB(db)

	f := NewDatasourceFixtures()
	f.TenantID = GenerateUniqueID("tenant-ds")
	f.UserID = GenerateUniqueID("user-ds")
	for i := range f.Datasources {
		f.Datasources[i].ID = GenerateUniqueID("ds")
		f.Datasources[i].TenantID = f.TenantID
		f.Datasources[i].CreatedBy = f.UserID
	}

	err := f.Setup(db)
	require.NoError(t, err)

	for _, ds := range f.Datasources {
		var found models.DataSource
		err = db.Where("id = ?", ds.ID).First(&found).Error
		assert.NoError(t, err)
	}

	err = f.Cleanup(db)
	require.NoError(t, err)
	for _, ds := range f.Datasources {
		var found models.DataSource
		err = db.Where("id = ?", ds.ID).First(&found).Error
		assert.Error(t, err)
	}
}

func TestAuthFixturesSetupAndCleanup(t *testing.T) {
	db := openTestDB(t)
	defer CloseDB(db)

	f := NewAuthFixtures()
	tenantID := GenerateUniqueID("tenant-auth")
	f.Tenants[0].ID = tenantID
	f.Tenants[0].Code = tenantID
	f.Users[0].ID = GenerateUniqueID("user")
	f.Users[0].Username = GenerateUniqueID("user")
	f.Users[0].TenantID = tenantID
	f.Users[1].ID = GenerateUniqueID("user")
	f.Users[1].Username = GenerateUniqueID("admin")
	f.Users[1].TenantID = tenantID

	err := f.Setup(db)
	require.NoError(t, err)

	for _, u := range f.Users {
		var found models.User
		err = db.Where("id = ?", u.ID).First(&found).Error
		assert.NoError(t, err)
	}

	err = f.Cleanup(db)
	require.NoError(t, err)

	_ = db.Where("id IN ?", []string{f.Users[0].ID, f.Users[1].ID}).Delete(&models.User{}).Error
	_ = db.Where("id = ?", tenantID).Delete(&models.Tenant{}).Error
}

func TestNewTenantTestContext(t *testing.T) {
	db := openTestDB(t)

	tc := NewTenantTestContext(t, db)
	require.NotNil(t, tc)
	assert.NotEmpty(t, tc.TenantID)
	assert.NotEmpty(t, tc.UserID)
	assert.NotNil(t, tc.DB)

	var tenant models.Tenant
	err := db.Where("id = ?", tc.TenantID).First(&tenant).Error
	assert.NoError(t, err)
}

func TestSetupMySQLTestDB(t *testing.T) {
	db := SetupMySQLTestDB(t)
	require.NotNil(t, db)
	CloseDB(db)
}

func TestSetupRepositoryTestDB(t *testing.T) {
	t.Skip("Skipping due to migration issues with existing database schema")
	db := SetupRepositoryTestDB(t)
	require.NotNil(t, db)

	hasTenants := db.Migrator().HasTable("tenants")
	assert.True(t, hasTenants)
}

func TestEnsureTenants(t *testing.T) {
	db := openTestDB(t)
	defer CloseDB(db)

	EnsureTenants(db, t)

	for _, tenantID := range []string{"test-tenant", "tenant-a", "tenant-1"} {
		var tenant models.Tenant
		err := db.Where("id = ?", tenantID).First(&tenant).Error
		assert.NoError(t, err, "tenant %s should exist", tenantID)
	}
}

func TestCleanupTenantData(t *testing.T) {
	db := openTestDB(t)
	defer CloseDB(db)

	tenantID := GenerateUniqueID("cleanup-test")

	err := EnsureTenant(db, tenantID)
	require.NoError(t, err)

	dataset := &models.Dataset{
		ID:       GenerateUniqueID("dataset"),
		TenantID: tenantID,
		Name:     "Test Dataset",
		Type:     "sql",
		Status:   1,
		Config:   "{}",
	}
	err = db.Create(dataset).Error
	require.NoError(t, err)

	CleanupTenantData(db, []string{tenantID})

	var count int64
	db.Model(&models.Dataset{}).Where("tenant_id = ?", tenantID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestCloseDB(t *testing.T) {
	db := openTestDB(t)
	CloseDB(db)

	sqlDB, err := db.DB()
	if err == nil {
		assert.Error(t, sqlDB.Ping(), "connection should be closed")
	}
}

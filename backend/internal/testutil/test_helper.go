package testutil

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupMySQLTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	if dsn == "" {
		t.Skip("TEST_DB_DSN or DB_DSN not set")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return db
}

func SetupRepositoryTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := SetupMySQLTestDB(t)

	err := db.AutoMigrate(
		&models.Tenant{},
		&models.User{},
		&models.DataSource{},
		&models.Dataset{},
		&models.DatasetField{},
		&models.DatasetSource{},
		&models.Dashboard{},
		&models.Chart{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	t.Cleanup(func() {
		CloseDB(db)
	})

	return db
}

func EnsureTenants(db *gorm.DB, t *testing.T) {
	t.Helper()
	tenantIDs := []string{"test-tenant", "tenant-a", "tenant-1"}
	now := time.Now()

	for _, tenantID := range tenantIDs {
		var tenant models.Tenant
		err := db.Where("id = ?", tenantID).First(&tenant).Error
		if err == gorm.ErrRecordNotFound {
			tenant = models.Tenant{
				ID:        tenantID,
				Name:      "Test " + tenantID,
				Code:      tenantID,
				Status:    1,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if err := db.Create(&tenant).Error; err != nil {
				t.Fatalf("Failed to prepare tenant fixture: %v", err)
			}
		} else if err != nil {
			t.Fatalf("Failed to check tenant: %v", err)
		}
	}
}

func CleanupTenantData(db *gorm.DB, tenantIDs []string) {
	hasDashboards := db.Migrator().HasTable("dashboards")
	hasDatasets := db.Migrator().HasTable("datasets")
	hasDatasetFields := db.Migrator().HasTable("dataset_fields")
	hasDatasetSources := db.Migrator().HasTable("dataset_sources")
	hasReports := db.Migrator().HasTable("reports")

	for _, tenantID := range tenantIDs {
		if hasDatasetFields && hasDatasets {
			_ = db.Exec("DELETE df FROM dataset_fields df INNER JOIN datasets d ON df.dataset_id = d.id WHERE d.tenant_id = ?", tenantID).Error
		}
		if hasDatasetSources && hasDatasets {
			_ = db.Exec("DELETE ds FROM dataset_sources ds INNER JOIN datasets d ON ds.dataset_id = d.id WHERE d.tenant_id = ?", tenantID).Error
		}
		if hasDatasets {
			_ = db.Exec("DELETE FROM datasets WHERE tenant_id = ?", tenantID).Error
		}
		if hasDashboards {
			_ = db.Exec("DELETE FROM dashboards WHERE tenant_id = ?", tenantID).Error
		}
		if hasReports {
			_ = db.Exec("DELETE FROM reports WHERE tenant_id = ?", tenantID).Error
		}
	}
}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		_ = sqlDB.Close()
	}
}

// TestFixture is the interface for test data factories.
// Fixture implementations provide setup and cleanup operations for test data.
type TestFixture interface {
	// Setup creates test data in the database
	Setup(db *gorm.DB) error
	// Cleanup removes test data from the database
	Cleanup(db *gorm.DB) error
}

// TenantTestContext provides an isolated test context for tenant-scoped tests.
// Each context has unique tenant and user IDs to support parallel test execution.
type TenantTestContext struct {
	TenantID string
	UserID   string
	DB       *gorm.DB
}

// NewTenantTestContext creates a new isolated test context with unique IDs.
// The tenant ID and user ID are generated using the test name to ensure uniqueness.
func NewTenantTestContext(t *testing.T, db *gorm.DB) *TenantTestContext {
	t.Helper()

	// Generate unique IDs based on test name
	testName := t.Name()
	sanitizedName := sanitizeTestName(testName)
	tenantID := "test-tenant-" + sanitizedName
	userID := "test-user-" + sanitizedName

	ctx := &TenantTestContext{
		TenantID: tenantID,
		UserID:   userID,
		DB:       db,
	}

	// Ensure tenant exists
	if err := EnsureTenant(db, tenantID); err != nil {
		t.Fatalf("failed to ensure tenant %s: %v", tenantID, err)
	}

	// Register cleanup
	t.Cleanup(func() {
		CleanupTenantData(db, []string{tenantID})
	})

	return ctx
}

// sanitizeTestName converts test name to a safe string for use in IDs
func sanitizeTestName(name string) string {
	result := make([]byte, 0, len(name))
	for _, c := range name {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' {
			result = append(result, byte(c))
		} else if c == '_' || c == '/' {
			result = append(result, '-')
		}
	}
	// Limit length
	if len(result) > 32 {
		result = result[:32]
	}
	return string(result)
}

// EnsureTenant ensures a tenant exists in the database
func EnsureTenant(db *gorm.DB, tenantID string) error {
	now := time.Now()
	var tenant models.Tenant
	err := db.Where("id = ?", tenantID).First(&tenant).Error
	if err == gorm.ErrRecordNotFound {
		tenant = models.Tenant{
			ID:        tenantID,
			Name:      "Test " + tenantID,
			Code:      tenantID,
			Status:    1,
			CreatedAt: now,
			UpdatedAt: now,
		}
		return db.Create(&tenant).Error
	}
	return err
}

// GenerateUniqueID generates a unique ID for test data
func GenerateUniqueID(prefix string) string {
	b := make([]byte, 3)
	rand.Read(b)
	return prefix + "-" + time.Now().Format("20060102150405") + "-" + hex.EncodeToString(b)
}

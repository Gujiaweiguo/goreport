package testutil

import (
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
		if hasDashboards {
			_ = db.Exec("DELETE FROM dashboards WHERE tenant_id = ?", tenantID).Error
		}
		if hasDatasets {
			_ = db.Exec("DELETE FROM datasets WHERE tenant_id = ?", tenantID).Error
		}
		if hasDatasetFields && hasDatasets {
			_ = db.Exec("DELETE FROM dataset_fields WHERE dataset_id IN (SELECT id FROM datasets WHERE tenant_id = ?)", tenantID).Error
		}
		if hasDatasetSources && hasDatasets {
			_ = db.Exec("DELETE FROM dataset_sources WHERE dataset_id IN (SELECT id FROM datasets WHERE tenant_id = ?)", tenantID).Error
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

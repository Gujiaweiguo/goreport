package testutil

import (
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SetupMySQLTestDB 在 MySQL 环境中创建测试数据库连接
// 优先使用 TEST_DB_DSN，回退到 DB_DSN
// 如果两者都不存在，则跳过测试
func SetupMySQLTestDB(t *testing.T) *gorm.DB {
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

// EnsureTenants 确保测试所需的租户数据存在
// 如果租户不存在，则创建它们
func EnsureTenants(db *gorm.DB, t *testing.T) {
	tenantIDs := []string{"test-tenant", "tenant-a", "tenant-1"}

	for _, tenantID := range tenantIDs {
		err := db.Exec(
			"INSERT IGNORE INTO tenants (id, name, code, status) VALUES (?, ?, ?, 1)",
			tenantID,
			"Test "+tenantID,
			tenantID,
		).Error
		if err != nil {
			t.Fatalf("Failed to prepare tenant fixture: %v", err)
		}
	}
}

// CleanupTenantData 清理指定租户的测试数据
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

// CloseDB 安全关闭数据库连接
func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		_ = sqlDB.Close()
	}
}

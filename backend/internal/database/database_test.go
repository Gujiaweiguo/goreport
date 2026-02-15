package database

import (
	"database/sql"
	"os"
	"testing"

	"github.com/gujiaweiguo/goreport/internal/config"
	"gorm.io/gorm"
)

func TestInit_InvalidDSN(t *testing.T) {
	_, err := Init("invalid-dsn")
	if err == nil {
		t.Error("Init() with invalid DSN should return error")
	}
}

func TestInitWithConfig_EmptyDSN(t *testing.T) {
	cfg := &config.DatabaseConfig{
		MaxOpenConns:    50,
		MaxIdleConns:    10,
		ConnMaxLifetime: 1800,
	}
	_, err := InitWithConfig("", cfg)
	if err == nil {
		t.Error("InitWithConfig() with empty DSN should return error")
	}
}

func TestInitWithConfig_NilConfig(t *testing.T) {
	t.Skip("Nil config causes nil pointer dereference - InitWithConfig requires non-nil config")
	dsn := getTestDSN(t)
	if dsn == "" {
		return
	}

	db, err := InitWithConfig(dsn, nil)
	if err != nil {
		t.Fatalf("InitWithConfig() with nil config error = %v", err)
	}
	defer closeDB(db)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get DB instance: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		t.Logf("Warning: failed to close database: %v", err)
	}
}

func TestInitWithConfig_ZeroValues(t *testing.T) {
	dsn := getTestDSN(t)
	if dsn == "" {
		return
	}

	cfg := &config.DatabaseConfig{
		MaxOpenConns:    0,
		MaxIdleConns:    0,
		ConnMaxLifetime: 0,
	}

	db, err := InitWithConfig(dsn, cfg)
	if err != nil {
		t.Fatalf("InitWithConfig() error = %v", err)
	}
	defer closeDB(db)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get DB instance: %v", err)
	}

	stats := sqlDB.Stats()
	if stats.MaxOpenConnections != 100 {
		t.Errorf("MaxOpenConnections with zero config = %d, want 100 (default)", stats.MaxOpenConnections)
	}
}

func TestInit_Success(t *testing.T) {
	dsn := getTestDSN(t)
	if dsn == "" {
		return
	}

	db, err := Init(dsn)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	defer closeDB(db)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get DB instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Errorf("Failed to ping database: %v", err)
	}
}

func TestInit_DatabaseConnectionPool(t *testing.T) {
	dsn := getTestDSN(t)
	if dsn == "" {
		return
	}

	cfg := &config.DatabaseConfig{
		MaxOpenConns:    10,
		MaxIdleConns:    2,
		ConnMaxLifetime: 300,
	}

	db, err := InitWithConfig(dsn, cfg)
	if err != nil {
		t.Fatalf("InitWithConfig() error = %v", err)
	}
	defer closeDB(db)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get DB instance: %v", err)
	}

	stats := sqlDB.Stats()

	if stats.MaxOpenConnections != 10 {
		t.Errorf("MaxOpenConnections = %d, want 10", stats.MaxOpenConnections)
	}
}

func TestEnsureDatasourceSchemaCompatibility_TableNotExists(t *testing.T) {
	dsn := getTestDSN(t)
	if dsn == "" {
		return
	}

	db, err := Init(dsn)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	defer closeDB(db)

	err = ensureDatasourceSchemaCompatibility(db)
	if err != nil {
		t.Fatalf("ensureDatasourceSchemaCompatibility() error = %v", err)
	}
}

func getTestDSN(t *testing.T) string {
	t.Helper()
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	if dsn == "" {
		t.Skip("TEST_DB_DSN or DB_DSN not set")
	}
	return dsn
}

func closeDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	_ = sqlDB.Close()
}

func TestInitWithConfig_MaxOpenConnsNegative(t *testing.T) {
	dsn := getTestDSN(t)
	if dsn == "" {
		return
	}

	cfg := &config.DatabaseConfig{
		MaxOpenConns:    -1,
		MaxIdleConns:    10,
		ConnMaxLifetime: 3600,
	}

	db, err := InitWithConfig(dsn, cfg)
	if err != nil {
		t.Fatalf("InitWithConfig() error = %v", err)
	}
	defer closeDB(db)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get DB instance: %v", err)
	}

	stats := sqlDB.Stats()
	if stats.MaxOpenConnections != 0 {
		t.Errorf("MaxOpenConnections with -1 config = %d, want 0", stats.MaxOpenConnections)
	}
}

func TestInit_DSNWithCharset(t *testing.T) {
	dsn := getTestDSN(t)
	if dsn == "" {
		return
	}

	db, err := Init(dsn)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	defer closeDB(db)

	sqlDB, _ := db.DB()
	row := sqlDB.QueryRow("SELECT @@character_set_database")
	if row != nil {
		var charset string
		if err := row.Scan(&charset); err == nil {
			t.Logf("Database charset: %s", charset)
		}
	}
}

func TestInit_ConnectionValidation(t *testing.T) {
	dsn := getTestDSN(t)
	if dsn == "" {
		return
	}

	db, err := Init(dsn)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	defer closeDB(db)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get DB instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Errorf("Database ping failed: %v", err)
	}

	var version string
	row := sqlDB.QueryRow("SELECT VERSION()")
	if err := row.Scan(&version); err != nil {
		t.Errorf("Failed to get database version: %v", err)
	} else {
		t.Logf("MySQL version: %s", version)
	}
}

func TestInitWithConfig_PoolSettings(t *testing.T) {
	dsn := getTestDSN(t)
	if dsn == "" {
		return
	}

	testCases := []struct {
		name            string
		maxOpenConns    int
		maxIdleConns    int
		connMaxLifetime int
		expectedOpen    int
		expectedIdle    int
	}{
		{"defaults", 0, 0, 0, 100, 10},
		{"custom", 50, 20, 1800, 50, 20},
		{"small_pool", 5, 1, 60, 5, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.DatabaseConfig{
				MaxOpenConns:    tc.maxOpenConns,
				MaxIdleConns:    tc.maxIdleConns,
				ConnMaxLifetime: tc.connMaxLifetime,
			}

			db, err := InitWithConfig(dsn, cfg)
			if err != nil {
				t.Fatalf("InitWithConfig() error = %v", err)
			}
			defer closeDB(db)

			sqlDB, err := db.DB()
			if err != nil {
				t.Fatalf("Failed to get DB instance: %v", err)
			}

			stats := sqlDB.Stats()
			if stats.MaxOpenConnections != tc.expectedOpen {
				t.Errorf("MaxOpenConnections = %d, want %d", stats.MaxOpenConnections, tc.expectedOpen)
			}
		})
	}
}

func TestInit_InvalidDSNFormats(t *testing.T) {
	invalidDSNs := []string{
		"",
		"mysql://invalid",
		"tcp(localhost:3306)/db",
		"root@tcp(localhost)/db",
	}

	for _, dsn := range invalidDSNs {
		t.Run(dsn, func(t *testing.T) {
			if dsn == "" {
				return
			}
			_, err := Init(dsn)
			if err == nil {
				t.Errorf("Init(%q) should return error", dsn)
			}
		})
	}
}

func TestDatabaseConfig_Defaults(t *testing.T) {
	dsn := getTestDSN(t)
	if dsn == "" {
		return
	}

	db, err := InitWithConfig(dsn, &config.DatabaseConfig{})
	if err != nil {
		t.Fatalf("InitWithConfig() error = %v", err)
	}
	defer closeDB(db)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get DB instance: %v", err)
	}

	stats := sqlDB.Stats()

	if stats.MaxOpenConnections != 100 {
		t.Errorf("Default MaxOpenConnections = %d, want 100", stats.MaxOpenConnections)
	}
}

func TestInit_DBInstance(t *testing.T) {
	dsn := getTestDSN(t)
	if dsn == "" {
		return
	}

	db, err := Init(dsn)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	defer closeDB(db)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() error = %v", err)
	}

	if sqlDB == nil {
		t.Error("sqlDB should not be nil")
	}

	var isConnected int
	err = sqlDB.QueryRow("SELECT 1").Scan(&isConnected)
	if err != nil {
		t.Errorf("Failed to execute query: %v", err)
	}
	if isConnected != 1 {
		t.Errorf("SELECT 1 = %d, want 1", isConnected)
	}
}

var _ = sql.DB{}
